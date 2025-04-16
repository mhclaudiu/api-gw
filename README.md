# API-GW

A high-performance, scalable API Gateway built in Go, designed for microservices. It features:

- Enhanced metrics (timestamps, response time, endpoint usage, request count)
- JWT authentication
- Rate limiting
- Docker containerization (optional)
- Modular architecture (easy to extend/replace)
- Unit testing

## Setup & Usage

### Prerequisites

- Go 1.22+
- Docker (optional)

### Running Locally
```bash
go run .
```
### Running on Docker
```bash
docker build -t api-gw .
docker run -p 8080:8080 api-gw
```
### Run Tests
```bash
go run ./tests
```
### Example Request
```bash
curl -H "Authorization: Bearer <token>" http://<host>:<port>/<path>
```