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
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	if error != nil {
		log.Fatalf("unable to check if product exists %v", error)
	}
	if product.ProductId != 0 {
		res := response{
			ID:      400,
			Message: "creating products with existing ID's not allowed",
		}
		w.WriteHeader(400)
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
		error = rows.Scan(&product.ProductId, &product.Name, &product.ShortDescription, &product.Description, &product.Price,
			&product.CreatedAt, &product.UpdatedAt)
		if error != nil {
			log.Fatalf("unable to scan row %v", error)
		}
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
	error := row.Scan(&product.ProductId, &product.Name, &product.ShortDescription, &product.Description, &product.Price,
		&product.CreatedAt, &product.UpdatedAt)
	switch error {
	case sql.ErrNoRows:
		fmt.Println("no rows were returned")
		return product, nil
	case nil:
		return product, nil
	default:
		log.Fatalf("unable to scan rows %v", error)
	}
	return product, error

}
func createProduct(product models.Product) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO products (name, short_description, description, price) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	error := db.QueryRow(sqlStatement, product.Name, product.ShortDescription, product.Description, product.Price).Scan(&id)
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
	sqlStatement := `UPDATE products SET name = $2, short_description = $3, description = $4, price = $5, updated_at = $6 WHERE id = $1`
	res, error := db.Exec(sqlStatement, product.ProductId, product.Name, product.ShortDescription, product.Description, product.Price, product.UpdatedAt)
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
