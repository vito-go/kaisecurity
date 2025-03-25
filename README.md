# Kai Security - Vulnerability Scan API

This is a Go-based REST API service that provides two main endpoints:

- `POST /scan`: Fetches and stores vulnerability scan data from a GitHub repository.
- `POST /query`: Queries stored vulnerabilities by severity.

This service uses SQLite for storage and supports parallel file processing.

---

## ✅ Features

- Concurrent scanning of up to 3 JSON files from GitHub
- Stores metadata and vulnerability information in SQLite
- Query vulnerabilities by severity
- Written in Go, tested with 60%+ code coverage
- Dockerized and production-ready

---

## 🚀 How to Run

### 1. Run locally with Go:
```bash
go run cmd/main.go -db=./kai_security.db -port=8080
```

### 2. Run with Docker:
```bash
# Build the image
docker build -t kai_security:latest .

# Run the container
docker run -p 8080:8080 kai_security:latest
```

---

## 🔬 Testing Instructions

### 1. Run all unit tests:
```bash
make test
```

### 2. View coverage report:
```bash
make cover
```

### 3. Run CI-friendly test (verbose, without coverage):
```bash
make ci
```

---

## 📦 API Examples

### POST /scan
```bash
curl -X POST http://localhost:8080/scan \
  -H 'Content-Type: application/json' \
  -d '{
    "repo": "https://github.com/velancio/vulnerability_scans",
    "files": ["vulnscan15.json"]
  }'
```

### POST /query
```bash
curl -X POST http://localhost:8080/query \
  -H 'Content-Type: application/json' \
  -d '{
    "filters": {
      "severity": "HIGH"
    }
  }'
```

---

## 📁 Project Structure (Key Directories)

```
.
├── cmd/             # Application entry point
├── internal/        # App logic,httpsrv, handlers, DAO, models
├── pkg/             # Reusable utilities (db, utils)
├── Makefile         # Build and test commands
├── Dockerfile       # Docker image definition
├── README.md        # Project documentation
```

---

## 📄 License

- TODO: Add license
```
