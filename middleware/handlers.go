package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-products-rest/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type response struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

func createConnection() *sql.DB {
	error := godotenv.Load(".env")
	if error != nil {
		log.Fatal("error loading .env file")
	}
	db, error := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if error != nil {
		panic(error)
	}
	error = db.Ping()
	if error != nil {
		panic(error)
	}
	fmt.Println("successfully connected to postgres")
	return db
}
func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products, error := getProducts()
	if error != nil {
		log.Fatalf("unable to fetch products %v", error)
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(products)
}
func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse product id %v", error)
	}
	product, error := getProduct(int64(id))
	if error != nil {
		log.Fatalf("unable to fetch product %v", error)
	}
	if product.ProductId == 0 {
		res := response{
			ID:      404,
			Message: "Not found",
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(product)
	}
}
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product models.Product
	error := json.NewDecoder(r.Body).Decode(&product)
	if error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	if product.ProductId != 0 {
		res := response{
			ID:      400,
			Message: "creating products with existing ID's not allowed",
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
	} else {
		categoryExisting, error := getCategory(int64(product.Category.Id))
		if error != nil {
			log.Fatalf("unable to check if category exists %v", error)
		}
		if categoryExisting.Id == 0 {
			res := response{
				ID:      404,
				Message: "category with that ID doesn't exist",
			}
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(res)
		} else {
			id := createProduct(product)
			res := response{
				ID:      id,
				Message: "inserted a product",
			}
			json.NewEncoder(w).Encode(res)
		}
	}
}
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse id %v", error)
	}
	prodId := deleteProduct(int64(id))
	msg1 := fmt.Sprintf("product deleted successfully. total rows affected %v", prodId)
	msg2 := fmt.Sprintf("could not find a project with an id of : %v. total rows affected %v", id, prodId)
	if prodId == 0 {
		res := response{
			ID:      404,
			Message: msg2,
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
	} else {
		res := response{
			ID:      int64(prodId),
			Message: msg1,
		}
		w.WriteHeader(204)
		json.NewEncoder(w).Encode(res)
	}
}
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse id %v", error)
	}
	var product models.Product
	error = json.NewDecoder(r.Body).Decode(&product)
	if error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	msg2 := fmt.Sprintf("path variable %v not equal to product id %v", id, product.ProductId)
	exists, error := getProduct(int64(product.ProductId))
	if error != nil {
		log.Fatalf("unable to check if product exists %v", error)
	}
	if id != product.ProductId {
		res := response{
			ID:      400,
			Message: msg2,
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
	} else if exists.ProductId == 0 {
		res := response{
			ID:      404,
			Message: "product with that ID doesn't exist",
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
		return
	} else {
		categoryExisting, error := getCategory(int64(product.Category.Id))
		if error != nil {
			log.Fatalf("unable to check if category exists %v", error)
		}
		if categoryExisting.Id == 0 {
			res := response{
				ID:      404,
				Message: "Category with that ID doesn't exist",
			}
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(res)
		} else {
			rows := updateProduct(product)
			msg2 := fmt.Sprintf("product successfully updated, rows affected %v", rows)
			res := response{
				ID:      int64(id),
				Message: msg2,
			}
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(res)
		}
	}
}
func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, error := getCategories()
	if error != nil {
		log.Fatalf("unable to fetch categories %v", error)
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(categories)
}
func GetCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse category id %v", error)
	}
	category, error := getCategory(int64(id))
	if error != nil {
		log.Fatalf("unable to fetch category %v", error)
	}
	if category.Id == 0 {
		res := response{
			ID:      404,
			Message: "Not found",
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(category)
	}
}
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var category models.Category
	error := json.NewDecoder(r.Body).Decode(&category)
	if error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	if category.Id != 0 {
		res := response{
			ID:      400,
			Message: "creating categories with existing ID's not allowed",
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
	} else {
		id := createCategory(category)
		res := response{
			ID:      id,
			Message: "inserted a category",
		}
		json.NewEncoder(w).Encode(res)
	}
}
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse id %v", error)
	}
	categoryId := deleteCategory(int64(id))
	msg1 := fmt.Sprintf("category deleted successfully. total rows affected %v", categoryId)
	msg2 := fmt.Sprintf("could not find a category with an id of : %v. total rows affected %v", id, categoryId)
	if categoryId == 0 {
		res := response{
			ID:      404,
			Message: msg2,
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
	} else {
		res := response{
			ID:      int64(categoryId),
			Message: msg1,
		}
		w.WriteHeader(204)
		json.NewEncoder(w).Encode(res)
	}
}
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse id %v", error)
	}
	var category models.Category
	error = json.NewDecoder(r.Body).Decode(&category)
	if error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	msg2 := fmt.Sprintf("path variable %v not equal to category id %v", id, category.Id)
	exists, error := getCategory(int64(category.Id))
	if error != nil {
		log.Fatalf("unable to check if category exists %v", error)
	}
	if id != category.Id {
		res := response{
			ID:      400,
			Message: msg2,
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
	} else if exists.Id == 0 {
		res := response{
			ID:      404,
			Message: "category with that ID doesn't exist",
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
	} else {
		rows := updateCategory(category)
		msg2 := fmt.Sprintf("category successfully updated, rows affected %v", rows)
		res := response{
			ID:      int64(id),
			Message: msg2,
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
	}
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, error := strconv.Atoi(params["id"])
	if error != nil {
		log.Fatalf("unable to parse id %v", error)
	}
	var req models.UpdateRequest
	if error := json.NewDecoder(r.Body).Decode(&req); error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	exists, error := getUserById(int64(id))
	if error != nil {
		log.Fatalf("unable to fetch user %v", error)
	}
	if exists.Id == 0 {
		res := response{
			ID:      404,
			Message: "user with that ID doesn't exist",
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(res)
	} else {
		rows := updateUser(int64(id), req)
		msg := fmt.Sprintf("user successfully updated, rows affected %v", rows)
		res := response{
			ID:      int64(id),
			Message: msg,
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
	}

}
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.LoginRequest
	if error := json.NewDecoder(r.Body).Decode(&req); error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	user, error := getUserByEmail(req.Email)
	if error != nil {
		log.Fatalf("unable to fetch user by email %v", error)
	}
	if user.Email != req.Email || !validPassword(user.Password, req.Password) {
		res := response{
			ID:      400,
			Message: "email or password are incorrect",
		}
		fmt.Printf("%s", user.Password)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
	} else {
		token, err := createJwt(user)
		if error != nil {
			log.Fatalf("unable to generate jwt %v", err)
		}
		res := models.LoginResponse{
			Token: token,
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res)
	}

}
func UserRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if error := json.NewDecoder(r.Body).Decode(&user); error != nil {
		log.Fatalf("unable to decode request body %v", error)
	}
	if user.Id != 0 {
		res := response{
			ID:      400,
			Message: "creating users with existing ID's not allowed",
		}
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(res)
	} else if !validEmail(user.Email) {
		w.WriteHeader(400)
		res := response{
			ID:      400,
			Message: "email not valid",
		}
		json.NewEncoder(w).Encode(res)
	} else if existing, _ := getUserByEmail(user.Email); existing.Id != 0 {
		w.WriteHeader(400)
		res := response{
			ID:      400,
			Message: "user already exists",
		}
		json.NewEncoder(w).Encode(res)
	} else if len(user.Password) < 1 {
		w.WriteHeader(400)
		res := response{
			ID:      400,
			Message: "password too short",
		}
		json.NewEncoder(w).Encode(res)
	} else {
		id := createUser(user)
		res := response{
			ID:      id,
			Message: "registration successful",
		}
		json.NewEncoder(w).Encode(res)
	}
}

func createUser(user models.User) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	encPw, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if error := db.QueryRow(sqlStatement, user.FirstName, user.LastName, user.Email, string(encPw)).Scan(&id); error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	return id
}
func getUserById(id int64) (models.User, error) {
	db := createConnection()
	defer db.Close()
	var user models.User
	sqlStatement := `SELECT * FROM users WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	error := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	switch error {
	case sql.ErrNoRows:
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("unable to scan the rows %v", error)
	}
	return user, error
}
func getUserByEmail(email string) (models.User, error) {
	db := createConnection()
	defer db.Close()
	var user models.User
	sqlStatement := `SELECT * FROM users WHERE email = $1`
	row := db.QueryRow(sqlStatement, email)
	error := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	switch error {
	case sql.ErrNoRows:
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("unable to scan the rows %v", error)
	}
	return user, error
}
func updateUser(id int64, req models.UpdateRequest) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE users SET first_name = $2, last_name = $3 WHERE id = $1`
	res, error := db.Exec(sqlStatement, id, req.FirstName, req.LastName)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	rows, error := res.RowsAffected()
	if error != nil {
		log.Fatalf("unable to scan rows affected %v", error)
	}
	fmt.Printf("total rows affected %v", rows)
	return rows
}
func getProducts() ([]models.Product, error) {
	db := createConnection()
	defer db.Close()
	var products []models.Product
	sqlStatement := `SELECT * FROM products`
	rows, error := db.Query(sqlStatement)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		var categoryId int64
		error = rows.Scan(&product.ProductId, &product.Name, &product.ShortDescription, &product.Description, &product.Price,
			&product.CreatedAt, &product.UpdatedAt, &product.Quantity, &categoryId)
		if error != nil {
			log.Fatalf("unable to scan row %v", error)
		}
		product.Category, error = getCategory(categoryId)
		products = append(products, product)
	}
	return products, error
}
func getProduct(id int64) (models.Product, error) {
	db := createConnection()
	defer db.Close()
	var product models.Product
	sqlStatement := `SELECT * FROM products WHERE id = $1`
	row := db.QueryRow(sqlStatement, id)
	var categoryId int64
	error := row.Scan(&product.ProductId, &product.Name, &product.ShortDescription, &product.Description, &product.Price,
		&product.CreatedAt, &product.UpdatedAt, &product.Quantity, &categoryId)
	switch error {
	case sql.ErrNoRows:
		fmt.Println("no rows were returned")
		return product, nil
	case nil:
		product.Category, error = getCategory(categoryId)
		if error != nil {
			log.Fatalf("unable to get product's category %v", error)
		}
		return product, nil
	default:
		log.Fatalf("unable to scan rows %v", error)
	}
	return product, error
}
func createProduct(product models.Product) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO products (name, short_description, description, price, quantity, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int64
	error := db.QueryRow(sqlStatement, product.Name, product.ShortDescription, product.Description, product.Price, product.Quantity, product.Category.Id).Scan(&id)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	return id
}
func deleteProduct(id int64) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM products WHERE id = $1`
	res, error := db.Exec(sqlStatement, id)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	rows, error := res.RowsAffected()
	if error != nil {
		log.Fatalf("unable to get affected rows %v", error)
	}
	return rows
}
func updateProduct(product models.Product) int64 {
	db := createConnection()
	defer db.Close()
	product.UpdatedAt = time.Now()
	sqlStatement := `UPDATE products SET name = $2, short_description = $3, description = $4, price = $5, updated_at = $6, quantity = $7, category_id = $8 WHERE id = $1`
	res, error := db.Exec(sqlStatement, product.ProductId, product.Name, product.ShortDescription, product.Description, product.Price, product.UpdatedAt, product.Quantity, product.Category.Id)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	rows, error := res.RowsAffected()
	if error != nil {
		log.Fatalf("unable to scan rows affected %v", error)
	}
	fmt.Printf("total rows affected %v", rows)
	return rows
}
func getCategories() ([]models.Category, error) {
	db := createConnection()
	defer db.Close()
	var categories []models.Category
	sqlStatement := `SELECT * FROM categories`
	rows, error := db.Query(sqlStatement)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		error = rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if error != nil {
			log.Fatalf("unable to scan row %v", error)
		}
		categories = append(categories, category)
	}
	return categories, error
}
func getCategory(id int64) (models.Category, error) {
	db := createConnection()
	defer db.Close()
	var category models.Category
	sqlStatement := `SELECT * FROM categories WHERE category_id = $1`
	row := db.QueryRow(sqlStatement, id)
	error := row.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	switch error {
	case sql.ErrNoRows:
		fmt.Println("no rows were returned")
		return category, nil
	case nil:
		return category, nil
	default:
		log.Fatalf("unable to scan rows %v", error)
	}
	return category, error
}
func createCategory(category models.Category) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO categories (category_name) VALUES ($1) RETURNING category_id`
	var id int64
	error := db.QueryRow(sqlStatement, category.Name).Scan(&id)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	return id
}
func deleteCategory(id int64) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM categories WHERE category_id = $1`
	res, error := db.Exec(sqlStatement, id)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	rows, error := res.RowsAffected()
	if error != nil {
		log.Fatalf("unable to get affected rows %v", error)
	}
	return rows
}
func updateCategory(category models.Category) int64 {
	db := createConnection()
	defer db.Close()
	category.UpdatedAt = time.Now()
	sqlStatement := `UPDATE categories SET category_name = $2, updated_at = $3 WHERE category_id = $1`
	res, error := db.Exec(sqlStatement, category.Id, category.Name, category.UpdatedAt)
	if error != nil {
		log.Fatalf("unable to execute the query %v", error)
	}
	rows, error := res.RowsAffected()
	if error != nil {
		log.Fatalf("unable to scan rows affected %v", error)
	}
	fmt.Printf("total rows affected %v", rows)
	return rows
}
func createJwt(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt": 15000,
		"email":     user.Email,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func JwtAuth(foo http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("jwt")
		token, err := validJwt(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		foo(w, r)
	}
}
func validJwt(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}
func validPassword(encPw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encPw), []byte(pw)) == nil
}
func validEmail(email string) bool {
	if !strings.ContainsAny(email, "@") {
		return false
	} else {
		return true
	}
}
func permissionDenied(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(403)
	res := response{
		ID:      403,
		Message: "permission denied",
	}
	json.NewEncoder(w).Encode(res)
}
