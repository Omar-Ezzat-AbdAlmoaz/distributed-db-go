package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// Node represents a single node in the distributed system (master or slave).
type Node struct {
	ID       string
	Port     int
	IsMaster bool
	DB       *sql.DB
	Peers    []string
	mu       sync.Mutex
}

// Command represents an action request (create_db, insert, update, etc.) from the user.
type Command struct {
	Action   string   `json:"action"`
	Database string   `json:"database"`
	Table    string   `json:"table"`
	Columns  []string `json:"columns"`
	Values   []string `json:"values"`
	Where    string   `json:"where"`
}

func main() {
	// Ensure correct number of arguments (node ID, port, master port and IP if slave)
	if len(os.Args) < 5 {
		log.Fatal("Usage: go run main.go <node_id> <port> <master_port_if_slave> <master_ip_if_slave>")
	}

	// Read arguments for node configuration
	nodeID := os.Args[1]
	port, _ := strconv.Atoi(os.Args[2])
	masterPort := os.Args[3]
	masterIP := os.Args[4]

	// Initialize node
	node := &Node{
		ID:       nodeID,
		Port:     port,
		IsMaster: masterPort == "0",
	}

	// Connect to MySQL database
	db, err := sql.Open("mysql", "root:rootroot@tcp(127.0.0.1:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	node.DB = db

	// Set up peers (slaves) if node is master, otherwise set master IP for slave
	if node.IsMaster {
		node.Peers = []string{
			fmt.Sprintf("http://%s:%d", "192.168.1.6", port), // IP Slave 1
			fmt.Sprintf("http://%s:%d", "192.168.1.7", port), // IP Slave 2
		}
	} else {
		node.Peers = []string{fmt.Sprintf("http://%s:%s", masterIP, masterPort)}
	}

	// Ensure database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	} else {
		log.Println("Successfully connected to MySQL!")
	}

	// Set up HTTP routes for commands and replication
	http.HandleFunc("/execute", node.handleCommand)
	http.HandleFunc("/replicate", node.handleReplication)

	// Start the HTTP server to listen for incoming requests
	go func() {
		log.Printf("Node %s listening on %s:%d (Master: %v)", nodeID, getLocalIP(), port, node.IsMaster)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	// Block the main goroutine indefinitely
	select {}
}

// getLocalIP returns the local IP address of the node (used for logging purposes).
func getLocalIP() string {
	return "0.0.0.0"
}

// handleCommand processes the commands (create_db, insert, etc.) from the client.
func (n *Node) handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cmd Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result interface{}
	var err error

	switch cmd.Action {
	case "create_db":
		if !n.IsMaster {
			http.Error(w, "Only master can create databases", http.StatusForbidden)
			return
		}
		err = n.createDatabase(cmd.Database)
		result = fmt.Sprintf("Database %s created", cmd.Database)

	case "create_table":
		err = n.createTable(cmd.Database, cmd.Table, cmd.Columns)
		result = fmt.Sprintf("Table %s created in database %s", cmd.Table, cmd.Database)

	case "insert":
		err = n.insert(cmd.Database, cmd.Table, cmd.Columns, cmd.Values)
		result = fmt.Sprintf("Record inserted into %s.%s", cmd.Database, cmd.Table)

	case "select":
		result, err = n.selectData(cmd.Database, cmd.Table, cmd.Where)

	case "update":
		err = n.update(cmd.Database, cmd.Table, cmd.Columns, cmd.Values, cmd.Where)
		result = fmt.Sprintf("Record updated in %s.%s", cmd.Database, cmd.Table)

	case "delete":
		err = n.delete(cmd.Database, cmd.Table, cmd.Where)
		result = fmt.Sprintf("Record deleted from %s.%s", cmd.Database, cmd.Table)

	case "drop_table":
		err = n.dropTable(cmd.Database, cmd.Table)
		result = fmt.Sprintf("Table %s dropped from database %s", cmd.Table, cmd.Database)

	case "drop_db":
		if !n.IsMaster {
			http.Error(w, "Only master can drop databases", http.StatusForbidden)
			return
		}
		err = n.dropDatabase(cmd.Database)
		result = fmt.Sprintf("Database %s dropped", cmd.Database)

	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if n.IsMaster && isWriteOperation(cmd.Action) {
		go n.replicateToSlaves(cmd)
	}

	json.NewEncoder(w).Encode(result)
}

func isWriteOperation(action string) bool {
	return action == "create_db" || action == "create_table" ||
		action == "insert" || action == "update" ||
		action == "delete" || action == "drop_table" ||
		action == "drop_db"
}

func (n *Node) replicateToSlaves(cmd Command) {
	for _, peer := range n.Peers {
		url := fmt.Sprintf("%s/replicate", peer)
		jsonData, _ := json.Marshal(cmd)

		resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonData)))
		if err != nil {
			log.Printf("Failed to replicate to %s: %v", peer, err)
			continue
		}
		resp.Body.Close()
	}
}

func (n *Node) handleReplication(w http.ResponseWriter, r *http.Request) {
	var cmd Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	switch cmd.Action {
	case "create_db":
		err = n.createDatabase(cmd.Database)
	case "create_table":
		err = n.createTable(cmd.Database, cmd.Table, cmd.Columns)
	case "insert":
		err = n.insert(cmd.Database, cmd.Table, cmd.Columns, cmd.Values)
	case "update":
		err = n.update(cmd.Database, cmd.Table, cmd.Columns, cmd.Values, cmd.Where)
	case "delete":
		err = n.delete(cmd.Database, cmd.Table, cmd.Where)
	case "drop_table":
		err = n.dropTable(cmd.Database, cmd.Table)
	case "drop_db":
		err = n.dropDatabase(cmd.Database)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Database operations
func (n *Node) createDatabase(name string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, err := n.DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	return err
}

func (n *Node) dropDatabase(name string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, err := n.DB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	return err
}

func (n *Node) createTable(db, table string, columns []string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, err := n.DB.Exec(fmt.Sprintf("USE %s", db))
	if err != nil {
		return err
	}

	cols := strings.Join(columns, ", ")
	_, err = n.DB.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", table, cols))
	return err
}

func (n *Node) dropTable(db, table string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, err := n.DB.Exec(fmt.Sprintf("USE %s", db))
	if err != nil {
		return err
	}

	_, err = n.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
	return err
}

func (n *Node) insert(db, table string, columns, values []string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, err := n.DB.Exec(fmt.Sprintf("USE %s", db))
	if err != nil {
		return err
	}

	cols := strings.Join(columns, ", ")
	vals := strings.Join(values, "', '")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES ('%s')", table, cols, vals)
	_, err = n.DB.Exec(query)
	return err
}

func (n *Node) selectData(db, table, where string) ([]map[string]interface{}, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, err := n.DB.Exec(fmt.Sprintf("USE %s", db))
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("SELECT * FROM %s", table)
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := n.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	result := []map[string]interface{}{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			row[colName] = *val
		}

		result = append(result, row)
	}

	return result, nil
}

func (n *Node) update(db, table string, columns, values []string, where string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, err := n.DB.Exec(fmt.Sprintf("USE %s", db))
	if err != nil {
		return err
	}

	setClauses := []string{}
	for i := range columns {
		setClauses = append(setClauses, fmt.Sprintf("%s = '%s'", columns[i], values[i]))
	}

	query := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(setClauses, ", "))
	if where != "" {
		query += " WHERE " + where
	}

	_, err = n.DB.Exec(query)
	return err
}

func (n *Node) delete(db, table, where string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, err := n.DB.Exec(fmt.Sprintf("USE %s", db))
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s", table)
	if where != "" {
		query += " WHERE " + where
	}

	_, err = n.DB.Exec(query)
	return err
}
