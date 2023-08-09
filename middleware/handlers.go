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
	error := godotenv.Load("../.env")
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
	products, err := getProducts()
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to fetch the products"})
		return
	}
	writeJson(w, http.StatusOK, products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	product, err := getProduct(int64(id))
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to fetch the product"})
		return
	}
	if product.ProductId == 0 {
		writeJson(w, http.StatusNotFound, response{ID: int64(http.StatusNotFound), Message: "product not found"})
	} else {
		writeJson(w, http.StatusOK, product)
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.ProductRequest
	if err := readJson(w, r, &req); err != nil {
		return
	}
	categoryExisting, err := getCategory(int64(req.Category.Id))
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to check if category exists"})
		return
	}
	if categoryExisting.Id == 0 {
		writeJson(w, http.StatusNotFound, response{ID: int64(http.StatusNotFound), Message: "category with that ID doesn't exist"})
	} else {
		id := createProduct(req)
		writeJson(w, http.StatusOK, response{ID: id, Message: "inserted a product"})
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	prodId := deleteProduct(int64(id))
	if prodId == 0 {
		msg := fmt.Sprintf("could not find a product with an id of : %v. total rows affected %v", id, prodId)
		writeJson(w, http.StatusNotFound, response{ID: int64(prodId), Message: msg})
	} else {
		msg := fmt.Sprintf("product deleted successfully. total rows affected %v", prodId)
		writeJson(w, http.StatusNoContent, response{ID: int64(prodId), Message: msg})
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	var req models.ProductRequest
	if err := readJson(w, r, &req); err != nil {
		return
	}
	exists, err := getProduct(int64(id))
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to check if product exists"})
		return
	}
	if exists.ProductId == 0 {
		writeJson(w, http.StatusNotFound, response{ID: int64(id), Message: "product with that ID doesn't exist"})
	} else {
		categoryExisting, err := getCategory(int64(req.Category.Id))
		if err != nil {
			requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to check if category exists"})
			return
		}
		if categoryExisting.Id == 0 {
			writeJson(w, http.StatusNotFound, response{ID: int64(id), Message: "category with that ID doesn't exist"})
		} else {
			rows := updateProduct(req, int64(id))
			msg := fmt.Sprintf("product successfully updated, rows affected %v", rows)
			writeJson(w, http.StatusOK, response{ID: int64(id), Message: msg})
		}
	}
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := getCategories()
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to fetch the categories"})
		return
	}
	writeJson(w, http.StatusOK, categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	category, err := getCategory(int64(id))
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to fetch the category"})
		return
	}
	if category.Id == 0 {
		writeJson(w, http.StatusNotFound, response{ID: int64(http.StatusNotFound), Message: "category not found"})
	} else {
		writeJson(w, http.StatusOK, category)
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CategoryRequest
	if err := readJson(w, r, &req); err != nil {
		return
	}
	id := createCategory(req)
	writeJson(w, http.StatusOK, response{ID: id, Message: "inserted a category"})
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	categoryId := deleteCategory(int64(id))
	if categoryId == 0 {
		msg := fmt.Sprintf("could not find a category with an id of : %v. total rows affected %v", id, categoryId)
		writeJson(w, http.StatusNotFound, response{ID: int64(http.StatusNotFound), Message: msg})
	} else {
		msg := fmt.Sprintf("category deleted successfully. total rows affected %v", categoryId)
		writeJson(w, http.StatusNoContent, response{ID: int64(categoryId), Message: msg})
	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	var req models.CategoryRequest
	if err := readJson(w, r, &req); err != nil {
		return
	}
	exists, err := getCategory(int64(id))
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "cannot fetch category"})
		return
	}
	if exists.Id == 0 {
		writeJson(w, http.StatusNotFound, response{ID: int64(http.StatusNotFound), Message: "category with that ID doesn't exist"})
	} else {
		rows := updateCategory(req, int64(id))
		msg := fmt.Sprintf("category successfully updated, rows affected %v", rows)
		writeJson(w, http.StatusOK, response{ID: int64(id), Message: msg})
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := getId(w, r)
	if err != nil {
		return
	}
	var req models.UserRequest
	if err := readJson(w, r, &req); err != nil {
		return
	}
	exists, err := getUserById(int64(id))
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to fetch the user by id"})
		return
	}
	if exists.Id == 0 {
		writeJson(w, http.StatusNotFound, response{ID: int64(http.StatusNotFound), Message: "user with that ID doesn't exist"})
	} else {
		rows := updateUser(int64(id), req)
		msg := fmt.Sprintf("user successfully updated, rows affected %v", rows)
		writeJson(w, http.StatusOK, response{ID: int64(id), Message: msg})
	}
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := readJson(w, r, &req); err != nil {
		return
	}
	user, err := getUserByEmail(req.Email)
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to fetch the user by email"})
		return
	}
	if user.Email != req.Email || !validPassword(user.Password, req.Password) {
		writeJson(w, http.StatusBadRequest, response{ID: int64(http.StatusBadRequest), Message: "email or password are incorrect"})
		return
	}
	token, err := createJwt(user)
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "unable to generate jwt"})
		return
	}
	res := models.LoginResponse{
		Token: token,
	}
	writeJson(w, http.StatusOK, res)
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := readJson(w, r, &user); err != nil {
		return
	}
	if !validEmail(user.Email) {
		writeJson(w, http.StatusBadRequest, response{ID: int64(http.StatusBadRequest), Message: "email not valid"})
	} else if existing, _ := getUserByEmail(user.Email); existing.Id != 0 {
		writeJson(w, http.StatusBadRequest, response{ID: int64(http.StatusBadRequest), Message: "user already exists"})
	} else if len(user.Password) < 1 {
		writeJson(w, http.StatusBadRequest, response{ID: int64(http.StatusBadRequest), Message: "password too short"})
	} else {
		id := createUser(user)
		writeJson(w, http.StatusOK, response{ID: id, Message: "registration successful"})
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
		return user, error
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

func updateUser(id int64, req models.UserRequest) int64 {
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

func createProduct(new models.ProductRequest) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO products (name, short_description, description, price, quantity, category_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int64
	error := db.QueryRow(sqlStatement, new.Name, new.ShortDescription, new.Description, new.Price, new.Quantity, new.Category.Id).Scan(&id)
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

func updateProduct(update models.ProductRequest, id int64) int64 {
	db := createConnection()
	defer db.Close()
	updatedAt := time.Now()
	sqlStatement := `UPDATE products SET name = $2, short_description = $3, description = $4, price = $5, updated_at = $6, quantity = $7, category_id = $8 WHERE id = $1`
	res, error := db.Exec(sqlStatement, id, update.Name, update.ShortDescription, update.Description, update.Price, updatedAt, update.Quantity, update.Category.Id)
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

func createCategory(new models.CategoryRequest) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO categories (category_name) VALUES ($1) RETURNING category_id`
	var id int64
	error := db.QueryRow(sqlStatement, new.Name).Scan(&id)
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

func updateCategory(update models.CategoryRequest, id int64) int64 {
	db := createConnection()
	defer db.Close()
	updatedAt := time.Now()
	sqlStatement := `UPDATE categories SET category_name = $2, updated_at = $3 WHERE category_id = $1`
	res, error := db.Exec(sqlStatement, id, update.Name, updatedAt)
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

func JwtAuth(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("jwt")
		token, err := validJwt(tokenString)
		if err != nil {
			requestError(w, err, response{ID: int64(http.StatusForbidden), Message: "permission denied"})
			return
		}
		if !token.Valid {
			requestError(w, err, response{ID: int64(http.StatusForbidden), Message: "permission denied"})
			return
		}
		// if we dont want to allow the user to change other accounts
		//
		// id, err := getId(w, r)
		// if err != nil {
		// 	return
		// }
		// user, err := getUserById(int64(id))
		// if err != nil {
		// 	requestError(w, err, response{ID: int64(http.StatusNotFound), Message: "cannot fetch the user"})
		// 	return
		// }
		// claims := token.Claims.(jwt.MapClaims)
		// if user.Email != claims["email"] {
		// 	writeJson(w, http.StatusForbidden, response{ID: int64(http.StatusForbidden), Message: "permission denied"})
		// 	return
		// }
		handleFunc(w, r)
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

func getId(w http.ResponseWriter, r *http.Request) (int, error) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		requestError(w, err, response{ID: int64(http.StatusBadRequest), Message: "cannot parse given id"})
		return id, err
	}
	return id, err
}

func readJson(w http.ResponseWriter, r *http.Request, v any) error {
	error := json.NewDecoder(r.Body).Decode(&v)
	if error != nil {
		requestError(w, error, response{ID: int64(http.StatusBadRequest), Message: "cannot decode request body"})
	}
	return error
}

func writeJson(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func requestError(w http.ResponseWriter, err error, res response) {
	fmt.Println(err)
	writeJson(w, int(res.ID), res)
}
