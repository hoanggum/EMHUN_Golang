package utility

import (
	"emhun/models"
	"fmt"
)

func CalculateTransactionUtility(transaction *models.Transaction) int {
	totalUtility := 0
	for _, utility := range transaction.Utilities {
		totalUtility += utility
	}
	return totalUtility
}

func CalculateAndPrintAllTransactionUtilities(transactions []*models.Transaction) {
	for i, transaction := range transactions {
		tu := CalculateTransactionUtility(transaction)
		fmt.Printf("Transaction %d TU: %d\n", i+1, tu)
	}
}
func CalculateRLUForAllItemsRhoAnDenta(transactions []*models.Transaction, rho, delta map[int]bool, utilityArray *models.UtilityArray) {
	combinedSet := UnionMaps(rho, delta)

	for item := range combinedSet {
		totalRLU := 0
		fmt.Printf("\nCalculating RLU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsItem(transaction, item) {
				fmt.Printf("  Found item %d in transaction: %v\n", item, transaction.Items)
				rlu := CalculateRemainingResidualUtility(transaction, item)
				totalRLU += rlu
				fmt.Printf("  RLU for this transaction: %d (cumulative RLU: %d)\n", rlu, totalRLU)
			}
		}

		utilityArray.SetRLU(item, totalRLU)
		fmt.Printf("Calculated total RLU for item %d: %d\n", item, totalRLU)
	}
}
func CalculateRLUForAllItems(transactions []*models.Transaction, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRLU := 0
		fmt.Printf("\nCalculating RLU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsItem(transaction, item) {
				index := GetItemIndex(transaction, item)
				itemUtility := transaction.Utilities[index]
				remainingUtility := CalculateRemainingUtility(transaction, index+1)
				totalRLU += itemUtility + remainingUtility

				fmt.Printf("  Found item %d in transaction %v with utility: %d, Remaining Residual Utility: %d\n",
					item, transaction.Items, itemUtility, remainingUtility)
			}
		}

		utilityArray.SetRLU(item, totalRLU)
		fmt.Printf("Calculated total RLU for item %d: %d\n", item, totalRLU)
	}
}
func CalculateRemainingResidualUtility(transaction *models.Transaction, currentItem int) int {
	foundCurrentItem := false
	rru := 0
	fmt.Printf("    Remaining items after %d: ", currentItem)

	for i, item := range transaction.Items {
		utility := transaction.Utilities[i]

		if foundCurrentItem && utility > 0 {
			rru += utility
			fmt.Printf("%d(%d) ", item, utility)
		}

		if item == currentItem {
			foundCurrentItem = true
			if utility > 0 {
				rru += utility
				fmt.Printf("    Adding utility of currentItem %d: %d\n", currentItem, utility)
			}
		}
	}
	fmt.Println()
	return rru
}
func CalculateRTWUForAllItems(transactions []*models.Transaction, rho, delta, eta map[int]bool, utilityArray *models.UtilityArray) {
	combinedSet := UnionMaps(rho, delta)
	combinedSet = UnionMaps(combinedSet, eta)

	for item := range combinedSet {
		totalRTWU := 0
		// fmt.Printf("\nCalculating RTWU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsItem(transaction, item) {
				// fmt.Printf("  Found item %d in transaction: %v\n", item, transaction.Items)
				rtwu := CalculateRTUForTransaction(transaction)
				totalRTWU += rtwu
				// fmt.Printf(" = %d (cumulative RTWU: %d)\n", rtwu, totalRTWU)
			}
		}

		utilityArray.SetRTWU(item, totalRTWU)
		// fmt.Printf("Calculated total RTWU for item %d: %d\n", item, totalRTWU)
	}
	utilityArray.PrintUtilityArray()
}
func CalculateRTUForTransaction(transaction *models.Transaction) int {
	rtwu := 0
	for _, utility := range transaction.Utilities {
		if utility > 0 {
			rtwu += utility
		}
	}
	return rtwu
}
func CalculateRSUForAllItems(transactions []*models.Transaction, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRSU := 0
		// fmt.Printf("\nCalculating RSU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsItem(transaction, item) {
				index := GetItemIndex(transaction, item)
				itemUtility := transaction.Utilities[index]
				remainingUtility := CalculateRemainingUtility(transaction, index+1)
				totalRSU += itemUtility + remainingUtility

				// fmt.Printf("  Found item %d in transaction %v with utility: %d, Remaining Residual Utility: %d\n",
				// 	item, transaction.Items, itemUtility, remainingUtility)
			}
		}

		utilityArray.SetRSU(item, totalRSU)
		// fmt.Printf("Calculated total RSU for item %d: %d\n", item, totalRSU)
	}
}

// Tính Remaining Utility
func CalculateRemainingUtility(transaction *models.Transaction, startIndex int) int {
	remainingUtility := 0
	for i := startIndex; i < len(transaction.Items); i++ {
		if transaction.Utilities[i] > 0 {
			remainingUtility += transaction.Utilities[i]
		}
	}
	return remainingUtility
}

// Tính RSU cho tất cả các item trong tập X và secondary
func CalculateRSUForAllItem(transactions []*models.Transaction, X []int, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRSU := 0
		fmt.Printf("\nCalculating RSU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsAllItems(transaction, X) && ContainsItem(transaction, item) {
				utilityX := CalculateUtilityForSet(transaction, X)
				indexZ := GetItemIndex(transaction, item)
				utilityZ := transaction.Utilities[indexZ]
				rru := CalculateRemainingUtility(transaction, indexZ+1)
				totalRSU += utilityX + utilityZ + rru

				fmt.Printf("  Found set X %v and item %d in transaction %v with utility of X: %d, utility of z: %d, Remaining Residual Utility: %d, Calculated RSU: %d\n",
					X, item, transaction.Items, utilityX, utilityZ, rru, utilityX+utilityZ+rru)
			}
		}

		utilityArray.SetRSU(item, totalRSU)
		fmt.Printf("Calculated total RSU for item %d: %d\n", item, totalRSU)
	}
}
func CalculateRLUForAllItem(transactions []*models.Transaction, X []int, secondary []int, utilityArray *models.UtilityArray) {
	for _, item := range secondary {
		totalRLU := 0
		fmt.Printf("\nCalculating RLU for item: %d\n", item)

		for _, transaction := range transactions {
			if ContainsAllItems(transaction, X) && ContainsItem(transaction, item) {
				utilityX := CalculateUtilityForSet(transaction, X)
				maxIndexX := FindLocationMaxIndexForSet(transaction, X)
				index := GetItemIndex(transaction, maxIndexX)

				remainingUtility := CalculateRemainingUtility(transaction, index+1)

				totalRLU += utilityX + remainingUtility

				fmt.Printf("  Found item %d in transaction %v with utility of X: %d, Remaining Residual Utility (RRU): %d, Calculated RLU: %d\n",
					item, transaction.Items, utilityX, remainingUtility, utilityX+remainingUtility)
			}
		}

		utilityArray.SetRLU(item, totalRLU)
		fmt.Printf("Calculated total RLU for item %d: %d\n", item, totalRLU)
	}
}

// Tính utility cho một tập hợp item
func CalculateUtilityForSet(transaction *models.Transaction, X []int) int {
	totalUtility := 0
	for _, item := range X {
		if ContainsItem(transaction, item) {
			index := GetItemIndex(transaction, item)
			totalUtility += transaction.Utilities[index]
		}
	}
	return totalUtility
}

// Tìm chỉ số lớn nhất trong tập X
func FindLocationMaxIndexForSet(transaction *models.Transaction, X []int) int {
	maxIndex := -1
	for _, item := range X {
		index := GetItemIndex(transaction, item)
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex
}

// Kiểm tra xem một transaction có chứa item không
func ContainsItem(transaction *models.Transaction, item int) bool {
	for _, tItem := range transaction.Items {
		if tItem == item {
			return true
		}
	}
	return false
}

// Kiểm tra xem transaction có chứa tất cả các item trong tập X không
func ContainsAllItems(transaction *models.Transaction, X []int) bool {
	for _, item := range X {
		if !ContainsItem(transaction, item) {
			return false
		}
	}
	return true
}

// Lấy chỉ số của item trong transaction
func GetItemIndex(transaction *models.Transaction, item int) int {
	for i, tItem := range transaction.Items {
		if tItem == item {
			return i
		}
	}
	return -1
}

// Hợp nhất hai tập hợp lại với nhau (rho và delta)
func UnionMaps(a, b map[int]bool) map[int]bool {
	result := make(map[int]bool)
	for k := range a {
		result[k] = true
	}
	for k := range b {
		result[k] = true
	}
	return result
}
