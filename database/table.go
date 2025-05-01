package database

import (
	"fmt"
)

// جدول يحتوي على الأعمدة (Columns) و السجلات (Rows)
type Table struct {
	Name    string
	Columns []string
	Rows    map[string]map[string]string // RowID -> (Column -> Value)
}

// إضافة سجل جديد للجدول
func (t *Table) Insert(rowID string, data map[string]string) error {
	// تأكد من الأعمدة
	for _, col := range t.Columns {
		if _, ok := data[col]; !ok {
			return fmt.Errorf(" Column %s is missing in the data", col)
		}
	}
	t.Rows[rowID] = data
	return nil
}

func (t *Table) Update(rowID string, newData map[string]string) error {
	row, exists := t.Rows[rowID]
	if !exists {
		return fmt.Errorf(" Row %s does not exist ", rowID)
	}
	for k, v := range newData {
		row[k] = v
	}
	return nil
}

func (t *Table) DeleteRow(rowID string) error {
	if _, exists := t.Rows[rowID]; !exists {
		return fmt.Errorf(" Row %s does not exist", rowID)
	}
	delete(t.Rows, rowID)
	return nil
}

func (t *Table) GetAll() []map[string]string {
	var results []map[string]string
	for _, row := range t.Rows {
		results = append(results, row)
	}
	return results
}

func (t *Table) SearchByColumn(column string, value string) []map[string]string {
	var results []map[string]string
	for _, row := range t.Rows {
		if row[column] == value {
			results = append(results, row)
		}
	}
	return results
}
