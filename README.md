# gopherator

## Running

### With Docker

```bash
docker build -t gopherator .
docker run --rm -p 8080:8080 gopherator
```

The server will be available at http://localhost:8080.

### Locally

Requires Go 1.16+.

```bash
git clone https://github.com/yottooo/gopherator.git
cd gopherator
go mod init gopherator
go mod tidy
PORT=8080 go run .
```
