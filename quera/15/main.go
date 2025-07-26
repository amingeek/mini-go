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

func toLower(s string) string {
	return strings.ToLower(s)
}

func (s *Store) AddProduct(name string, price float64, count int) error {
	normalized := toLower(name)

	for _, existing := range s.ProductsNames {
		if existing == normalized {
			return fmt.Errorf("%s already exists", name) // دقت کن که اسم اصلی ورودی رو باید استفاده کنی
		}
	}

	if price <= 0 {
		return fmt.Errorf("price should be positive")
	}
	if count <= 0 {
		return fmt.Errorf("count should be positive")
	}

	s.ProductsNames = append(s.ProductsNames, normalized)
	s.ProductsPrices = append(s.ProductsPrices, price)
	s.ProductsCount = append(s.ProductsCount, count)

	return nil
}

func (s *Store) GetProductCount(name string) (int, error) {
	normalized := toLower(name)

	for i, n := range s.ProductsNames {
		if n == normalized {
			return s.ProductsCount[i], nil
		}
	}
	return 0, fmt.Errorf("invalid product name")
}

func (s *Store) GetProductPrice(name string) (float64, error) {
	normalized := toLower(name)

	for i, n := range s.ProductsNames {
		if n == normalized {
			return s.ProductsPrices[i], nil
		}
	}
	return 0, fmt.Errorf("invalid product name")
}

func (s *Store) Order(name string, count int) error {
	normalized := toLower(name)

	if count <= 0 {
		return fmt.Errorf("count should be positive")
	}

	for i, n := range s.ProductsNames {
		if n == normalized {
			if s.ProductsCount[i] == 0 {
				return fmt.Errorf("there is no %s in the store", name)
			}
			if count > s.ProductsCount[i] {
				return fmt.Errorf("not enough %s in the store. there are %d left", name, s.ProductsCount[i])
			}
			s.ProductsCount[i] -= count
			return nil
		}
	}
	return fmt.Errorf("invalid product name")
}

func (s *Store) ProductsList() ([]string, error) {
	var available []string

	for i, n := range s.ProductsNames {
		if s.ProductsCount[i] > 0 {
			available = append(available, n)
		}
	}

	if len(available) == 0 {
		return nil, fmt.Errorf("store is empty")
	}

	sort.Strings(available)
	return available, nil
}
