package middleware

import (
	"bytes"
	"encoding/json"
	"go-products-rest/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetProducts(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProducts)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("ERROR - Status code is %v, wanted status code is %v", status, http.StatusOK)
	}
}

func TestGetProduct(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/products/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProduct)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("ERROR - Status code is %v, wanted status code is %v", status, http.StatusOK)
	}
}
func TestCreateProducts(t *testing.T) {
	testP := models.ProductRequest{Name: "Test", ShortDescription: "Test", Description: "Test", Price: 1.0, Quantity: 1, Category: models.Category{Id: 1}}
	marshallP, _ := json.Marshal(testP)
	r, err := http.NewRequest("POST", "/api/products", bytes.NewReader(marshallP))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("ERROR - Status code is %v, wanted status code is %v", status, http.StatusOK)
	}
}
func TestDeleteProducts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(DeleteProduct))
	r, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if r.StatusCode != http.StatusOK {
		t.Errorf("expected 204, got %v", r.StatusCode)
	}

}
func TestUpdateProducts(t *testing.T) {
	testP := models.ProductRequest{Name: "Test", ShortDescription: "Test", Description: "Test", Price: 1.0, Quantity: 1, Category: models.Category{Id: 1}}
	marshallP, _ := json.Marshal(testP)
	r, err := http.NewRequest("PUT", "/api/products/{id}", bytes.NewReader(marshallP))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	handler := http.HandlerFunc(UpdateProduct)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("ERROR - Status code is %v, wanted status code is %v", status, http.StatusOK)
	}
}
func TestGetCategories(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCategories)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("ERROR - Status code is %v, wanted status code is %v", status, http.StatusOK)
	}
}
func TestGetCategory(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/categories/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	handler := http.HandlerFunc(GetCategory)
	handler.ServeHTTP(w, r)
	if status := w.Code; status != http.StatusOK {
		t.Errorf("ERROR - Status code is %v, wanted status code is %v", status, http.StatusOK)
	}
}

func TestCreateCategory(t *testing.T) {

}
func TestDeleteCategory(t *testing.T) {

}
func TestUpdateCategory(t *testing.T) {

}
func TestUserLogin(t *testing.T) {

}
func TestUserRegister(t *testing.T) {

}
func TestUpdateUser(t *testing.T) {

}
