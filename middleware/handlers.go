package middleware

import (
	
	"database/sql"
	"encoding/json"		// package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models"	// models package where Stock schema is defined
	"log"
	"net/http"	// used to access the request and response object of the api
	"os"		// used to read the environment variable
	"strconv"	// package used to covert string into int type

	"github.com/gorilla/mux"	// used to get the params from the route
	"github.com/joho/godotenv"	// package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	err := godotenv.Load(".env")	
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))		//os.Getenv("POSTGRES_URL") => fetches the value of the POSTGRES_URL environment variable.
	//I want to connect to a PostgreSQL database using the Go database/sql package.
	//_ "github.com/lib/pq"  --> a PostgreSQL driver which registers itself under the name "postgres" when imported.
	if err != nil {
		// panic(err)
		log.Println("Failed to open DB:", err)
		return nil // or handle properly
	}

	//sql.Open does not immediately connect to the database! It prepares the connection pool.
	//You should test the connection with db.Ping():
	err = db.Ping()		//sql.Open(...) -> Creates a DB connection object, doesn't connect yet
						//db.Ping() -> Actually tries to connect, useful for validation/debugging
	if err != nil {
		fmt.Println("Database Unreachable")
		panic(err)
	}

	fmt.Println("Successfully connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request){
	var stock models.Stock

	err :=json.NewDecoder(r.Body).Decode(&stock)  //reads from JSON body (from Postman or frontend) and decodes it into the stock variable

	if err != nil{
		log.Fatalf("Unable to decide the request body. %v", err)
	}

	insertID := insertStock(stock)

	res:=response{
		ID: insertID,
		Message: "stock created successfully",
	}

	json.NewEncoder(w).Encode(res)	//Encodes the res object into JSON and writes it to the response (w)
}

func GetStock(w http.ResponseWriter, r *http.Request){
	params:= mux.Vars(r)

	id , err := strconv.Atoi(params["id"])

	if err!= nil{
		log.Fatalf("Unable to convert the string into int. %v",err)
	}
	stock , err:= getStock(int64(id))

	if err!= nil{
		log.Fatalf("Unable to get stock %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request){
	stocks , err := getAllStocks()

	if err!= nil{
		log.Fatalf("Unable to get all stocks %v", err)
	}

	json.NewEncoder(w).Encode(stocks)

}

func UpdateStock(w http.ResponseWriter, r *http.Request){
	params :=mux.Vars(r)

	id ,err := strconv.Atoi(params["id"])

	if err!= nil{
		log.Fatalf("Unable to convert string into int %v", err)
	}

	var stock models.Stock

	err =json.NewDecoder(r.Body).Decode(&stock) //decode by taking reference as stock
	if err!= nil{
		log.Fatalf("Unable to decode the request body %v", err)
	}

	updateRows:= updateStock(int64(id), stock)

	msg:= fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updateRows)
	res:=response{
		ID :int64(id),
		Message:msg,
	}

	json.NewEncoder(w).Encode(res)
	
}

func DeleteStock(w http.ResponseWriter, r *http.Request){
	params:= mux.Vars(r)
	id,err := strconv.Atoi(params["id"])
	if err != nil{
		log.Fatalf("Unable to convert string to int %v",err)
	}

	deletedRows := deleteStock(int64(id))

	msg:= fmt.Sprintf("Stock deleted successfully. Total rows/records %v", deletedRows)

	res := response{
		ID:int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}


func insertStock(stock models.Stock) int64{
	db:= createConnection()
	defer db.Close()

	sqlStatement :=`INSERT INTO stocks (name, price, company) VALUES ($1,$2,$3) RETURNING stockid` //stocks -> name of a table in your PostgreSQL database stocksdb.
	// Note: the “RETURNING stockid” clause is required here—without it, there’s no value for row.Scan(&id) to read
	var id int64

	row := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company)	
	//db.QueryRow(...) executes the INSERT … RETURNING statement.
		//It substitutes:
		//	$1 → stock.Name
		//	$2 → stock.Price
		//	$3 → stock.Company
		//Returns a *sql.Row which you can call .Scan() on to retrieve the returned column(s).
	if err := row.Scan(&id) ;err != nil{	//row.Scan(&id) copies the first column of the result (stockid) into variable id.
		log.Fatalf("Unable to execute the query or scan the returned id: %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64)(models.Stock,error){  //getStock dont return array as in getAllStocks, no need of [].
	db:= createConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`	//No RETURNING is involved—SELECT by its very nature returns rows, and Scan simply reads those returned columns.

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default: 
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return stock, err
}

func getAllStocks() ([]models.Stock, error){
	db:=  createConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`

	rows, err:= db.Query(sqlStatement)
	if err!= nil{
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()

	for rows.Next(){
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err!= nil{
			log.Fatalf("Unable to scan the row. %v", err)
		}

		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64{
	db := createConnection()
	defer db.Close()

	sqlStatement  := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	res, err:= db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err!=nil{
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowAffected,err := res.RowsAffected()
	if err!= nil{
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total row/records affected %v", rowAffected)
	return rowAffected
}

func deleteStock(id int64) int64{
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil{ 
		log.Fatalf("Unable to execute the query, %v", err)
	}

	rowsAffected , err := res.RowsAffected()
	if err!=nil{
		log.Fatalf("Error while checking affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected. %v \n",rowsAffected)
	return rowsAffected
}