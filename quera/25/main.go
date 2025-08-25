package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
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
	params := r.URL.Query()
	name := params.Get("name")
	products := getProducts()
	res := []product{}

	for _, p := range products {
		if p.Name == name {
			res = append(res, p)
		}
	}

	jsonRes, _ := json.Marshal(Response{})
	if len(res) == 0 {
		jsonRes, _ = json.Marshal(Response{Message: "no product found", Status: "ok"})
	} else {
		jsonRes, _ = json.Marshal(res)
	}
	w.Header().Set("Content-Type", "application/json")
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
