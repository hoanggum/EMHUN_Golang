package algorithms

import (
	"emhun/models"
	"emhun/utility"
	"fmt"
)

type SearchAlgorithms struct {
	UtilityArray        *models.UtilityArray
	Beta                map[int]bool
	ItemList            []int
	FilteredPrimary     []int
	FilteredSecondary   []int
	HighUtilityItemsets []*models.HighUtilityItemset
}

func NewSearchAlgorithms(utilityArray *models.UtilityArray) *SearchAlgorithms {
	return &SearchAlgorithms{
		UtilityArray:        utilityArray,
		Beta:                make(map[int]bool),
		HighUtilityItemsets: []*models.HighUtilityItemset{},
	}
}

func (s *SearchAlgorithms) Search(eta []int, X map[int]bool, transactions []*models.Transaction, primary []int, secondary []int, minU int) {
	if len(primary) == 0 {
		return
	}

	for _, item := range primary {
		s.Beta = copyMap(X)
		s.Beta[item] = true
		s.ItemList = mapKeys(s.Beta)

		utilityBeta := s.calculateUtility(transactions, s.Beta)
		fmt.Printf("Utility of %v: %d\n", s.Beta, utilityBeta)

		projectedDB := s.projectDatabase(transactions, s.ItemList)
		s.printProjectedDatabase(projectedDB, item)

		if utilityBeta >= minU {
			fmt.Printf("U(%d) = %d >= %d HUI Found: %v\n", item, utilityBeta, minU, s.Beta)
			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(s.ItemList, utilityBeta))
		} else {
			fmt.Printf("%d < %d so %d is not a HUI.\n", utilityBeta, minU, item)
		}

		if utilityBeta > minU {
			s.SearchN(eta, s.Beta, transactions, minU)
		}

		s.FilteredPrimary = []int{}
		s.FilteredSecondary = []int{}
		utility.CalculateRSUForAllItem(transactions, s.ItemList, secondary, s.UtilityArray)
		utility.CalculateRLUForAllItem(transactions, s.ItemList, secondary, s.UtilityArray)

		for _, secItem := range secondary {
			rsu := s.UtilityArray.GetRSU(secItem)
			rlu := s.UtilityArray.GetRLU(secItem)

			if rsu >= minU {
				s.FilteredPrimary = append(s.FilteredPrimary, secItem)
			}
			if rlu >= minU {
				s.FilteredSecondary = append(s.FilteredSecondary, secItem)
			}
		}

		fmt.Printf("Primary%v = %v\n", s.ItemList, s.FilteredPrimary)
		fmt.Printf("Secondary%v = %v\n", s.ItemList, s.FilteredSecondary)
		s.processSecondary(s.FilteredSecondary, s.ItemList, transactions, minU)
		s.Search(eta, s.Beta, projectedDB, s.FilteredPrimary, s.FilteredSecondary, minU)
	}
}

func (s *SearchAlgorithms) processSecondary(secondary []int, beta []int, transactions []*models.Transaction, minU int) {
	for i := 0; i < len(secondary); i++ {
		secItem := secondary[i]
		betaNew := copyMapFromList(beta)
		betaNew[secItem] = true

		utilityBetaNew := s.calculateUtility(transactions, betaNew)
		fmt.Printf("Utility of combination %v: %d\n", betaNew, utilityBetaNew)

		if utilityBetaNew >= minU {
			fmt.Printf("U(%d) = %d >= %d HUI Found: %v\n", secItem, utilityBetaNew, minU, betaNew)
			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(mapKeys(betaNew), utilityBetaNew))
		} else {
			fmt.Printf("%d < %d so %d is not a HUI.\n", utilityBetaNew, minU, secItem)
		}

		for j := i + 1; j < len(secondary); j++ {
			nextSecItem := secondary[j]
			betaExtended := copyMap(betaNew)
			betaExtended[nextSecItem] = true

			utilityBetaExtended := s.calculateUtility(transactions, betaExtended)
			fmt.Printf("Utility of extended combination %v: %d\n", betaExtended, utilityBetaExtended)

			if utilityBetaExtended >= minU {
				fmt.Printf("U(%d) = %d >= %d HUI Found: %v\n", nextSecItem, utilityBetaExtended, minU, betaExtended)
				s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(mapKeys(betaExtended), utilityBetaExtended))
			} else {
				fmt.Printf("%d < %d so %d is not a HUI.\n", utilityBetaExtended, minU, nextSecItem)
			}
		}
	}
}

func (s *SearchAlgorithms) SearchN(eta []int, beta map[int]bool, transactions []*models.Transaction, minU int) {
	if len(eta) == 0 {
		return
	}

	for _, item := range eta {
		betaNew := copyMap(beta)
		betaNew[item] = true

		itemList := mapKeys(betaNew)
		projectedDB := s.projectDatabase(transactions, itemList)
		s.printProjectedDatabase(projectedDB, item)

		utilityBetaNew := s.calculateUtility(transactions, betaNew)
		fmt.Printf("Utility of (negative) %v: %d\n", betaNew, utilityBetaNew)

		if utilityBetaNew >= minU {
			fmt.Printf("U(%d) = %d >= %d HUI Found: %v\n", item, utilityBetaNew, minU, betaNew)
			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(mapKeys(betaNew), utilityBetaNew))
		} else {
			fmt.Printf("%d < %d so %v is not a HUI.\n", utilityBetaNew, minU, betaNew)
		}

		filteredPrimary := []int{}
		utility.CalculateRSUForAllItem(transactions, itemList, eta, s.UtilityArray)
		for _, secItem := range eta {
			rsu := s.UtilityArray.GetRSU(secItem)
			if rsu >= minU {
				filteredPrimary = append(filteredPrimary, secItem)
			}
		}

		remainingEta := removeItem(eta, item)
		s.SearchN(remainingEta, betaNew, transactions, minU)
	}
}

func (s *SearchAlgorithms) projectDatabase(transactions []*models.Transaction, items []int) []*models.Transaction {
	var projectedDB []*models.Transaction

	for _, transaction := range transactions {
		// Kiểm tra nếu giao dịch chứa tất cả các items cần thiết
		if containsAllItems(transaction.Items, items) {
			var projectedItems []int
			var projectedUtilities []int
			lastItemIndex := -1

			// Tìm vị trí của item cuối cùng trong danh sách items
			for _, item := range items {
				itemIndex := indexOf(transaction.Items, item)
				if itemIndex > lastItemIndex {
					lastItemIndex = itemIndex
				}
			}

			// Lấy các item và utilities sau item cuối cùng
			for i := lastItemIndex + 1; i < len(transaction.Items); i++ {
				projectedItems = append(projectedItems, transaction.Items[i])
				projectedUtilities = append(projectedUtilities, transaction.Utilities[i])
			}

			// Nếu có các item sau item cuối cùng, thêm giao dịch đã được chiếu vào kết quả
			if len(projectedItems) > 0 {
				projectedDB = append(projectedDB, models.NewTransaction(projectedItems, projectedUtilities, calculateTransactionUtility(projectedUtilities)))
			}
		}
	}

	return projectedDB
}

func (s *SearchAlgorithms) calculateUtility(transactions []*models.Transaction, itemset map[int]bool) int {
	totalUtility := 0
	for _, transaction := range transactions {
		if containsAllItemsMap(transaction.Items, itemset) {
			itemsetUtility := 0
			for item := range itemset {
				index := indexOf(transaction.Items, item)
				if index != -1 {
					itemsetUtility += transaction.Utilities[index]
				}
			}
			totalUtility += itemsetUtility
		}
	}
	return totalUtility
}

// Helper functions for copying, checking, and map handling
func copyMap(original map[int]bool) map[int]bool {
	copy := make(map[int]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

func mapKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func containsAllItems(items []int, itemset []int) bool {
	for _, item := range itemset {
		if indexOf(items, item) == -1 {
			return false
		}
	}
	return true
}

func indexOf(items []int, item int) int {
	for i, v := range items {
		if v == item {
			return i
		}
	}
	return -1
}

func removeItem(slice []int, item int) []int {
	result := []int{}
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

func containsAllItemsMap(items []int, itemset map[int]bool) bool {
	for item := range itemset {
		if indexOf(items, item) == -1 {
			return false
		}
	}
	return true
}
func copyMapFromList(list []int) map[int]bool {
	copy := make(map[int]bool)
	for _, v := range list {
		copy[v] = true
	}
	return copy
}
func (s *SearchAlgorithms) printProjectedDatabase(projectedDB []*models.Transaction, item int) {
	fmt.Printf("\nProjected Database after item %d:\n", item)
	for _, transaction := range projectedDB {
		fmt.Printf("Items: %v, Utilities: %v, Transaction Utility: %d\n",
			transaction.Items, transaction.Utilities, calculateTransactionUtility(transaction.Utilities))
	}
	fmt.Println("----------------------------------")
}
func calculateTransactionUtility(utilities []int) int {
	totalUtility := 0
	for _, utility := range utilities {
		totalUtility += utility
	}
	return totalUtility
}
