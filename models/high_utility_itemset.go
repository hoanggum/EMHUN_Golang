package models

import (
	"fmt"
	"strings"
)

type HighUtilityItemset struct {
	Itemset []int
	Utility int
}

func NewHighUtilityItemset(itemset []int, utility int) *HighUtilityItemset {
	return &HighUtilityItemset{
		Itemset: itemset,
		Utility: utility,
	}
}

func (hui *HighUtilityItemset) GetItemset() []int {
	return hui.Itemset
}

func (hui *HighUtilityItemset) GetUtility() int {
	return hui.Utility
}

func (hui *HighUtilityItemset) String() string {
	return fmt.Sprintf("Itemset: [%s], Utility: %d", strings.Trim(fmt.Sprint(hui.Itemset), "[]"), hui.Utility)
}
