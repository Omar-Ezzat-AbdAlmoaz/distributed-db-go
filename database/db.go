package database

import "fmt"

type Database struct {
	Tables map[string]*Table
}

// إنشاء قاعدة بيانات جديدة
func NewDatabase() *Database {
	return &Database{
		Tables: make(map[string]*Table),
	}
}

// إنشاء جدول جديد
func (db *Database) CreateTable(name string, columns []string) error {
	if _, exists := db.Tables[name]; exists {
		return fmt.Errorf("Table %s already exists", name)
	}
	db.Tables[name] = &Table{
		Name:    name,
		Columns: columns,
		Rows:    make(map[string]map[string]string),
	}
	return nil
}

func (db *Database) DeleteTable(tableName string) error {
	if _, exists := db.Tables[tableName]; !exists {
		return fmt.Errorf(" Table %s does not exist ", tableName)
	}
	delete(db.Tables, tableName)
	return nil
}
