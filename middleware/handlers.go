package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models" // models package where Stock schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route
	

	"go-postgres/database"
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}


func CreateStock(w http.ResponseWriter, r *http.Request){
	var stock models.Stock

	err :=json.NewDecoder(r.Body).Decode(&stock)  //reads from JSON body (from Postman or frontend) and decodes it into the stock variable

	if err != nil{
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	insertID, err := insertStock(stock)
	if err!=nil{
		log.Printf("Error inserting stock: %v", err)
		http.Error(w, "Failed to insert stock", http.StatusInternalServerError)
		return
	}

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
		http.Error(w, fmt.Sprintf("ID must be an integer: %v", err), http.StatusBadRequest)
		return
	}

	stock , err:= getStock(int64(id))
	if err!= nil{
		if err == sql.ErrNoRows {
			http.Error(w, "No stock found", http.StatusNotFound)
		} else {
			log.Printf("Error fetching stock: %v", err)
			http.Error(w, "Error fetching stock", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request){
	stocks , err := getAllStocks()

	if err!= nil{
		log.Printf("Error fetching all stocks: %v", err)
		http.Error(w, "Error fetching all stocks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stocks)

}

func UpdateStock(w http.ResponseWriter, r *http.Request){
	params :=mux.Vars(r)

	id ,err := strconv.Atoi(params["id"])

	if err!= nil{
		http.Error(w, fmt.Sprintf("ID must be an integer: %v", err), http.StatusBadRequest)
		return
	}

	var stock models.Stock

	err =json.NewDecoder(r.Body).Decode(&stock) //decode by taking reference as stock
	if err!= nil{
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	updatedRows,err:= updateStock(int64(id), stock)
	if err != nil {
		log.Printf("Error updating stock: %v", err)
		http.Error(w, "Failed to update stock", http.StatusInternalServerError)
		return
	}
	if updatedRows == 0 {
		http.Error(w, "No stock found to update", http.StatusNotFound)
		return
	}

	msg:= fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updatedRows)
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
		http.Error(w, fmt.Sprintf("ID must be an integer: %v", err), http.StatusBadRequest)
		return
	}

	deletedRows,err  := deleteStock(int64(id))
	if err!=nil{
		log.Printf("Error deleting stock: %v", err)
		http.Error(w, "Failed to delete stock", http.StatusInternalServerError)
		return
	}

	msg:= fmt.Sprintf("Stock deleted successfully. Total rows/records %v", deletedRows)

	res := response{
		ID:int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}


func insertStock(stock models.Stock) (int64, error){
	sqlStatement :=`
		INSERT INTO stocks (name, price, company) 
		VALUES ($1,$2,$3) 
		RETURNING stockid
	` //stocks -> name of a table in your PostgreSQL database stocksdb.
	// Note: the “RETURNING stockid” clause is required here—without it, there’s no value for row.Scan(&id) to read
	var id int64

	row := database.DB.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company)	
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
	return id,nil
}

func getStock(id int64)(models.Stock,error){  //getStock dont return array as in getAllStocks, no need of [].
	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`	//No RETURNING is involved—SELECT by its very nature returns rows, and Scan simply reads those returned columns.

	row := database.DB.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	return stock, err
}

func getAllStocks() ([]models.Stock, error){
	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`

	rows, err:= database.DB.Query(sqlStatement)
	if err!= nil{
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()

	for rows.Next(){
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err!= nil{
			return nil, fmt.Errorf("scan row: %w", err)
		}

		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func updateStock(id int64, stock models.Stock) (int64,error){
	// sqlStatement  := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	sqlStatement := `
		UPDATE stocks
		SET
			name = COALESCE(NULLIF($2,''), name),
			price = CASE WHEN $3 = 0 THEN price ELSE $3 END,
			company = COALESCE(NULLIF($4,''),company)
		WHERE stockid = $1
	`			
	res, err:= database.DB.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err!=nil{
		return 0, fmt.Errorf("execute update: %w", err)
	}

	rowAffected,err := res.RowsAffected()
	if err!= nil{
		return 0, fmt.Errorf("rows affected: %w", err)
	}
	fmt.Printf("Total row/records affected %v", rowAffected)
	return rowAffected, nil
}

func deleteStock(id int64) (int64, error){

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := database.DB.Exec(sqlStatement, id)
	if err != nil{ 
		return 0, fmt.Errorf("execute delete: %w", err)
	}

	rowsAffected , err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	fmt.Printf("Total rows/records affected: %v \n",rowsAffected)
	return rowsAffected, nil
}