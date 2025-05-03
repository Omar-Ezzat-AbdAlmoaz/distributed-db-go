package handlers

import (
	"encoding/json"
	"net/http"

	"distributed-db-go/database"
	"distributed-db-go/utils"
)

type InitDatabaseRequest struct {
	DBName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"` // مثل: localhost:3306
}

func InitDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InitDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	var err error
	if utils.IsMaster {
		// الماستر ينشئ القاعدة لو مش موجودة
		err = database.ConnectOrCreateDatabase(req.User, req.Password, req.Host, req.DBName)
		if err != nil {
			http.Error(w, "❌ Master DB init failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// يرسل لباقي الـ Slaves
		for _, node := range utils.OtherNodes {
			go utils.SendPostRequest("http://"+node.Address+"/init_database", req)
		}

	} else {
		// السليف بس يتصل بقاعدة موجودة
		err = database.ConnectToExistingDatabase(req.User, req.Password, req.Host, req.DBName)
		if err != nil {
			http.Error(w, "❌ Slave DB connect failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("✅ Database connected: " + req.DBName))
}

type CreateTableRequest struct {
	TableName string   `json:"table_name"`
	Columns   []string `json:"columns"`
}

func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := database.CreateTable(req.TableName, req.Columns); err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("✅ Table created successfully"))

	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
			go utils.SendPostRequest("http://"+node.Address+"/create_table", req)
		}
	}
}

type InsertRequest struct {
	TableName string            `json:"table_name"`
	Data      map[string]string `json:"data"`
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req InsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := database.Insert(req.TableName, req.Data); err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("✅ Inserted successfully"))

	if utils.IsMaster {
		for _, node := range utils.OtherNodes {
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := database.Update(req.TableName, req.RowID, req.NewData); err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("✅ Updated successfully"))

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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := database.DeleteRow(req.TableName, req.RowID); err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("✅ Row deleted"))

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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := database.DeleteTable(req.TableName); err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
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
	table := r.URL.Query().Get("table")
	if table == "" {
		http.Error(w, "❌ Missing table name", http.StatusBadRequest)
		return
	}

	records, err := database.GetAll(table)
	if err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(records)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	column := r.URL.Query().Get("column")
	value := r.URL.Query().Get("value")

	if table == "" || column == "" || value == "" {
		http.Error(w, "❌ Missing parameters", http.StatusBadRequest)
		return
	}

	records, err := database.SearchByColumn(table, column, value)
	if err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(records)
}

func DropDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.IsMaster {
		http.Error(w, "❌ Only master can drop database", http.StatusForbidden)
		return
	}

	err := database.DropDatabase("distributed_db") // replace with actual name
	if err != nil {
		http.Error(w, "❌ "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("🧨 Database dropped"))

	for _, node := range utils.OtherNodes {
		go utils.SendPostRequest("http://"+node.Address+"/drop_database", nil)
	}
}
