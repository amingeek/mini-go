package main

import (
	"fmt"
	"sort"
	"strings"
)

type Store struct {
	ProductsNames  []string
	ProductsPrices []float64
	ProductsCount  []int
}

func NewStore() *Store {
	return &Store{
		ProductsNames:  []string{},
		ProductsPrices: []float64{},
		ProductsCount:  []int{},
	}
}

func (s *Store) SumCount() int {
	x := 0
	for _, n := range s.ProductsCount {
		x += n
	}

	return x
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func (s *Store) AddProduct(name string, price float64, count int) error {
	normalizedName := ToLower(name)
	for _, n := range s.ProductsNames {
		if n == normalizedName {
			return fmt.Errorf("product '%s' already exists", normalizedName)
		}
	}
	if price <= 0 {
		return fmt.Errorf("price should be positive")
	}
	if count <= 0 {
		return fmt.Errorf("count should be positive")
	}
	s.ProductsNames = append(s.ProductsNames, normalizedName)
	s.ProductsPrices = append(s.ProductsPrices, float64(price))
	s.ProductsCount = append(s.ProductsCount, count)

	return nil
}

func (s *Store) GetProductCount(name string) (int, error) {
	normalizedName := ToLower(name)

	for i, n := range s.ProductsNames {
		if n == normalizedName {
			return s.ProductsCount[i], nil
		}
	}
	return 0, fmt.Errorf("invalid product name")

}

func (s *Store) GetProductPrice(name string) (float64, error) {
	normalizedName := ToLower(name)

	for i, n := range s.ProductsNames {
		if n == normalizedName {
			return float64(s.ProductsPrices[i]), nil
		}
	}
	return 0, fmt.Errorf("invalid product name")
}

func (s *Store) Order(name string, count int) error {
	normalizedName := ToLower(name)

	if count <= 0 {
		return fmt.Errorf("count should be positive")
	}
	for i, n := range s.ProductsNames {
		if n == normalizedName {
			if s.ProductsCount[i] == 0 {
				return fmt.Errorf("there is no %s in the store ", normalizedName)
			}
			if count > s.ProductsCount[i] {
				return fmt.Errorf("not enough %s in the store. there are %s left", normalizedName, s.ProductsCount[i])
			}
			s.ProductsCount[i] -= count
			return nil
		}
	}
	return fmt.Errorf("invalid product name")

}

func (s *Store) ProductsList() ([]string, error) {
	if len(s.ProductsNames) == 0 {
		return nil, fmt.Errorf("store is empty")
	}
	if s.SumCount() >= 0 {
		sliceNames := make([]string, len(s.ProductsNames))
		for i, n := range s.ProductsNames {
			if s.ProductsCount[i] > 0 {
				sliceNames = append(sliceNames, ToLower(n))
			}
		}
		sort.Strings(sliceNames)
		return sliceNames, nil

	}
	return nil, nil
}
