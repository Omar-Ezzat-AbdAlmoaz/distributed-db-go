
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
   go run main.go slave1 8081 8080 192.168.1.2
   go run main.go slave2 8082 8080 192.168.1.2


### ✅ Prerequisites

- Go installed (version `>=1.19`)
- Terminal or VS Code
- MySQL Server (default root/rootroot)
---

## 🧪 Usage Examples
### 📌 Create database (master only)

```http
POST /init_database
{
  "db_name": "DDB0",
  "user": "root",
  "password": "rootroot",
  "host": "localhost:3306"
}
```

### 📌 Create Table (master only)

```curl -X POST http:// IPv4 Address:8080/execute -d '{
  "action": "create_db",
  "database": "testdb"
}'

```

### 📌 Insert Record (any node)

```curl -X POST http://IPv4 Address:8080/execute -d '{
"action": "insert",
"database": "testdb",
"table": "users",
"columns": ["name", "email"],
"values": ["John Doe", "john@example.com"]
}'
```
 




## 🔁 Replication

- The master node automatically replicates write operations to the slave nodes via the `/replicate` endpoint.

---

## ⚠️ Fault Tolerance

- In the event of a master node failure, the slave nodes can be promoted to master through manual intervention.

---


## 📈 Future Improvements

- Save/load DB from files or BoltDB.
- Add Web GUI for visualization.

---

