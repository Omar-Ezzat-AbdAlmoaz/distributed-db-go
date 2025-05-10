package database

import (
	"database/sql"
	"fmt"
	"log"

	"distributed-db-go/utils"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectOrCreateDatabase(user, password, host, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", user, password, host)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return fmt.Errorf("CREATE DATABASE failed: %v", err)
	}

	db.Close()

	fullDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err = sql.Open("mysql", fullDSN)
	if err != nil {
		return fmt.Errorf("Reconnect failed: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Ping failed: %v", err)
	}

	DB = db

	_, err = DB.Exec("USE " + dbName)
	if err != nil {
		return fmt.Errorf("Database doesn't exist or not selected")
	}
	//fmt.Println("Connected to database:", dbName)
	return nil
}

func ConnectToExistingDatabase(user, password, host, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("Failed to connect to existing DB: %v", err)
	}

	DB = db

	_, err = DB.Exec("USE " + dbName)

	if err != nil {
		return fmt.Errorf("Database doesn't exist or not selected")
	}
	//fmt.Println("Connected to existing database:", dbName)
	return nil
}

func CreateTable(table string, columns []string) error {
	if !utils.IsMaster {
		return fmt.Errorf("Permission denied: Slaves cannot create tables")
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY", table)
	for _, col := range columns {
		query += fmt.Sprintf(", %s VARCHAR(255)", col)
	}
	query += ")"

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("CreateTable failed: %v", err)
	}
	//fmt.Println("Create:", table,"table")
	return nil
}

func DeleteTable(table string) error {
	if !utils.IsMaster {
		return fmt.Errorf("Permission denied: Slaves cannot delete tables")
	}
	_, err := DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	if err != nil {
		return fmt.Errorf("DeleteTable failed: %v", err)
	}
	//fmt.Println("Delete:", table,"table")
	return nil
}

func DropDatabase(dbName string) error {
	if !utils.IsMaster {
		return fmt.Errorf("Permission denied: Slaves cannot drop database")
	}

	dsn := "root:rootroot@tcp(localhost:3306)/"
	tmpDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Failed to connect without DB: %v", err)
	}
	defer tmpDB.Close()

	// ننفذ DROP DATABASE
	_, err = tmpDB.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		return fmt.Errorf("DROP DATABASE failed: %v", err)
	}

	//fmt.Println("Dropped database:", dbName)

	if DB != nil {
		DB.Close()
	}
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Reconnect failed after drop:", err)
	}

	return nil
}
