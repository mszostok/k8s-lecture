# K8s lecture

Example services contains following services:
- quote - service which return random quote/tip that definitely enrich quality of your life.
- meme - generate random meme. It takes random quote from *quote* service and display it on random image.

### Simplest way
```bash
env GOOS=linux GOARCH=amd64 go build -o quote.bin ./cmd/quote/main.go
docker build -t quote .
```

TBD
