## Getting Started

### Prerequisites
- Go 1.8+
- MySQL 8.0.x
- Google Cloud Service Account key in JSON "credentials.json"

### Installation or Configure
# tidy up and run
- Install the dependencies
```bash
$ go mod tidy
```
- Setup your env\
copy the `.env.example`, rename to `.env`, and define your own env
- Run the app
```bash
$ go run main.go
```