package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"distributed-db-go/database"
	"distributed-db-go/handlers"
	"distributed-db-go/utils"
)

var globalDB *database.Database

func main() {

	// تحميل إعدادات النودز
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("❌ لازم تحدد البورت كـ argument زي: go run main.go 8080")
	}
	port := os.Args[1]
	address := "localhost:" + port

	// نحدد بيانات النود الحالية
	var currentNode *utils.NodeConfig
	for _, node := range config.Nodes {
		if node.Address == address {
			currentNode = &node
			break
		}
	}

	if currentNode == nil {
		log.Fatalf("❌ مفيش نود في config.json على العنوان %s", address)
	}

	// نحدد الدور (ماستر أو سليف)
	utils.InitRoles(address, config)
	if utils.IsMaster {
		fmt.Println("🎖️ This node is MASTER")
	} else {
		fmt.Println("👷 This node is SLAVE - Monitoring master...")
		utils.StartMasterMonitor()
	}

	// نبدأ السيرفر
	startServer(address)

}

func startServer(address string) {
	// إنشاء قاعدة البيانات
	globalDB = database.NewDatabase()
	handlers.DB = globalDB // نربطها بجزء الـ Handlers

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Node is alive ✅\n👋 Welcome to the Distributed DB Node"))
	})

	http.HandleFunc("/create_table", handlers.CreateTableHandler)
	http.HandleFunc("/insert", handlers.InsertHandler)
	http.HandleFunc("/update", handlers.UpdateHandler)
	http.HandleFunc("/delete_record", handlers.DeleteRecordHandler)
	http.HandleFunc("/delete_table", handlers.DeleteTableHandler)
	http.HandleFunc("/select", handlers.SelectHandler)
	http.HandleFunc("/search", handlers.SearchHandler)

	fmt.Println("🚀 Starting server on", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("⚠️ Server failed:", err)
	}
}
