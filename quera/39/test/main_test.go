package test

import (
	"TestEcho/handlers"
	"TestEcho/models"
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateTestDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseTestDatabase(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}

func TestAddProduct(t *testing.T) {
	e := echo.New()

	requestBody := map[string]interface{}{
		"name":  "potato",
		"count": 200,
		"price": 30000,
	}
	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/add-product", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	db, err := CreateTestDatabase()
	assert.NoError(t, err)
	defer CloseTestDatabase(db)

	err = handlers.AddProduct(c, db)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response models.Product
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "potato", response.Name)
	assert.Equal(t, int64(200), response.Count)
	assert.Equal(t, int64(30000), response.Price)
}

func TestListProducts(t *testing.T) {
	e := echo.New()
	db, err := CreateTestDatabase()

	db.Create(&models.Product{Name: "egg", Count: 100, Price: 150})
	db.Create(&models.Product{Name: "tomato", Count: 200, Price: 50})

	assert.NoError(t, err)
	defer CloseTestDatabase(db)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err = handlers.ListProducts(ctx, db)

	assert.NoError(t, err, "Expected no error in ListProducts")
	assert.Equal(t, http.StatusOK, rec.Code, "Expected status code to be 200")

	expectedBody := `[{"Name":"egg","Price":150,"Count":100},{"Name":"tomato","Price":50,"Count":200}]`
	assert.JSONEq(t, expectedBody, rec.Body.String())
}
