package main

import "fmt"

func main() {
	db := GetConnection()

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Connected to MySQL database 'quera' successfully!")

	// می‌توانید دیتابیس‌های قابل دسترسی را چک کنید
	var result string
	db.Raw("SELECT DATABASE()").Scan(&result)
	fmt.Printf("Current database: %s\n", result)
}
