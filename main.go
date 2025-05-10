package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"distributed-db-go/handlers"
	"distributed-db-go/utils"
)

func main() {

	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	if len(os.Args) < 2 {
		log.Fatal(" لازم تحدد البورت كـ argument زي: go run main.go 8080")
	}
	port := os.Args[1]
	address := "localhost:" + port

	var currentNode *utils.NodeConfig
	for _, node := range config.Nodes {
		if node.Address == address {
			currentNode = &node
			break
		}
	}
	if currentNode == nil {
		log.Fatalf("مفيش نود في config.json على العنوان %s", address)
	}

	utils.InitRoles(address, config)
	if utils.IsMaster {
		fmt.Println("This node is MASTER")
	} else {
		fmt.Println("This node is SLAVE - Monitoring master...")
		utils.StartMasterMonitor()
	}

	// // الاتصال بـ MySQL
	// database.ConnectOrCreateDatabase("root", "rootroot", "localhost:3306", "distributed_db") // عدّل الـ credentials حسب بيئتك

	startServer(address)
}

func startServer(address string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Node is alive \n Welcome to the Distributed DB Node"))
	})
	http.HandleFunc("/notify", handlers.NotifyHandler)

	http.HandleFunc("/init_database", handlers.InitDatabaseHandler)
	http.HandleFunc("/create_table", handlers.CreateTableHandler)
	http.HandleFunc("/insert", handlers.InsertHandler)
	http.HandleFunc("/update", handlers.UpdateHandler)
	http.HandleFunc("/delete_record", handlers.DeleteRecordHandler)
	http.HandleFunc("/delete_table", handlers.DeleteTableHandler)
	http.HandleFunc("/select", handlers.SelectHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/drop_database", handlers.DropDatabaseHandler)

	fmt.Println("Starting server on", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
