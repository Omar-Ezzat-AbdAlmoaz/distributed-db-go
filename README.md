
# 🗄️ Distributed Database System in Go

## 📌 Overview

This is a basic **distributed database system** implemented in **Go (Golang)**.  
It demonstrates core distributed system concepts including **data replication**, **master-slave architecture**, and **fault tolerance**.

---

## 🚀 Features

- Master node creates,drop databases and tables dynamically.
- All nodes (master and slaves) can:
  - `INSERT`, `UPDATE`, `DELETE`, `SELECT`, `SEARCH` records.
- Automatic data replication to all nodes.
- Basic fault tolerance: If master fails, a slave promotes itself temporarily.
- HTTP-based communication using RESTful APIs.

---

## 🧱 Architecture                                                            
- 3+ Nodes communicating via HTTP.
- One **Master Node** (can create & drop DB , create tables and manage schema).
- Multiple **Slave Nodes** (can perform data queries).
- All nodes store data MYSql AppServer.
- Configuration is defined in `config.json`.

                ┌──────────────────────┐
                │      Client/API      │
                └─────────┬────────────┘
                          │
                          ▼
                ┌──────────────────────┐
                │      Master Node     │
                │  (Write & Read Ops)  │
                └─────────┬────────────┘
         ┌───────────────┼───────────────┐
         ▼                               ▼
┌────────────────────┐       ┌────────────────────┐
│     Slave Node 1   │       │     Slave Node 2   │
│   (Read & Replica) │       │   (Read & Replica) │
└────────────────────┘       └────────────────────┘

Replication: Master → Slaves

### 📁 Folder Structure

```
distributed-db-go/
│
├── main.go
├── go.sum
├── go.mod

```
---

## ⚙️ Setup Instructions
1. Clone the repository and navigate to the project directory.

  ### 📦 Run the Nodes
2. For the **Master Node**:
   ```bash
   go run main.go master 8080 0 0
   go run main.go slave1 8080 8080 <IPv4 Master>
   go run main.go slave2 8080 8080 <IPv4 Master>


### ✅ Prerequisites

- Go installed (version `>=1.19`)
- Terminal or VS Code
- MySQL Server (default root/rootroot)
---

## 🧪 Usage Examples
### 📌 Create database (master only)

```
curl -X POST http:// IPv4 Address:8080/execute -d '{
  "action": "create_db",
  "database": "testdb"
}'
```

### 📌 Create Table (master only)

```
curl -X POST http://IPv4 Address:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "action": "create_table",
    "database": "school",
    "table": "students",
    "columns": ["id INT PRIMARY KEY", "name VARCHAR(50)", "age INT"]
}'


```

### 📌 Insert Record (any node)

```
curl -X POST http://IPv4 Address:8080/execute -d '{
"action": "insert",
"database": "testdb",
"table": "users",
"columns": ["name", "email"],
"values": ["John Doe", "john@example.com"]
}'
```
 
### 📌 Delete Record (any node)

```
curl -X POST http://IPv4 Address:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "action": "delete",
    "database": "school",
    "table": "students",
    "where": "id = 1"
}'

```
### 📌 Select Record (any node)

```
curl -X POST http://IPv4 Address:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "action": "select",
    "database": "school",
    "table": "students",
    "where": ""
}'
```
### 📌 Update Record (any node)
```
curl -X POST http://IPv4 Address:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "database": "school",
    "table": "students",
    "columns": ["name"],
    "values": ["Alicia"],
    "where": "id = 1"
}'
```
### 📌 Drop Table ((master only))
```
curl -X POST http://192.168.107.51:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "action": "drop_table",
    "database": "school",
    "table": "students"
}'
```

### 📌 Drop Database ((master only))
```
curl -X POST http://192.168.107.51:8080/execute \
  -H "Content-Type: application/json" \
  -d '{
    "action": "drop_db",
    "database": "school"
}'
```
## 🔁 Replication

- The master (sebnd & recice) node automatically replicates write operations to the slave nodes via the `/replicate` endpoint.

---

## ⚠️ Fault Tolerance

- In the event of a master node failure, the slave nodes can be promoted to master through manual intervention.

---


## 📈 Future Improvements

- Save/load DB from files or BoltDB.
- Add Web GUI for visualization.

---

