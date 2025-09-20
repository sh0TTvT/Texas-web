// 德州扑克牌型判断算法
// 作用：实现完整的德州扑克牌型识别和比较算法，支持从7张牌中找出最佳5张牌组合

package poker

import (
	"sort"
)

// EvaluateHand 评估7张牌的最佳5张牌组合
func EvaluateHand(cards []Card) Hand {
	if len(cards) != 7 {
		panic("德州扑克必须用7张牌评估（2张底牌+5张公共牌）")
	}

	// 生成所有可能的5张牌组合
	combinations := generateCombinations(cards, 5)
	
	var bestHand Hand
	bestHandValue := -1

	// 评估每种组合，找出最强的手牌
	for _, combination := range combinations {
		hand := evaluateFiveCards(combination)
		handValue := getHandValue(hand)
		
		if handValue > bestHandValue {
			bestHandValue = handValue
			bestHand = hand
		}
	}

	return bestHand
}

// evaluateFiveCards 评估5张牌的牌型
func evaluateFiveCards(cards []Card) Hand {
	// 复制并排序卡牌
	sortedCards := Cards(cards).Clone()
	sortedCards.Sort()

	// 检查各种牌型（按强度从高到低）
	if hand := checkRoyalFlush(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkStraightFlush(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkFourOfAKind(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkFullHouse(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkFlush(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkStraight(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkThreeOfAKind(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkTwoPair(sortedCards); hand.Type != -1 {
		return hand
	}
	if hand := checkOnePair(sortedCards); hand.Type != -1 {
		return hand
	}

	// 高牌
	return checkHighCard(sortedCards)
}

// checkRoyalFlush 检查皇家同花顺（A K Q J 10同花）
func checkRoyalFlush(cards Cards) Hand {
	if hand := checkStraightFlush(cards); hand.Type == StraightFlush {
		// 检查是否为A K Q J 10
		if hand.Ranks[0] == Ace {
			return Hand{
				Cards: cards,
				Type:  RoyalFlush,
				Ranks: []Rank{Ace},
			}
		}
	}
	return Hand{Type: -1}
}

// checkStraightFlush 检查同花顺
func checkStraightFlush(cards Cards) Hand {
	if flushHand := checkFlush(cards); flushHand.Type == Flush {
		if straightHand := checkStraight(cards); straightHand.Type == Straight {
			return Hand{
				Cards: cards,
				Type:  StraightFlush,
				Ranks: straightHand.Ranks,
			}
		}
	}
	return Hand{Type: -1}
}

// checkFourOfAKind 检查四条
func checkFourOfAKind(cards Cards) Hand {
	rankCounts := getRankCounts(cards)
	
	for rank, count := range rankCounts {
		if count == 4 {
			// 找到四条，找出剩余的一张牌
			var kicker Rank
			for r, c := range rankCounts {
				if r != rank && c == 1 {
					kicker = r
					break
				}
			}
			
			return Hand{
				Cards: cards,
				Type:  FourOfAKind,
				Ranks: []Rank{rank, kicker},
			}
		}
	}
	
	return Hand{Type: -1}
}

// checkFullHouse 检查葫芦（三条+一对）
func checkFullHouse(cards Cards) Hand {
	rankCounts := getRankCounts(cards)
	
	var threeRank, pairRank Rank
	var hasThree, hasPair bool
	
	for rank, count := range rankCounts {
		if count == 3 {
			threeRank = rank
			hasThree = true
		} else if count == 2 {
			pairRank = rank
			hasPair = true
		}
	}
	
	if hasThree && hasPair {
		return Hand{
			Cards: cards,
			Type:  FullHouse,
			Ranks: []Rank{threeRank, pairRank},
		}
	}
	
	return Hand{Type: -1}
}

// checkFlush 检查同花
func checkFlush(cards Cards) Hand {
	suitCounts := getSuitCounts(cards)
	
	for _, count := range suitCounts {
		if count == 5 {
			// 按点数降序排列
			ranks := make([]Rank, 5)
			for i, card := range cards {
				ranks[i] = card.Rank
			}
			sort.Slice(ranks, func(i, j int) bool {
				return ranks[i] > ranks[j]
			})
			
			return Hand{
				Cards: cards,
				Type:  Flush,
				Ranks: ranks,
			}
		}
	}
	
	return Hand{Type: -1}
}

// checkStraight 检查顺子
func checkStraight(cards Cards) Hand {
	ranks := make([]Rank, len(cards))
	for i, card := range cards {
		ranks[i] = card.Rank
	}
	
	// 去重并排序
	uniqueRanks := removeDuplicateRanks(ranks)
	sort.Slice(uniqueRanks, func(i, j int) bool {
		return uniqueRanks[i] > uniqueRanks[j]
	})
	
	// 检查是否有连续的5张牌
	if len(uniqueRanks) >= 5 {
		for i := 0; i <= len(uniqueRanks)-5; i++ {
			if isConsecutive(uniqueRanks[i:i+5]) {
				return Hand{
					Cards: cards,
					Type:  Straight,
					Ranks: []Rank{uniqueRanks[i]}, // 最高牌
				}
			}
		}
	}
	
	// 特殊情况：A-2-3-4-5（轮子）
	if containsRanks(uniqueRanks, []Rank{Ace, Five, Four, Three, Two}) {
		return Hand{
			Cards: cards,
			Type:  Straight,
			Ranks: []Rank{Five}, // A-2-3-4-5中，5是最高牌
		}
	}
	
	return Hand{Type: -1}
}

// checkThreeOfAKind 检查三条
func checkThreeOfAKind(cards Cards) Hand {
	rankCounts := getRankCounts(cards)
	
	for rank, count := range rankCounts {
		if count == 3 {
			// 找到三条，找出剩余的两张牌作为kicker
			kickers := []Rank{}
			for r, c := range rankCounts {
				if r != rank && c == 1 {
					kickers = append(kickers, r)
				}
			}
			
			// 按点数降序排列kickers
			sort.Slice(kickers, func(i, j int) bool {
				return kickers[i] > kickers[j]
			})
			
			ranks := []Rank{rank}
			ranks = append(ranks, kickers...)
			
			return Hand{
				Cards: cards,
				Type:  ThreeOfAKind,
				Ranks: ranks,
			}
		}
	}
	
	return Hand{Type: -1}
}

// checkTwoPair 检查两对
func checkTwoPair(cards Cards) Hand {
	rankCounts := getRankCounts(cards)
	
	pairs := []Rank{}
	var kicker Rank
	
	for rank, count := range rankCounts {
		if count == 2 {
			pairs = append(pairs, rank)
		} else if count == 1 {
			kicker = rank
		}
	}
	
	if len(pairs) == 2 {
		// 按点数降序排列对子
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i] > pairs[j]
		})
		
		return Hand{
			Cards: cards,
			Type:  TwoPair,
			Ranks: []Rank{pairs[0], pairs[1], kicker},
		}
	}
	
	return Hand{Type: -1}
}

// checkOnePair 检查一对
func checkOnePair(cards Cards) Hand {
	rankCounts := getRankCounts(cards)
	
	for rank, count := range rankCounts {
		if count == 2 {
			// 找到对子，找出剩余的三张牌作为kicker
			kickers := []Rank{}
			for r, c := range rankCounts {
				if r != rank && c == 1 {
					kickers = append(kickers, r)
				}
			}
			
			// 按点数降序排列kickers
			sort.Slice(kickers, func(i, j int) bool {
				return kickers[i] > kickers[j]
			})
			
			ranks := []Rank{rank}
			ranks = append(ranks, kickers...)
			
			return Hand{
				Cards: cards,
				Type:  OnePair,
				Ranks: ranks,
			}
		}
	}
	
	return Hand{Type: -1}
}

// checkHighCard 高牌
func checkHighCard(cards Cards) Hand {
	ranks := make([]Rank, len(cards))
	for i, card := range cards {
		ranks[i] = card.Rank
	}
	
	// 按点数降序排列
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] > ranks[j]
	})
	
	return Hand{
		Cards: cards,
		Type:  HighCard,
		Ranks: ranks,
	}
}

// 辅助函数

// getRankCounts 获取每个点数的出现次数
func getRankCounts(cards Cards) map[Rank]int {
	counts := make(map[Rank]int)
	for _, card := range cards {
		counts[card.Rank]++
	}
	return counts
}

// getSuitCounts 获取每个花色的出现次数
func getSuitCounts(cards Cards) map[Suit]int {
	counts := make(map[Suit]int)
	for _, card := range cards {
		counts[card.Suit]++
	}
	return counts
}

// removeDuplicateRanks 去除重复的点数
func removeDuplicateRanks(ranks []Rank) []Rank {
	keys := make(map[Rank]bool)
	var result []Rank
	
	for _, rank := range ranks {
		if !keys[rank] {
			keys[rank] = true
			result = append(result, rank)
		}
	}
	
	return result
}

// isConsecutive 检查点数是否连续
func isConsecutive(ranks []Rank) bool {
	for i := 1; i < len(ranks); i++ {
		if ranks[i-1]-ranks[i] != 1 {
			return false
		}
	}
	return true
}

// containsRanks 检查是否包含指定的点数
func containsRanks(ranks []Rank, target []Rank) bool {
	rankMap := make(map[Rank]bool)
	for _, rank := range ranks {
		rankMap[rank] = true
	}
	
	for _, rank := range target {
		if !rankMap[rank] {
			return false
		}
	}
	
	return true
}

// generateCombinations 生成所有可能的组合
func generateCombinations(cards []Card, k int) [][]Card {
	var result [][]Card
	var combination []Card
	
	var backtrack func(start int)
	backtrack = func(start int) {
		if len(combination) == k {
			// 复制当前组合
			combo := make([]Card, k)
			copy(combo, combination)
			result = append(result, combo)
			return
		}
		
		for i := start; i < len(cards); i++ {
			combination = append(combination, cards[i])
			backtrack(i + 1)
			combination = combination[:len(combination)-1]
		}
	}
	
	backtrack(0)
	return result
}

// getHandValue 获取手牌的数值（用于比较）
func getHandValue(hand Hand) int {
	value := int(hand.Type) * 1000000 // 基础牌型分数
	
	// 加上关键点数的分数
	for i, rank := range hand.Ranks {
		value += int(rank) * (100000 / (i + 1)) // 重要性递减
	}
	
	return value
}

// CompareHands 比较两个手牌，返回1表示hand1获胜，-1表示hand2获胜，0表示平局
func CompareHands(hand1, hand2 Hand) int {
	value1 := getHandValue(hand1)
	value2 := getHandValue(hand2)
	
	if value1 > value2 {
		return 1
	} else if value1 < value2 {
		return -1
	} else {
		return 0
	}
} 