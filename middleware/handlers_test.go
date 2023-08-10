package middleware

import (
	"bytes"
	"encoding/json"
	"go-products-rest/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetProducts(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProducts)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted status code %v, got %v", http.StatusOK, rr.Code)
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
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProduct)
	handler.ServeHTTP(rr, r)
	var product models.Product
	err = json.Unmarshal(rr.Body.Bytes(), &product)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}
	id, _ := strconv.Atoi(vars["id"])
	if id != product.ProductId {
		t.Errorf("Wanted a product with an ID of %v, got %v", id, product)
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted status code %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestCreateProduct(t *testing.T) {
	test := models.ProductRequest{
		Name:             "test",
		ShortDescription: "test",
		Description:      "test",
		Price:            1,
		Quantity:         10,
		Category:         models.Category{Id: 1},
	}
	jsonProduct, err := json.Marshal(test)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("POST", "/api/products", bytes.NewBuffer(jsonProduct))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusCreated {
		log.Fatalf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
	var res response
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}
	product, err := getProduct(res.ID)
	id := product.ProductId
	expectedRes := response{
		ID:      int64(id),
		Message: "inserted a product",
	}
	if res.ID < 0 {
		log.Fatalf("Id is less than zero %v", err)
	}
	if res.ID != expectedRes.ID || res.Message != expectedRes.Message {
		log.Fatalf("Expected response %v, but got %v", expectedRes, res)
	}
	deleteProduct(int64(id))
}

func TestUpdateProduct(t *testing.T) {
	test := models.ProductRequest{
		Name:             "test",
		ShortDescription: "test",
		Description:      "test",
		Price:            1.0,
		Quantity:         1,
		Category:         models.Category{Id: 1},
	}
	id := createProduct(test)
	updateTest := models.ProductRequest{
		Name:             "test 2",
		ShortDescription: "test 2",
		Description:      "test 2",
		Price:            2.0,
		Quantity:         2,
		Category:         models.Category{Id: 2},
	}
	updateJson, err := json.Marshal(updateTest)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("PUT", "/api/products/{id}", bytes.NewBuffer(updateJson))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	ids := strconv.Itoa(int(id))
	vars := map[string]string{
		"id": ids,
	}
	r = mux.SetURLVars(r, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateProduct)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		log.Fatalf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
	product, err := getProduct(id)
	if updateTest.Name != product.Name ||
		updateTest.ShortDescription != product.ShortDescription ||
		updateTest.Description != product.Description ||
		updateTest.Price != product.Price ||
		updateTest.Quantity != product.Quantity ||
		updateTest.Category.Id != product.Category.Id {
		log.Fatalf("Expected result %v, but got %v", updateTest, product)
	}
	deleteProduct(int64(id))
}

func TestDeleteProduct(t *testing.T) {
	test := models.ProductRequest{
		Name:             "test",
		ShortDescription: "test",
		Description:      "test",
		Price:            1.0,
		Quantity:         1,
		Category:         models.Category{Id: 1},
	}
	id := createProduct(test)
	r, err := http.NewRequest("DELETE", "/api/products/{id}", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	ids := strconv.Itoa(int(id))
	vars := map[string]string{
		"id": ids,
	}
	r = mux.SetURLVars(r, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteProduct)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusNoContent {
		log.Fatalf("Expected status code %v, but got %v", http.StatusNoContent, rr.Code)
	}
	product, _ := getProduct(id)
	if product.ProductId != 0 {
		t.Errorf("Product is not deleted %v", product)
	}
}
func TestGetCategories(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/categories", GetCategories)
	router.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted status code %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestGetCategory(t *testing.T) {
	r, err := http.NewRequest("GET", "/api/categories/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "1",
	}
	r = mux.SetURLVars(r, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCategory)
	handler.ServeHTTP(rr, r)
	var category models.Category
	err = json.Unmarshal(rr.Body.Bytes(), &category)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}
	id, _ := strconv.Atoi(vars["id"])
	if id != category.Id {
		t.Errorf("Wanted a product with an ID of %v, got %v", id, category)
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Wanted status code %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestCreateCategory(t *testing.T) {
	test := models.CategoryRequest{
		Name: "test",
	}
	jsonProduct, err := json.Marshal(test)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(jsonProduct))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCategory)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusCreated {
		log.Fatalf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
	var res response
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}
	category, err := getCategory(res.ID)
	id := category.Id
	expectedRes := response{
		ID:      int64(id),
		Message: "inserted a category",
	}
	if res.ID < 0 {
		log.Fatalf("Id is less than zero %v", err)
	}
	if res.ID != expectedRes.ID || res.Message != expectedRes.Message {
		log.Fatalf("Expected response %v, but got %v", expectedRes, res)
	}
	deleteCategory(int64(id))
}

func TestUpdateCategory(t *testing.T) {
	test := models.CategoryRequest{
		Name: "test",
	}
	id := createCategory(test)
	updateTest := models.CategoryRequest{
		Name: "updated test",
	}
	updateJson, err := json.Marshal(updateTest)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("PUT", "/api/categories/{id}", bytes.NewBuffer(updateJson))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	ids := strconv.Itoa(int(id))
	vars := map[string]string{
		"id": ids,
	}
	r = mux.SetURLVars(r, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateCategory)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		log.Fatalf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
	category, err := getCategory(id)
	if updateTest.Name != category.Name {
		log.Fatalf("Expected name %v, but got %v", updateTest.Name, category.Name)
	}
	deleteCategory(int64(id))
}

func TestDeleteCategory(t *testing.T) {
	test := models.CategoryRequest{
		Name: "test",
	}
	id := createCategory(test)
	r, err := http.NewRequest("DELETE", "/api/categories/{id}", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	ids := strconv.Itoa(int(id))
	vars := map[string]string{
		"id": ids,
	}
	r = mux.SetURLVars(r, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteCategory)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusNoContent {
		log.Fatalf("Expected status code %v, but got %v", http.StatusNoContent, rr.Code)
	}
	category, _ := getCategory(id)
	if category.Id != 0 {
		t.Errorf("Category is not deleted %v", category)
	}
}

func TestUserRegister(t *testing.T) {
	test := models.User{
		FirstName: "Test",
		LastName:  "Testic",
		Email:     "t@t.com",
		Password:  "ttt",
	}
	jsonProduct, err := json.Marshal(test)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonProduct))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserRegister)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusCreated {
		log.Fatalf("Expected status code %v, but got %v", http.StatusCreated, rr.Code)
	}
	var res response
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}
	user, err := getUserById(res.ID)
	id := user.Id
	expectedRes := response{
		ID:      int64(id),
		Message: "registration successful",
	}
	if res.ID < 0 {
		log.Fatalf("Id is less than zero %v", err)
	}
	if res.ID != expectedRes.ID || res.Message != expectedRes.Message {
		log.Fatalf("Expected response %v, but got %v", expectedRes, res)
	}
	deleteUser(int64(id))
}

func TestUserLogin(t *testing.T) {
	user := models.User{
		FirstName: "Test",
		LastName:  "Testic",
		Email:     "t@t.com",
		Password:  "ttt",
	}
	id := createUser(user)
	logReq := models.LoginRequest{
		Email:    "t@t.com",
		Password: "ttt",
	}
	jsonProduct, err := json.Marshal(logReq)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonProduct))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserLogin)
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		log.Fatalf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
	var res models.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body: %v", err)
	}
	token, err := validJwt(res.Token)
	if err != nil {
		log.Fatal(err)
	}
	if !token.Valid {
		t.Error("Token is not valid")
	}
	deleteUser(int64(id))
}

func TestUpdateUser(t *testing.T) {
	newUser := models.User{
		FirstName: "Test",
		LastName:  "Testic",
		Email:     "t@t.com",
		Password:  "ttt",
	}
	id := createUser(newUser)
	updateUser := models.UserRequest{
		FirstName: "Test 2",
		LastName:  "Testic 2",
	}
	updateJson, err := json.Marshal(updateUser)
	if err != nil {
		log.Fatalf("Failed to marshal product to JSON: %v", err)
	}
	r, err := http.NewRequest("PUT", "/api/users/{id}", bytes.NewBuffer(updateJson))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	ids := strconv.Itoa(int(id))
	vars := map[string]string{
		"id": ids,
	}
	token, err := createJwt(newUser)
	if err != nil {
		log.Fatalf("Unable to generate token %v", err)
	}
	r.Header.Set("Authorization", token)
	r = mux.SetURLVars(r, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(JwtAuth(UpdateUser))
	handler.ServeHTTP(rr, r)
	if rr.Code != http.StatusOK {
		log.Fatalf("Expected status code %v, but got %v", http.StatusOK, rr.Code)
	}
	user, err := getUserById(id)
	if updateUser.FirstName != user.FirstName || updateUser.LastName != user.LastName {
		log.Fatalf("Expected %v, but got %v", updateUser, user)
	}
	deleteUser(int64(id))
}
