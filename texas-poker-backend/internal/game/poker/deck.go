// 牌堆管理系统
// 作用：创建、洗牌、发牌，管理德州扑克的52张标准牌

package poker

import (
	"math/rand"
	"time"
)

// Deck 牌堆结构
type Deck struct {
	cards []Card
	index int // 当前发牌位置
}

// NewDeck 创建新的标准52张牌堆
func NewDeck() *Deck {
	cards := make([]Card, 0, 52)
	
	// 创建52张标准扑克牌
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	
	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, NewCard(rank, suit))
		}
	}
	
	return &Deck{
		cards: cards,
		index: 0,
	}
}

// Shuffle 洗牌
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	
	// Fisher-Yates洗牌算法
	for i := len(d.cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
	
	d.index = 0 // 重置发牌位置
}

// Deal 发一张牌
func (d *Deck) Deal() Card {
	if d.index >= len(d.cards) {
		panic("牌堆已空，无法继续发牌")
	}
	
	card := d.cards[d.index]
	d.index++
	return card
}

// CanDeal 检查是否还能发牌
func (d *Deck) CanDeal() bool {
	return d.index < len(d.cards)
}

// Remaining 剩余牌数
func (d *Deck) Remaining() int {
	return len(d.cards) - d.index
}

// Reset 重置牌堆（重新洗牌）
func (d *Deck) Reset() {
	d.index = 0
	d.Shuffle()
}

// PeekNext 查看下一张牌（不发出）
func (d *Deck) PeekNext() Card {
	if d.index >= len(d.cards) {
		panic("牌堆已空，无法查看下一张牌")
	}
	
	return d.cards[d.index]
}

// GetAllCards 获取所有牌（用于调试）
func (d *Deck) GetAllCards() []Card {
	return d.cards
} 