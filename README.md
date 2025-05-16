# StockTrack: Building a RESTful Stock Management API with Go and PostgreSQL

**StockTrack** is a backend service developed using the Go programming language and PostgreSQL as the relational database. This project is ideal for computer science students and developers who are learning about web development, RESTful APIs, SQL, and backend data management.

---

## 📚 Table of Contents

- [StockTrack: Building a RESTful Stock Management API with Go and PostgreSQL](#stocktrack-building-a-restful-stock-management-api-with-go-and-postgresql)
  - [📚 Table of Contents](#-table-of-contents)
  - [✅ Prerequisites](#-prerequisites)
  - [🗄️ Database Setup](#️-database-setup)
    - [Installing PostgreSQL and `psql`](#installing-postgresql-and-psql)
    - [Accessing the `postgres` System User](#accessing-the-postgres-system-user)
    - [Launching the `psql` Shell](#launching-the-psql-shell)
    - [Creating a Role and Database](#creating-a-role-and-database)
    - [Setting a Password](#setting-a-password)
    - [Granting Privileges](#granting-privileges)
  - [🧱 Defining the `stocks` Table](#-defining-the-stocks-table)
  - [🔍 Exploring Tables and Data](#-exploring-tables-and-data)
  - [⚙️ Configuring Environment Variables](#️-configuring-environment-variables)
    - [Steps to Set Up](#steps-to-set-up)
      - [Why Use a `.env` File?](#why-use-a-env-file)
    - [Load `.env` in Go:](#load-env-in-go)
    - [Alternative: Manual Export](#alternative-manual-export)
      - [When to Use `export`?](#when-to-use-export)
      - [Verifying Your Database Connection](#verifying-your-database-connection)
  - [🚀 Running the Server](#-running-the-server)
  - [Troubleshooting \& Tips](#troubleshooting--tips)

---

## ✅ Prerequisites

Before you begin, ensure the following tools are installed on your machine:

- Go (version 1.20 or newer)
- PostgreSQL (version 12 or newer)
- `psql` command-line interface
- Git

---

## 🗄️ Database Setup

### Installing PostgreSQL and `psql`

```bash
# Ubuntu/Debian:
sudo apt update && sudo apt install postgresql postgresql-contrib

# macOS (using Homebrew):
brew update && brew install postgresql

# Start the PostgreSQL service:
sudo service postgresql start   # Linux
brew services start postgresql  # macOS
```

### Accessing the `postgres` System User

```bash
sudo -i -u postgres
```

This command allows you to switch to PostgreSQL’s superuser account.

### Launching the `psql` Shell

```bash
psql
```

### Creating a Role and Database

Inside the `psql` shell, run the following:

```sql
CREATE ROLE stocks_user WITH LOGIN;
CREATE DATABASE stocksdb OWNER stocks_user;
```

### Setting a Password

To assign a password to the user:

```sql
\password stocks_user
```

Avoid using special characters such as `#`, `!`, or `$` in your password, as they may interfere with connection strings.

### Granting Privileges

(Optional but recommended)

```sql
GRANT ALL PRIVILEGES ON DATABASE stocksdb TO stocks_user;
```

To exit the shell:

```bash
\q
exit
```

---

## 🧱 Defining the `stocks` Table

Reconnect to the database using the following command:

```bash
psql -U stocks_user -d stocksdb
```

Then create the `stocks` table:

```sql
CREATE TABLE stocks (
  stockid SERIAL PRIMARY KEY,
  name     TEXT,
  price    INT,
  company  TEXT
);
```

This defines a table with:
- `stockid`: An auto-incrementing integer serving as the primary key
- `name`: The stock name (text)
- `price`: The stock price (integer)
- `company`: The company name (text)

You may insert a sample record for testing purposes:

```sql
INSERT INTO stocks (name, price, company)
  VALUES ('TCS', 3600, 'Tata Consultancy');
```

This command demonstrates how to add a new row to your table using SQL.

---

## 🔍 Exploring Tables and Data

Useful `psql` commands for database inspection:

```sql
#List all databases
\l

#Connect to your database
\c stocksdb

#List all tables
\dt

#Display table schema
\d stocks

#Retrieve all data in the table
SELECT * FROM stocks;
```

Sample output of `\d stocks`:

```
                               Table "public.stocks"
 Column  |  Type   | Collation | Nullable |                 Default                 
---------+---------+-----------+----------+-----------------------------------------
 stockid | integer |           | not null | nextval('stocks_stockid_seq'::regclass)
 name    | text    |           |          | 
 price   | integer |           |          | 
 company | text    |           |          | 
Indexes:
    "stocks_pkey" PRIMARY KEY, btree (stockid)
```

---

## ⚙️ Configuring Environment Variables

Create a `.env` file in your project root:

<!-- ### .env.example -->

```bash
POSTGRES_URL="postgres://<username>:<password>@localhost:5432/stocksdb"
```

### Steps to Set Up

```bash
cp .env.example .env
nano .env  # Edit the file to update your credentials
```



#### Why Use a `.env` File?
- Keeps sensitive information out of your source code
- Enables easy configuration management
- Improves portability

### Load `.env` in Go:

```bash
import "github.com/joho/godotenv"

func init() {
    godotenv.Load(".env")
}
```

### Alternative: Manual Export

```bash
export POSTGRES_URL="postgres://<username>:<password>@localhost:5432/stocksdb"
```
#### When to Use `export`? 

- If `godotenv.Load()` is not used

- When running tools that rely on the variable

- In environments like Docker or CI/CD


####  Verifying Your Database Connection

Check if your `.env` connection string works:

```bash
psql "$POSTGRES_URL"
```

If successful, you’ll see the `stocksdb=#` prompt.

---

## 🚀 Running the Server

Execute the following commands to run your Go server:

```bash
#Install dependencies
go mod tidy

#Run the application
go run main.go

#OR build and run the binary
go build -o stocks-app
./stocks-app
```

Your backend service should now be running on the port specified in your `.env` file (default: 8080).

---

<!-- ## 🧯 Troubleshooting Tips -->
## Troubleshooting & Tips

- **Switch-user errors**: If `sudo -i -u postgres` fails, try `sudo su - postgres`.
- **Password issues**: Use `\password stocks_user` as shown above to reset.
- **Connection errors**: Ensure PostgreSQL is running on port `5432`.
- **Go module errors**: Run `go mod verify` to check dependencies.
