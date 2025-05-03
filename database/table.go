package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// إدخال صف جديد
func Insert(table string, data map[string]string) error {
	cols := []string{}
	vals := []string{}
	args := []interface{}{}

	for k, v := range data {
		cols = append(cols, k)
		vals = append(vals, "?")
		args = append(args, v)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(vals, ", "))
	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("❌ Insert failed: %v", err)
	}
	return nil
}

// تحديث صف
func Update(table, rowID string, newData map[string]string) error {
	set := []string{}
	args := []interface{}{}

	for k, v := range newData {
		set = append(set, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	args = append(args, rowID)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", table, strings.Join(set, ", "))
	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("❌ Update failed: %v", err)
	}
	return nil
}

// حذف صف
func DeleteRow(table, rowID string) error {
	_, err := DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), rowID)
	if err != nil {
		return fmt.Errorf("❌ Delete failed: %v", err)
	}
	return nil
}

// استرجاع كل الصفوف
func GetAll(table string) ([]map[string]string, error) {
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var result []map[string]string

	for rows.Next() {
		values := make([]sql.RawBytes, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		rows.Scan(scanArgs...)

		row := map[string]string{}
		for i, col := range cols {
			row[col] = string(values[i])
		}
		result = append(result, row)
	}
	return result, nil
}

// البحث في جدول حسب قيمة عمود
func SearchByColumn(table, column, value string) ([]map[string]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", table, column)
	rows, err := DB.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var result []map[string]string

	for rows.Next() {
		values := make([]sql.RawBytes, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		rows.Scan(scanArgs...)

		row := map[string]string{}
		for i, col := range cols {
			row[col] = string(values[i])
		}
		result = append(result, row)
	}
	return result, nil
}
