package models

import "fmt"

type UtilityArray struct {
	RTWUs []int
	RLUs  []int
	RSUs  []int
	Size  int
}

func NewUtilityArray(size int) *UtilityArray {
	size += 1
	return &UtilityArray{
		RTWUs: make([]int, size),
		RLUs:  make([]int, size),
		RSUs:  make([]int, size),
		Size:  size,
	}
}

func (ua *UtilityArray) SetRTWU(index, value int) {
	if index >= 0 && index < ua.Size {
		ua.RTWUs[index] = value
	}
}

func (ua *UtilityArray) GetRTWU(index int) int {
	if index >= 0 && index < ua.Size {
		return ua.RTWUs[index]
	}
	return 0
}

func (ua *UtilityArray) SetRLU(index, value int) {
	if index >= 0 && index < ua.Size {
		ua.RLUs[index] = value
	}
}

func (ua *UtilityArray) GetRLU(index int) int {
	if index >= 0 && index < ua.Size {
		return ua.RLUs[index]
	}
	return 0
}

func (ua *UtilityArray) SetRSU(index, value int) {
	if index >= 0 && index < ua.Size {
		ua.RSUs[index] = value
	}
}

func (ua *UtilityArray) GetRSU(index int) int {
	if index >= 0 && index < ua.Size {
		return ua.RSUs[index]
	}
	return 0
}

func (ua *UtilityArray) PrintUtilityArray() {

	fmt.Println("RTWU Array:")
	fmt.Println(ua.RTWUs)

	fmt.Println("RLU Array:")
	fmt.Println(ua.RLUs)

	fmt.Println("RSU Array:")
	fmt.Println(ua.RSUs)
}
