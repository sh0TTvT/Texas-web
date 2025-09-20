// 房间管理系统
// 作用：管理游戏房间，处理玩家加入/离开，维护房间状态，协调游戏流程

package room

import (
	"fmt"
	"sync"
	"time"

	"texas-poker-backend/internal/game/poker"
	"texas-poker-backend/internal/game/statemachine"
)

// RoomStatus 房间状态
type RoomStatus string

const (
	RoomWaiting RoomStatus = "waiting" // 等待玩家
	RoomPlaying RoomStatus = "playing" // 游戏进行中
	RoomClosed  RoomStatus = "closed"  // 房间关闭
)

// Player 房间内玩家信息
type Player struct {
	ID       int64             `json:"id"`
	Username string            `json:"username"`
	Chips    int               `json:"chips"`
	Position int               `json:"position"` // 座位位置（0-5）
	Status   PlayerStatus      `json:"status"`
	Cards    []poker.Card      `json:"cards,omitempty"` // 手牌（只有本人可见）
	LastAction statemachine.PlayerAction `json:"last_action,omitempty"`
	BetAmount int               `json:"bet_amount"` // 本局已下注金额
	IsDealer bool              `json:"is_dealer"`  // 是否是庄家
	IsSmallBlind bool          `json:"is_small_blind"` // 是否是小盲注
	IsBigBlind bool            `json:"is_big_blind"`   // 是否是大盲注
	JoinTime time.Time         `json:"join_time"`
}

// PlayerStatus 玩家状态
type PlayerStatus string

const (
	PlayerActive   PlayerStatus = "active"   // 活跃
	PlayerFolded   PlayerStatus = "folded"   // 已弃牌
	PlayerAllIn    PlayerStatus = "allin"    // 全押
	PlayerSitting  PlayerStatus = "sitting"  // 坐下但未参与游戏
	PlayerWaiting  PlayerStatus = "waiting"  // 等待下一局
)

// Room 房间结构
type Room struct {
	ID              int64                         `json:"id"`
	Name            string                        `json:"name"`
	ChipLevel       string                        `json:"chip_level"`
	MinChips        int                           `json:"min_chips"`
	SmallBlind      int                           `json:"small_blind"`
	BigBlind        int                           `json:"big_blind"`
	MaxPlayers      int                           `json:"max_players"`
	IsPrivate       bool                          `json:"is_private"`
	Status          RoomStatus                    `json:"status"`
	Players         map[int64]*Player             `json:"players"`
	CommunityCards  []poker.Card                  `json:"community_cards"`
	Pot             int                           `json:"pot"` // 底池
	CurrentGame     *GameSession                  `json:"current_game,omitempty"`
	StateMachine    *statemachine.GameStateMachine `json:"-"`
	BettingRound    *statemachine.BettingRound    `json:"-"`
	Deck            *poker.Deck                   `json:"-"` // 牌堆
	DealerPosition  int                           `json:"dealer_position"` // 庄家位置
	CreatedAt       time.Time                     `json:"created_at"`
	UpdatedAt       time.Time                     `json:"updated_at"`
	
	// 并发安全
	mu sync.RWMutex `json:"-"`
}

// GameSession 游戏会话
type GameSession struct {
	ID          string                     `json:"id"`
	StartTime   time.Time                  `json:"start_time"`
	Participants []int64                    `json:"participants"` // 参与此局的玩家ID
	GameLog     []string                   `json:"game_log"`     // 游戏日志
	WinnerID    int64                      `json:"winner_id,omitempty"`
	WinAmount   int                        `json:"win_amount,omitempty"`
}

// NewRoom 创建新房间
func NewRoom(id int64, name, chipLevel string, minChips, smallBlind, bigBlind, maxPlayers int, isPrivate bool) *Room {
	room := &Room{
		ID:             id,
		Name:           name,
		ChipLevel:      chipLevel,
		MinChips:       minChips,
		SmallBlind:     smallBlind,
		BigBlind:       bigBlind,
		MaxPlayers:     maxPlayers,
		IsPrivate:      isPrivate,
		Status:         RoomWaiting,
		Players:        make(map[int64]*Player),
		CommunityCards: make([]poker.Card, 0, 5),
		Pot:            0,
		StateMachine:   statemachine.NewGameStateMachine(),
		DealerPosition: 0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	
	// 设置状态机回调
	room.setupStateMachineCallbacks()
	
	return room
}

// setupStateMachineCallbacks 设置状态机回调
func (r *Room) setupStateMachineCallbacks() {
	r.StateMachine.SetStateCallback(statemachine.PreFlop, func() error {
		return r.startPreFlop()
	})
	
	r.StateMachine.SetStateCallback(statemachine.Flop, func() error {
		return r.dealFlop()
	})
	
	r.StateMachine.SetStateCallback(statemachine.Turn, func() error {
		return r.dealTurn()
	})
	
	r.StateMachine.SetStateCallback(statemachine.River, func() error {
		return r.dealRiver()
	})
	
	r.StateMachine.SetStateCallback(statemachine.Showdown, func() error {
		return r.showdown()
	})
	
	r.StateMachine.SetStateCallback(statemachine.GameEnd, func() error {
		return r.endGame()
	})
}

// AddPlayer 添加玩家到房间
func (r *Room) AddPlayer(userID int64, username string, chips int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 检查房间是否已满
	if len(r.Players) >= r.MaxPlayers {
		return fmt.Errorf("房间已满")
	}
	
	// 检查玩家是否已在房间中
	if _, exists := r.Players[userID]; exists {
		return fmt.Errorf("玩家已在房间中")
	}
	
	// 检查筹码是否满足最低要求
	if chips < r.MinChips {
		return fmt.Errorf("筹码不足，最低需要 %d", r.MinChips)
	}
	
	// 找到空闲位置
	position := r.findAvailablePosition()
	if position == -1 {
		return fmt.Errorf("没有可用位置")
	}
	
	// 创建玩家
	player := &Player{
		ID:       userID,
		Username: username,
		Chips:    chips,
		Position: position,
		Status:   PlayerSitting,
		Cards:    make([]poker.Card, 0, 2),
		JoinTime: time.Now(),
	}
	
	r.Players[userID] = player
	r.UpdatedAt = time.Now()
	
	// 如果达到最少玩家数量且房间在等待状态，可以开始游戏
	if len(r.Players) >= 2 && r.Status == RoomWaiting {
		r.Status = RoomWaiting // 保持等待状态，等待手动开始
	}
	
	return nil
}

// RemovePlayer 从房间移除玩家
func (r *Room) RemovePlayer(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	player, exists := r.Players[userID]
	if !exists {
		return fmt.Errorf("玩家不在房间中")
	}
	
	// 如果游戏正在进行，需要特殊处理
	if r.Status == RoomPlaying && player.Status == PlayerActive {
		// 玩家自动弃牌
		player.Status = PlayerFolded
		r.logGameAction(fmt.Sprintf("玩家 %s 离开房间，自动弃牌", player.Username))
		
		// 检查是否需要触发游戏事件
		if r.BettingRound != nil {
			// 触发玩家离开事件
			r.StateMachine.Transition(statemachine.PlayerLeft)
		}
	}
	
	delete(r.Players, userID)
	r.UpdatedAt = time.Now()
	
	// 如果房间空了，设置为等待状态
	if len(r.Players) == 0 {
		r.Status = RoomWaiting
		r.StateMachine.Reset()
	}
	
	return nil
}

// StartGame 开始游戏
func (r *Room) StartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 检查是否可以开始游戏
	if len(r.Players) < 2 {
		return fmt.Errorf("至少需要2名玩家才能开始游戏")
	}
	
	if r.Status != RoomWaiting {
		return fmt.Errorf("房间状态不允许开始游戏")
	}
	
	// 创建新的游戏会话
	r.CurrentGame = &GameSession{
		ID:           fmt.Sprintf("game_%d_%d", r.ID, time.Now().Unix()),
		StartTime:    time.Now(),
		Participants: r.getActivePlayerIDs(),
		GameLog:      make([]string, 0),
	}
	
	// 初始化牌堆
	r.Deck = poker.NewDeck()
	r.Deck.Shuffle()
	
	// 重置房间状态
	r.resetRoomState()
	
	// 设置盲注
	r.setupBlinds()
	
	// 更新房间状态
	r.Status = RoomPlaying
	
	// 触发游戏开始事件
	return r.StateMachine.Transition(statemachine.StartGame)
}

// ProcessPlayerAction 处理玩家操作
func (r *Room) ProcessPlayerAction(userID int64, action statemachine.PlayerAction, amount int) (statemachine.ActionResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// 检查玩家是否在房间中
	player, exists := r.Players[userID]
	if !exists {
		return statemachine.ActionResult{}, fmt.Errorf("玩家不在房间中")
	}
	
	// 检查游戏状态
	if r.Status != RoomPlaying {
		return statemachine.ActionResult{}, fmt.Errorf("游戏未在进行中")
	}
	
	// 检查是否有活跃的下注轮
	if r.BettingRound == nil {
		return statemachine.ActionResult{}, fmt.Errorf("当前没有下注轮")
	}
	
	// 处理操作
	result := r.BettingRound.ProcessAction(userID, action, amount)
	
	// 更新玩家状态
	if result.Success {
		player.LastAction = action
		player.BetAmount = amount
		
		// 更新玩家状态
		switch action {
		case statemachine.Fold:
			player.Status = PlayerFolded
		case statemachine.AllIn:
			player.Status = PlayerAllIn
		}
		
		// 记录游戏日志
		r.logGameAction(fmt.Sprintf("玩家 %s %s", player.Username, action.String()))
		
		// 如果下注轮结束，触发相应事件
		if result.NextEvent != 0 {
			r.StateMachine.Transition(result.NextEvent)
		}
	}
	
	return result, nil
}

// 私有方法

// findAvailablePosition 找到可用的座位位置
func (r *Room) findAvailablePosition() int {
	occupied := make(map[int]bool)
	for _, player := range r.Players {
		occupied[player.Position] = true
	}
	
	for i := 0; i < r.MaxPlayers; i++ {
		if !occupied[i] {
			return i
		}
	}
	
	return -1
}

// getActivePlayerIDs 获取活跃玩家ID列表
func (r *Room) getActivePlayerIDs() []int64 {
	var playerIDs []int64
	for id, player := range r.Players {
		if player.Status == PlayerSitting || player.Status == PlayerActive {
			playerIDs = append(playerIDs, id)
		}
	}
	return playerIDs
}

// resetRoomState 重置房间状态
func (r *Room) resetRoomState() {
	r.CommunityCards = make([]poker.Card, 0, 5)
	r.Pot = 0
	
	// 重置所有玩家状态
	for _, player := range r.Players {
		player.Status = PlayerActive
		player.Cards = make([]poker.Card, 0, 2)
		player.LastAction = 0
		player.BetAmount = 0
		player.IsDealer = false
		player.IsSmallBlind = false
		player.IsBigBlind = false
	}
}

// setupBlinds 设置盲注
func (r *Room) setupBlinds() {
	playerIDs := r.getActivePlayerIDs()
	if len(playerIDs) < 2 {
		return
	}
	
	// 设置庄家位置（轮转）
	if r.DealerPosition >= len(playerIDs) {
		r.DealerPosition = 0
	}
	
	// 设置庄家、小盲注、大盲注
	dealerID := playerIDs[r.DealerPosition]
	smallBlindID := playerIDs[(r.DealerPosition+1)%len(playerIDs)]
	bigBlindID := playerIDs[(r.DealerPosition+2)%len(playerIDs)]
	
	// 对于只有2个玩家的情况，庄家是小盲注
	if len(playerIDs) == 2 {
		smallBlindID = dealerID
		bigBlindID = playerIDs[(r.DealerPosition+1)%len(playerIDs)]
	}
	
	r.Players[dealerID].IsDealer = true
	r.Players[smallBlindID].IsSmallBlind = true
	r.Players[smallBlindID].BetAmount = r.SmallBlind
	r.Players[bigBlindID].IsBigBlind = true
	r.Players[bigBlindID].BetAmount = r.BigBlind
	
	// 更新底池
	r.Pot = r.SmallBlind + r.BigBlind
}

// startPreFlop 开始发牌阶段
func (r *Room) startPreFlop() error {
	// 给每个活跃玩家发2张底牌
	for _, player := range r.Players {
		if player.Status == PlayerActive {
			player.Cards = []poker.Card{
				r.Deck.Deal(),
				r.Deck.Deal(),
			}
		}
	}
	
	// 创建下注轮
	r.BettingRound = statemachine.NewBettingRound(r.getActivePlayerIDs())
	
	r.logGameAction("开始发牌，每位玩家获得2张底牌")
	return nil
}

// dealFlop 发翻牌（3张公共牌）
func (r *Room) dealFlop() error {
	// 烧一张牌，然后发3张公共牌
	r.Deck.Deal() // 烧牌
	
	for i := 0; i < 3; i++ {
		r.CommunityCards = append(r.CommunityCards, r.Deck.Deal())
	}
	
	// 创建新的下注轮
	r.BettingRound = statemachine.NewBettingRound(r.getActivePlayerIDs())
	
	r.logGameAction("翻牌：发出3张公共牌")
	return nil
}

// dealTurn 发转牌（第4张公共牌）
func (r *Room) dealTurn() error {
	// 烧一张牌，然后发1张公共牌
	r.Deck.Deal() // 烧牌
	r.CommunityCards = append(r.CommunityCards, r.Deck.Deal())
	
	// 创建新的下注轮
	r.BettingRound = statemachine.NewBettingRound(r.getActivePlayerIDs())
	
	r.logGameAction("转牌：发出第4张公共牌")
	return nil
}

// dealRiver 发河牌（第5张公共牌）
func (r *Room) dealRiver() error {
	// 烧一张牌，然后发1张公共牌
	r.Deck.Deal() // 烧牌
	r.CommunityCards = append(r.CommunityCards, r.Deck.Deal())
	
	// 创建新的下注轮
	r.BettingRound = statemachine.NewBettingRound(r.getActivePlayerIDs())
	
	r.logGameAction("河牌：发出第5张公共牌")
	return nil
}

// showdown 摊牌阶段
func (r *Room) showdown() error {
	r.logGameAction("进入摊牌阶段")
	
	// 自动触发确定获胜者事件
	return r.StateMachine.Transition(statemachine.DetermineWinner)
}

// endGame 结束游戏
func (r *Room) endGame() error {
	// 确定获胜者
	winnerID := r.determineWinner()
	
	if winnerID != -1 {
		winner := r.Players[winnerID]
		winner.Chips += r.Pot
		
		// 更新游戏会话
		if r.CurrentGame != nil {
			r.CurrentGame.WinnerID = winnerID
			r.CurrentGame.WinAmount = r.Pot
		}
		
		r.logGameAction(fmt.Sprintf("玩家 %s 获胜，赢得 %d 筹码", winner.Username, r.Pot))
	}
	
	// 移动庄家位置到下一个玩家
	r.DealerPosition = (r.DealerPosition + 1) % len(r.getActivePlayerIDs())
	
	// 重置下注轮
	r.BettingRound = nil
	
	// 设置房间状态为等待
	r.Status = RoomWaiting
	
	r.logGameAction("游戏结束")
	return nil
}

// determineWinner 确定获胜者
func (r *Room) determineWinner() int64 {
	activePlayerIDs := r.getActivePlayerIDs()
	if len(activePlayerIDs) == 0 {
		return -1
	}
	
	// 如果只有一个活跃玩家，直接获胜
	if len(activePlayerIDs) == 1 {
		return activePlayerIDs[0]
	}
	
	// 比较手牌强度
	var bestHand poker.Hand
	var winnerID int64 = -1
	
	for _, playerID := range activePlayerIDs {
		player := r.Players[playerID]
		if player.Status != PlayerFolded {
			// 组合玩家手牌和公共牌
			allCards := append(player.Cards, r.CommunityCards...)
			if len(allCards) == 7 {
				hand := poker.EvaluateHand(allCards)
				
				if winnerID == -1 || poker.CompareHands(hand, bestHand) > 0 {
					bestHand = hand
					winnerID = playerID
				}
			}
		}
	}
	
	return winnerID
}

// logGameAction 记录游戏操作
func (r *Room) logGameAction(action string) {
	if r.CurrentGame != nil {
		timestamp := time.Now().Format("15:04:05")
		logEntry := fmt.Sprintf("[%s] %s", timestamp, action)
		r.CurrentGame.GameLog = append(r.CurrentGame.GameLog, logEntry)
	}
}

// GetRoomInfo 获取房间信息（用于广播）
func (r *Room) GetRoomInfo() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return map[string]interface{}{
		"id":              r.ID,
		"name":            r.Name,
		"status":          r.Status,
		"players":         r.Players,
		"community_cards": r.CommunityCards,
		"pot":             r.Pot,
		"current_state":   r.StateMachine.GetCurrentState().String(),
		"dealer_position": r.DealerPosition,
		"small_blind":     r.SmallBlind,
		"big_blind":       r.BigBlind,
		"updated_at":      r.UpdatedAt,
	}
} 