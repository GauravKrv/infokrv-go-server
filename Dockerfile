FROM golang:1.21-alpine
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -tags netgo -ldflags '-s -w' -o app ./cmd/api
CMD ["./app"]