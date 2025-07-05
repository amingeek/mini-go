package main

func AddElement(numbers *[]int, element int) {
	//TODO
	*numbers = append(*numbers, element)
}

func FindMin(numbers *[]int) int {
	//TODO
	if len(*numbers) == 0 {
		return 0
	}

	min := (*numbers)[0]
	for _, num := range *numbers {
		if num < min {
			min = num
		}
	}

	return min
}

func ReverseSlice(numbers *[]int) {
	if numbers == nil {
		return
	}
	for i := 0; i < len(*numbers)/2; i++ {
		j := len(*numbers) - 1 - i
		(*numbers)[i], (*numbers)[j] = (*numbers)[j], (*numbers)[i]
	}
}

func SwapElements(numbers *[]int, i, j int) {
	if numbers == nil {
		return
	}
	if i < 0 || j < 0 || i >= len(*numbers) || j >= len(*numbers) {
		return
	}
	if i == j {
		return
	}
	(*numbers)[i], (*numbers)[j] = (*numbers)[j], (*numbers)[i]
}
