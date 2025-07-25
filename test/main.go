package main

import (
	"fmt"
	"sort"
)

func main() {
	names := []string{"zahra", "ali", "reza", "mohammad"}

	// مرتب‌سازی حساس به بزرگی/کوچکی حروف
	sort.Strings(names)
	fmt.Println("Sorted (case-sensitive):", names)
	// خروجی: [Reza Zahra ali mohammad]
}
