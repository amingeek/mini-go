package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type APIResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}

type product struct {
	Name    string `json:"name"`
	Id      int    `json:"id"`
	Price   int    `json:"price"`
	Count   int    `json:"count"`
	Company string `json:"company"`
	Color   string `json:"color"`
}

func getProducts() []product {
	products := []product{
		// 📌 اسم تکراری
		{Name: "Laptop Pro", Id: 1, Price: 30000, Count: 5, Company: "Dell", Color: "Black"},
		{Name: "Laptop Pro", Id: 2, Price: 32000, Count: 3, Company: "Dell", Color: "Silver"},

		// 📌 قیمت یکسان
		{Name: "iPhone 14", Id: 3, Price: 60000, Count: 10, Company: "Apple", Color: "Silver"},
		{Name: "Galaxy S23", Id: 4, Price: 60000, Count: 7, Company: "Samsung", Color: "Gray"},

		// 📌 تعداد مساوی
		{Name: "AirPods Pro", Id: 5, Price: 8000, Count: 20, Company: "Apple", Color: "White"},
		{Name: "Galaxy Buds", Id: 6, Price: 7000, Count: 20, Company: "Samsung", Color: "Black"},

		// 📌 رنگ یکسان
		{Name: "Mechanical Keyboard", Id: 7, Price: 3500, Count: 15, Company: "Corsair", Color: "Black"},
		{Name: "Gaming Mouse", Id: 8, Price: 2500, Count: 25, Company: "Logitech", Color: "Black"},
		{Name: "Monitor 27inch", Id: 9, Price: 12000, Count: 8, Company: "LG", Color: "Black"},

		// 📌 خاص (فقط یکی)
		{Name: "PlayStation 5", Id: 10, Price: 20000, Count: 1, Company: "Sony", Color: "White"},
		{Name: "Xbox Series X", Id: 11, Price: 18000, Count: 1, Company: "Microsoft", Color: "Black"},

		// 📌 خاص (تعداد خیلی زیاد)
		{Name: "USB Cable", Id: 12, Price: 100, Count: 500, Company: "Anker", Color: "Black"},
		{Name: "Phone Case", Id: 13, Price: 300, Count: 1000, Company: "Spigen", Color: "Transparent"},
	}
	return products
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := r.URL.Query()
	name := params.Get("name")
	products := getProducts()
	res := []product{}

	if len(name) == 0 {
		// ❌ پارامتر خالی
		jsonRes, _ := json.Marshal(APIResponse{
			Data:    nil,
			Message: "Bad Request: 'name' parameter is required",
			Status:  "error",
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}

	// 🔍 جستجو (منعطف)
	for _, p := range products {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			res = append(res, p)
		}
	}

	if len(res) == 0 {
		// ❌ محصولی پیدا نشد
		jsonRes, _ := json.Marshal(APIResponse{
			Data:    nil,
			Message: "No product found",
			Status:  "ok",
		})
		w.WriteHeader(http.StatusOK)
		w.Write(jsonRes)
		return
	}

	// ✅ موفقیت
	jsonRes, _ := json.Marshal(APIResponse{
		Data:    res,
		Message: "Products found",
		Status:  "ok",
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func main() {
	http.HandleFunc("/search", search)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
