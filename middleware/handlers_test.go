package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-products-rest/models"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandlers(t *testing.T) {

	newProd := models.ProductRequest{Name: "Test", ShortDescription: "Test", Description: "Test", Price: 1.0, Quantity: 1, Category: models.Category{Id: 1}}
	newProdEnc, _ := json.Marshal(newProd)
	updProd := models.ProductRequest{Name: "Test 2", ShortDescription: "Test 2", Description: "Testing 2", Price: 2.0, Quantity: 2, Category: models.Category{Id: 2}}
	updProdEnc, _ := json.Marshal(updProd)
	newCat := models.CategoryRequest{Name: "Test"}
	updCat := models.CategoryRequest{Name: "Test 2"}
	newCatEnc, _ := json.Marshal(newCat)
	updCatEnc, _ := json.Marshal(updCat)
	var id string

	tests := []models.TestData{
		{Method: "GET", Url: "/api/products", Foo: GetProducts, WantStatus: 200, Body: nil, Vars: nil, HandlerName: "GetProducts"},
		{Method: "GET", Url: "/api/products/{id}", Foo: GetProduct, WantStatus: 200, Body: nil, Vars: map[string]string{"id": "1"}, HandlerName: "GetProduct"},
		{Method: "POST", Url: "/api/products", Foo: CreateProduct, WantStatus: 201, Body: bytes.NewReader(newProdEnc), Vars: nil, HandlerName: "CreateProduct"},
		{Method: "PUT", Url: "/api/products/{id}", Foo: UpdateProduct, WantStatus: 200, Body: bytes.NewReader(updProdEnc), Vars: nil, HandlerName: "UpdateProduct"},
		{Method: "DELETE", Url: "/api/products/{id}", Foo: DeleteProduct, WantStatus: 204, Body: nil, Vars: nil, HandlerName: "DeleteProduct"},
		{Method: "GET", Url: "/api/categories", Foo: GetCategories, WantStatus: 200, Body: nil, Vars: nil, HandlerName: "GetCategories"},
		{Method: "GET", Url: "/api/categories/{id}", Foo: GetCategory, WantStatus: 200, Body: nil, Vars: map[string]string{"id": "1"}, HandlerName: "GetCategory"},
		{Method: "POST", Url: "/api/categories", Foo: CreateCategory, WantStatus: 201, Body: bytes.NewReader(newCatEnc), Vars: nil, HandlerName: "CreateCategory"},
		{Method: "PUT", Url: "/api/categories/{id}", Foo: UpdateCategory, WantStatus: 200, Body: bytes.NewReader(updCatEnc), Vars: nil, HandlerName: "UpdateCategory"},
		{Method: "DELETE", Url: "/api/categories/{id}", Foo: DeleteCategory, WantStatus: 204, Body: nil, Vars: nil, HandlerName: "DeleteCategory"},
	}

	for _, test := range tests {
		if test.Method == "POST" {
			w, err := checkHandler(test)
			id = fetchId(w)
			if err != nil {
				t.Error(err)
			}
			continue
		}
		if test.Method == "DELETE" || test.Method == "PUT" {
			test.Vars = map[string]string{"id": id}
			_, err := checkHandler(test)
			if err != nil {
				t.Error(err)
			}
			continue
		}
		_, err := checkHandler(test)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestUserHandlers(t *testing.T) {
	user := models.LoginRequest{Email: "ab@a.com", Password: "a"}
	userRegEnc, _ := json.Marshal(user)
	// userUpdate := models.UserRequest{FirstName: "Marko", LastName: "Markovic"}
	tests := []models.TestData{
		{Method: "POST", Url: "/api/register", Foo: UserRegister, WantStatus: 200, Body: bytes.NewReader(userRegEnc), Vars: nil, HandlerName: "UserRegister"},
		{Method: "POST", Url: "/api/login", Foo: UserLogin, WantStatus: 200, Body: bytes.NewReader(userRegEnc), Vars: nil, HandlerName: "GetProducts"},
		// {Method: "PUT", Url: "/api/users/{id}", Foo: UpdateUser, WantStatus: 200, Body: bytes.NewReader(newProdEnc), Vars: nil, HandlerName: "CreateProduct"},
	}
	for _, test := range tests {
		_, err := checkHandler(test)
		if err != nil {
			t.Error(err)
		}
	}
}

func fetchId(w *httptest.ResponseRecorder) string {
	re := regexp.MustCompile("[0-9]+")
	bodyNumbers := re.FindAllString(w.Body.String(), -1)
	id := bodyNumbers[0]
	return id
}

func checkHandler(data models.TestData) (*httptest.ResponseRecorder, error) {
	r, err := http.NewRequest(data.Method, data.Url, data.Body)
	if err != nil {
		log.Fatal(err)
	}
	if data.Vars != nil {
		r = mux.SetURLVars(r, data.Vars)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(data.Foo)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != data.WantStatus {
		err = fmt.Errorf("%v status code is %v, wanted status code %v", data.HandlerName, status, data.WantStatus)
	}
	return w, err
}
