// 德州扑克牌型和牌面定义
// 作用：定义扑克牌的基本结构、花色、点数和牌型枚举

package poker

import (
	"fmt"
	"sort"
)

// Suit 花色枚举
type Suit int

const (
	Spades   Suit = 0 // 黑桃
	Hearts   Suit = 1 // 红桃
	Diamonds Suit = 2 // 方块
	Clubs    Suit = 3 // 梅花
)

// String 花色转字符串
func (s Suit) String() string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		return "?"
	}
}

// Rank 牌面点数枚举
type Rank int

const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

// String 牌面点数转字符串
func (r Rank) String() string {
	switch r {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return "?"
	}
}

// Card 扑克牌结构
type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

// String 扑克牌转字符串
func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank.String(), c.Suit.String())
}

// NewCard 创建新扑克牌
func NewCard(rank Rank, suit Suit) Card {
	return Card{Rank: rank, Suit: suit}
}

// ParseCard 从字符串解析扑克牌（如 "AS" = Ace of Spades）
func ParseCard(s string) (Card, error) {
	if len(s) != 2 {
		return Card{}, fmt.Errorf("invalid card string: %s", s)
	}

	var rank Rank
	switch s[0] {
	case '2':
		rank = Two
	case '3':
		rank = Three
	case '4':
		rank = Four
	case '5':
		rank = Five
	case '6':
		rank = Six
	case '7':
		rank = Seven
	case '8':
		rank = Eight
	case '9':
		rank = Nine
	case 'T':
		rank = Ten
	case 'J':
		rank = Jack
	case 'Q':
		rank = Queen
	case 'K':
		rank = King
	case 'A':
		rank = Ace
	default:
		return Card{}, fmt.Errorf("invalid rank: %c", s[0])
	}

	var suit Suit
	switch s[1] {
	case 'S':
		suit = Spades
	case 'H':
		suit = Hearts
	case 'D':
		suit = Diamonds
	case 'C':
		suit = Clubs
	default:
		return Card{}, fmt.Errorf("invalid suit: %c", s[1])
	}

	return NewCard(rank, suit), nil
}

// HandType 牌型枚举（按强度排序）
type HandType int

const (
	HighCard       HandType = 0 // 高牌
	OnePair        HandType = 1 // 一对
	TwoPair        HandType = 2 // 两对
	ThreeOfAKind   HandType = 3 // 三条
	Straight       HandType = 4 // 顺子
	Flush          HandType = 5 // 同花
	FullHouse      HandType = 6 // 葫芦
	FourOfAKind    HandType = 7 // 四条
	StraightFlush  HandType = 8 // 同花顺
	RoyalFlush     HandType = 9 // 皇家同花顺
)

// String 牌型转字符串
func (ht HandType) String() string {
	switch ht {
	case HighCard:
		return "高牌"
	case OnePair:
		return "一对"
	case TwoPair:
		return "两对"
	case ThreeOfAKind:
		return "三条"
	case Straight:
		return "顺子"
	case Flush:
		return "同花"
	case FullHouse:
		return "葫芦"
	case FourOfAKind:
		return "四条"
	case StraightFlush:
		return "同花顺"
	case RoyalFlush:
		return "皇家同花顺"
	default:
		return "未知"
	}
}

// Hand 手牌结构
type Hand struct {
	Cards []Card   `json:"cards"`
	Type  HandType `json:"type"`
	Ranks []Rank   `json:"ranks"` // 用于比较的关键点数（按重要性排序）
}

// String 手牌转字符串
func (h Hand) String() string {
	cardStrs := make([]string, len(h.Cards))
	for i, card := range h.Cards {
		cardStrs[i] = card.String()
	}
	return fmt.Sprintf("%s: %v", h.Type.String(), cardStrs)
}

// Cards 卡牌集合类型
type Cards []Card

// Len 实现sort接口
func (c Cards) Len() int { return len(c) }

// Swap 实现sort接口
func (c Cards) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less 实现sort接口（按点数降序排列）
func (c Cards) Less(i, j int) bool { return c[i].Rank > c[j].Rank }

// Sort 对卡牌按点数降序排序
func (c Cards) Sort() { sort.Sort(c) }

// Clone 复制卡牌集合
func (c Cards) Clone() Cards {
	clone := make(Cards, len(c))
	copy(clone, c)
	return clone
} 