# API-GW

Scalable API Gateway built in Go, designed for microservices. It features:

- Enhanced metrics (timestamps, response time, endpoint usage, request count)
- JWT authentication
- Rate limiting
- Docker containerization (optional)
- Modular architecture (easy to extend/modify)
- Unit testing

## Setup & Usage

### Prerequisites

- Go 1.22+
- Docker (optional)

### Running & Compile Locally

- Configuration file: api-gw.conf
```bash
go run .
go build -o api-gw . && chmod +x api-gw && ./api-gw
```
### Running on Docker
```bash
docker build -t api-gw .
docker run -p 8080:8080 api-gw
```
### Run Tests
```bash
cd tests && go test
```
### Example Request
```bash
curl -H "Authorization: Bearer <token>" http://<host>:<port>/<path>
```