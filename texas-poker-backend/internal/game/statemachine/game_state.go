// 德州扑克游戏状态机
// 作用：管理游戏状态转换，控制游戏流程，处理玩家操作

package statemachine

import (
	"fmt"
	"time"
)

// GameState 游戏状态枚举
type GameState int

const (
	WaitingForPlayers GameState = iota // 等待玩家
	PreFlop                           // 发底牌阶段
	Flop                              // 翻牌阶段（发3张公共牌）
	Turn                              // 转牌阶段（发第4张公共牌）
	River                             // 河牌阶段（发第5张公共牌）
	Showdown                          // 摊牌阶段
	GameEnd                           // 游戏结束
)

// String 游戏状态转字符串
func (gs GameState) String() string {
	switch gs {
	case WaitingForPlayers:
		return "等待玩家"
	case PreFlop:
		return "发牌阶段"
	case Flop:
		return "翻牌阶段"
	case Turn:
		return "转牌阶段"
	case River:
		return "河牌阶段"
	case Showdown:
		return "摊牌阶段"
	case GameEnd:
		return "游戏结束"
	default:
		return "未知状态"
	}
}

// GameEvent 游戏事件枚举
type GameEvent int

const (
	StartGame        GameEvent = iota // 开始游戏
	AllPlayersReady                  // 所有玩家准备完毕
	BettingComplete                  // 下注轮结束
	ShowCards                        // 摊牌
	DetermineWinner                  // 确定获胜者
	NextRound                        // 下一轮
	PlayerLeft                       // 玩家离开
	GameReset                        // 游戏重置
)

// String 游戏事件转字符串
func (ge GameEvent) String() string {
	switch ge {
	case StartGame:
		return "开始游戏"
	case AllPlayersReady:
		return "所有玩家准备完毕"
	case BettingComplete:
		return "下注轮结束"
	case ShowCards:
		return "摊牌"
	case DetermineWinner:
		return "确定获胜者"
	case NextRound:
		return "下一轮"
	case PlayerLeft:
		return "玩家离开"
	case GameReset:
		return "游戏重置"
	default:
		return "未知事件"
	}
}

// StateTransition 状态转换映射
type StateTransition map[GameState]map[GameEvent]GameState

// GameStateMachine 游戏状态机
type GameStateMachine struct {
	currentState GameState
	transitions  StateTransition
	callbacks    map[GameState]func() error // 状态进入回调
}

// NewGameStateMachine 创建新的游戏状态机
func NewGameStateMachine() *GameStateMachine {
	fsm := &GameStateMachine{
		currentState: WaitingForPlayers,
		transitions:  createTransitions(),
		callbacks:    make(map[GameState]func() error),
	}
	
	return fsm
}

// createTransitions 创建状态转换规则
func createTransitions() StateTransition {
	return StateTransition{
		WaitingForPlayers: {
			StartGame: PreFlop,
			GameReset: WaitingForPlayers,
		},
		PreFlop: {
			BettingComplete: Flop,
			PlayerLeft:      GameEnd,
			GameReset:       WaitingForPlayers,
		},
		Flop: {
			BettingComplete: Turn,
			ShowCards:       Showdown,
			PlayerLeft:      GameEnd,
			GameReset:       WaitingForPlayers,
		},
		Turn: {
			BettingComplete: River,
			ShowCards:       Showdown,
			PlayerLeft:      GameEnd,
			GameReset:       WaitingForPlayers,
		},
		River: {
			BettingComplete: Showdown,
			ShowCards:       Showdown,
			PlayerLeft:      GameEnd,
			GameReset:       WaitingForPlayers,
		},
		Showdown: {
			DetermineWinner: GameEnd,
			GameReset:       WaitingForPlayers,
		},
		GameEnd: {
			NextRound: WaitingForPlayers,
			GameReset: WaitingForPlayers,
		},
	}
}

// GetCurrentState 获取当前状态
func (fsm *GameStateMachine) GetCurrentState() GameState {
	return fsm.currentState
}

// CanTransition 检查是否可以进行指定的状态转换
func (fsm *GameStateMachine) CanTransition(event GameEvent) bool {
	if transitions, exists := fsm.transitions[fsm.currentState]; exists {
		_, canTransition := transitions[event]
		return canTransition
	}
	return false
}

// Transition 执行状态转换
func (fsm *GameStateMachine) Transition(event GameEvent) error {
	if !fsm.CanTransition(event) {
		return fmt.Errorf("无法从状态 %s 通过事件 %s 进行转换", 
			fsm.currentState.String(), event.String())
	}
	
	oldState := fsm.currentState
	fsm.currentState = fsm.transitions[oldState][event]
	
	// 执行状态进入回调
	if callback, exists := fsm.callbacks[fsm.currentState]; exists {
		if err := callback(); err != nil {
			// 如果回调失败，回滚状态
			fsm.currentState = oldState
			return fmt.Errorf("状态转换回调失败: %w", err)
		}
	}
	
	return nil
}

// SetStateCallback 设置状态进入回调
func (fsm *GameStateMachine) SetStateCallback(state GameState, callback func() error) {
	fsm.callbacks[state] = callback
}

// Reset 重置状态机到初始状态
func (fsm *GameStateMachine) Reset() {
	fsm.currentState = WaitingForPlayers
}

// GetValidEvents 获取当前状态下的有效事件
func (fsm *GameStateMachine) GetValidEvents() []GameEvent {
	var events []GameEvent
	
	if transitions, exists := fsm.transitions[fsm.currentState]; exists {
		for event := range transitions {
			events = append(events, event)
		}
	}
	
	return events
}

// PlayerAction 玩家操作类型
type PlayerAction int

const (
	Fold  PlayerAction = iota // 弃牌
	Call                      // 跟注
	Raise                     // 加注
	Check                     // 过牌
	Bet                       // 下注
	AllIn                     // 全押
)

// String 玩家操作转字符串
func (pa PlayerAction) String() string {
	switch pa {
	case Fold:
		return "弃牌"
	case Call:
		return "跟注"
	case Raise:
		return "加注"
	case Check:
		return "过牌"
	case Bet:
		return "下注"
	case AllIn:
		return "全押"
	default:
		return "未知操作"
	}
}

// ActionRequest 玩家操作请求
type ActionRequest struct {
	PlayerID int64        `json:"player_id"`
	Action   PlayerAction `json:"action"`
	Amount   int          `json:"amount"`
	Timestamp time.Time   `json:"timestamp"`
}

// ActionResult 操作结果
type ActionResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	NextEvent GameEvent `json:"next_event,omitempty"`
}

// BettingRound 下注轮管理
type BettingRound struct {
	currentBet    int               // 当前最高下注
	players       []int64           // 参与下注的玩家ID列表
	currentPlayer int               // 当前轮到的玩家索引
	playerBets    map[int64]int     // 每个玩家的下注金额
	playerActions map[int64]PlayerAction // 每个玩家的最后操作
	completed     bool              // 下注轮是否完成
}

// NewBettingRound 创建新的下注轮
func NewBettingRound(players []int64) *BettingRound {
	return &BettingRound{
		currentBet:    0,
		players:       players,
		currentPlayer: 0,
		playerBets:    make(map[int64]int),
		playerActions: make(map[int64]PlayerAction),
		completed:     false,
	}
}

// GetCurrentPlayer 获取当前应该操作的玩家
func (br *BettingRound) GetCurrentPlayer() int64 {
	if br.currentPlayer >= len(br.players) {
		return -1 // 无效玩家
	}
	return br.players[br.currentPlayer]
}

// ProcessAction 处理玩家操作
func (br *BettingRound) ProcessAction(playerID int64, action PlayerAction, amount int) ActionResult {
	// 验证是否为当前玩家
	if playerID != br.GetCurrentPlayer() {
		return ActionResult{
			Success: false,
			Message: "不是您的操作轮次",
		}
	}
	
	// 处理不同的操作
	switch action {
	case Fold:
		// 移除弃牌玩家
		br.removePlayer(playerID)
		br.playerActions[playerID] = Fold
		
	case Call:
		// 跟注到当前最高金额
		br.playerBets[playerID] = br.currentBet
		br.playerActions[playerID] = Call
		br.nextPlayer()
		
	case Raise:
		// 加注
		if amount <= br.currentBet {
			return ActionResult{
				Success: false,
				Message: "加注金额必须大于当前下注",
			}
		}
		br.currentBet = amount
		br.playerBets[playerID] = amount
		br.playerActions[playerID] = Raise
		br.nextPlayer()
		
	case Check:
		// 过牌（只有在没有人下注时才能过牌）
		if br.currentBet > 0 {
			return ActionResult{
				Success: false,
				Message: "已有人下注，无法过牌",
			}
		}
		br.playerActions[playerID] = Check
		br.nextPlayer()
		
	case Bet:
		// 下注（只有在没有人下注时才能下注）
		if br.currentBet > 0 {
			return ActionResult{
				Success: false,
				Message: "已有人下注，请选择跟注或加注",
			}
		}
		br.currentBet = amount
		br.playerBets[playerID] = amount
		br.playerActions[playerID] = Bet
		br.nextPlayer()
		
	case AllIn:
		// 全押
		br.currentBet = amount // 这里应该是玩家的全部筹码
		br.playerBets[playerID] = amount
		br.playerActions[playerID] = AllIn
		br.nextPlayer()
		
	default:
		return ActionResult{
			Success: false,
			Message: "无效的操作类型",
		}
	}
	
	// 检查下注轮是否完成
	if br.isBettingComplete() {
		br.completed = true
		return ActionResult{
			Success:   true,
			Message:   fmt.Sprintf("玩家 %d %s 成功，下注轮结束", playerID, action.String()),
			NextEvent: BettingComplete,
		}
	}
	
	return ActionResult{
		Success: true,
		Message: fmt.Sprintf("玩家 %d %s 成功", playerID, action.String()),
	}
}

// removePlayer 从下注轮中移除玩家
func (br *BettingRound) removePlayer(playerID int64) {
	for i, id := range br.players {
		if id == playerID {
			br.players = append(br.players[:i], br.players[i+1:]...)
			if br.currentPlayer > i {
				br.currentPlayer--
			}
			break
		}
	}
	
	// 如果当前玩家索引超出范围，回到开始
	if br.currentPlayer >= len(br.players) {
		br.currentPlayer = 0
	}
}

// nextPlayer 切换到下一个玩家
func (br *BettingRound) nextPlayer() {
	br.currentPlayer = (br.currentPlayer + 1) % len(br.players)
}

// isBettingComplete 检查下注轮是否完成
func (br *BettingRound) isBettingComplete() bool {
	if len(br.players) <= 1 {
		return true // 只剩一个或没有玩家
	}
	
	// 检查所有活跃玩家是否都已经操作并且下注金额一致
	activePlayers := 0
	for _, playerID := range br.players {
		if action, exists := br.playerActions[playerID]; exists && action != Fold {
			activePlayers++
			// 检查是否所有玩家的下注都等于当前最高下注
			if br.playerBets[playerID] != br.currentBet {
				return false
			}
		}
	}
	
	return activePlayers > 0
}

// IsCompleted 检查下注轮是否完成
func (br *BettingRound) IsCompleted() bool {
	return br.completed
}

// GetPlayerBets 获取所有玩家的下注
func (br *BettingRound) GetPlayerBets() map[int64]int {
	return br.playerBets
}

// GetCurrentBet 获取当前最高下注
func (br *BettingRound) GetCurrentBet() int {
	return br.currentBet
} 