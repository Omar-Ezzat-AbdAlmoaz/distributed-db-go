
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
- One **Master Node** (can create tables and manage schema).
- Multiple **Slave Nodes** (can perform data queries).
- All nodes store data in-memory independently.
- Configuration is defined in `config.json`.

---

## ⚙️ Setup Instructions

### ✅ Prerequisites

- Go installed (version `>=1.19`)
- Terminal or VS Code

### 📁 Folder Structure

```
distributed-db-go/
│
├── main.go
├── config.json
├── go.mod
│
├── /handlers        # HTTP API Handlers
├── /database        # Table & Data logic
├── /utils           # Config, networking, monitoring
```

### 📦 Run the Nodes

Each terminal runs a different node:

```bash
go run main.go 8080   # Master
go run main.go 8081   # Slave 1
go run main.go 8082   # Slave 2
```

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

```http
POST /create_table
{
  "table_name": "students",
  "columns": ["id", "name", "grade"]
}
```

### 📌 Insert Record (any node)

```http
POST /insert
{
  "table_name": "students",
  "row_id": "1",
  "data": {
    "id": "1",
    "name": "Ahmed",
    "grade": "A"
  }
}
```

### 📌 Update

```http
POST /update
{
  "table_name": "students",
  "row_id": "1",
  "new_data": {
    "grade": "A+"
  }
}
```

### 📌 Select All

```http
GET /select?table=students
```

### 📌 Search

```http
GET /search?table=students&column=name&value=Ahmed
```

---

### 📌 Delete Record

```http
POST /delete_record
{
  "table_name": "students",
  "row_id": "1"
}
```

### 📌 Delete Table (master only)

```http
POST /delete_table
{
  "table_name": "students"
}
```

## 🔁 Replication

- All insert/update/delete/create_table actions are automatically forwarded from master to all slaves.
- Slaves handle requests independently once data is replicated.

---

## ⚠️ Fault Tolerance

- If master node fails (ping unreachable), a slave promotes itself temporarily as a new master.
- Simple logic using node priority (`node.id`) is used for leader fallback.

---

## 📌 Limitations

- In-memory only (data lost after restart).
- No persistent storage (yet).
- No full consensus (e.g., Raft or Paxos).

---

## 📈 Future Improvements

- Save/load DB from files or BoltDB.
- Implement full leader election (Bully/Raft).
- Add Web GUI for visualization.


