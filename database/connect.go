package database

import (
	"database/sql"
	"fmt"
	"log"
	"os" // used to read the environment variable

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      //  Postgres driver registers itself under "postgres"
)

// DB is a global *sql.DB that holds our connection pool.  
// All middleware/handlers can now use database.DB directly.
var DB *sql.DB
// create connection with postgres db
func init(){
	// 1) Load environment variables from .env (only once).
	err := godotenv.Load(".env")	
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	
	// 2) Open a single connection pool to Postgres.  
	//    Note: sql.Open does NOT yet create any network sockets.
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))		//os.Getenv("POSTGRES_URL") => fetches the value of the POSTGRES_URL environment variable.
	//I want to connect to a PostgreSQL database using the Go database/sql package.
	//_ "github.com/lib/pq"  --> a PostgreSQL driver which registers itself under the name "postgres" when imported.

	//sql.Open("postgres", dsn)
    //Creates a new *sql.DB value, which is really a handle to a pool of zero or more connections.
    // No actual network sockets are opened at this point—Go simply prepares a “pool manager” that can open connections as needed.
	if err != nil {
		log.Fatalf("Failed to open DB: %v" , err)
	}

	//sql.Open does not immediately connect to the database! It prepares the connection pool.
	//You should test the connection with db.Ping():
	err = db.Ping()		//sql.Open(...) -> Creates a DB connection object, doesn't connect yet
						//db.Ping() -> Actually tries to connect, useful for validation/debugging(firsty real query)
	// This forces the pool to open at least one live connection to the database and verify it can communicate.
	// If Ping() succeeds, you know the pool can create and use connections. If it fails, you get an error right away.
	if err != nil {
		log.Fatalf("Database Unreachable/Error connecting to DB: %v", err)
	}

	DB=db
	
	fmt.Println("Successfully connected to postgreSQL")
}
