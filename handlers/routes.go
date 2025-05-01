package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"distributed-db-go/database"
	"distributed-db-go/utils"
)

var DB *database.Database // هنربطها من main.go

type CreateTableRequest struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
}

// POST /create_table
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error in the data sent ", http.StatusBadRequest)
		return
	}

	err = DB.CreateTable(req.TableName, req.Columns)
	if err != nil {
		http.Error(w, fmt.Sprintf("❌ Table creation failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ The table was created successfully "))

	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
			go utils.SendPostRequest("http://"+node.Address+"/create_table", req)
		}
	}

}

type InsertRequest struct {
	TableName string            `json:"table_name"`
	RowID     string            `json:"row_id"`
	Data      map[string]string `json:"data"`
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InsertRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "❌ Format error ", http.StatusBadRequest)
		return
	}

	table, exists := DB.Tables[req.TableName]
	if !exists {
		http.Error(w, "❌ Table not found ", http.StatusNotFound)
		return
	}

	err = table.Insert(req.RowID, req.Data)
	if err != nil {
		http.Error(w, "❌ Input failure: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Data entry was successful "))

	// لو النود الحالي Master، ابعت البيانات للنودز التانية  🔁 Replicate للـ Slaves لو Master
	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
			// نجهز نسخة من نفس الـ request ونبعتها لباقي النودز
			go utils.SendPostRequest("http://"+node.Address+"/insert", req)
		}
	}
}

type UpdateRequest struct {
	TableName string            `json:"table_name"`
	RowID     string            `json:"row_id"`
	NewData   map[string]string `json:"new_data"`
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "❌ Format error", http.StatusBadRequest)
		return
	}

	table, exists := DB.Tables[req.TableName]
	if !exists {
		http.Error(w, "❌ Table not found", http.StatusNotFound)
		return
	}

	err = table.Update(req.RowID, req.NewData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte("✅ Record updated"))

	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
			go utils.SendPostRequest("http://"+node.Address+"/update", req)
		}
	}
}

type DeleteRecordRequest struct {
	TableName string `json:"table_name"`
	RowID     string `json:"row_id"`
}

func DeleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteRecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "❌ Format error", http.StatusBadRequest)
		return
	}

	table, exists := DB.Tables[req.TableName]
	if !exists {
		http.Error(w, "❌ Table not found", http.StatusNotFound)
		return
	}

	err = table.DeleteRow(req.RowID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte("✅ Record deleted"))

	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
			go utils.SendPostRequest("http://"+node.Address+"/delete_record", req)
		}
	}
}

type DeleteTableRequest struct {
	TableName string `json:"table_name"`
}

func DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteTableRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "❌ Format error", http.StatusBadRequest)
		return
	}

	err = DB.DeleteTable(req.TableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte("✅ Table deleted"))

	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
			go utils.SendPostRequest("http://"+node.Address+"/delete_table", req)
		}
	}
}

func SelectHandler(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	if tableName == "" {
		http.Error(w, "❌ برجاء تحديد اسم الجدول", http.StatusBadRequest)
		return
	}

	table, exists := DB.Tables[tableName]
	if !exists {
		http.Error(w, "❌ Table not found", http.StatusNotFound)
		return
	}

	results := table.GetAll()
	json.NewEncoder(w).Encode(results)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	column := r.URL.Query().Get("column")
	value := r.URL.Query().Get("value")

	if tableName == "" || column == "" || value == "" {
		http.Error(w, "❌ برجاء تحديد table, column و value", http.StatusBadRequest)
		return
	}

	table, exists := DB.Tables[tableName]
	if !exists {
		http.Error(w, "❌ Table not found", http.StatusNotFound)
		return
	}

	results := table.SearchByColumn(column, value)
	json.NewEncoder(w).Encode(results)
}
