package main

import (
	"bufio"
	"emhun/algorithms"
	"emhun/models"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileName := "data/table3.txt"
	minUtility := 25

	transactions, err := readTransactionsFromFile(fileName)
	if err != nil {
		fmt.Println("Error reading transactions:", err)
		return
	}
	fmt.Println("Transactions vừa đọc được:")
	for i, transaction := range transactions {
		fmt.Printf("Transaction %d: %s\n", i+1, transaction)
	}
	emhun := algorithms.NewEMHUN(transactions, minUtility)

	emhun.Run()

	fmt.Println("\nFinished executing EMHUN algorithm.")
	outputFileName := "output/results.txt"
	err = writeResultsToFile(emhun, outputFileName)
	if err != nil {
		fmt.Println("Error writing results:", err)
		return
	}

	fmt.Println("Finished executing EMHUN algorithm. Results written to", outputFileName)
}

func readTransactionsFromFile(fileName string) ([]*models.Transaction, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var transactions []*models.Transaction

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		itemsStr := strings.Fields(parts[0])
		var items []int
		for _, item := range itemsStr {
			itemInt, err := strconv.Atoi(item)
			if err != nil {
				return nil, err
			}
			items = append(items, itemInt)
		}

		transUtility, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, err
		}

		utilitiesStr := strings.Fields(parts[2])
		var utilities []int
		for _, utility := range utilitiesStr {
			utilityInt, err := strconv.Atoi(utility)
			if err != nil {
				return nil, err
			}
			utilities = append(utilities, utilityInt)
		}

		transaction := models.NewTransaction(items, utilities, transUtility)
		transactions = append(transactions, transaction)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
func writeResultsToFile(emhun *algorithms.EMHUN, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, hui := range emhun.SearchAlgorithms.HighUtilityItemsets {
		line := fmt.Sprintf("Itemset: %v, Utility: %d\n", hui.Itemset, hui.Utility)
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
