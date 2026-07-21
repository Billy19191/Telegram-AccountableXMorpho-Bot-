FROM golang:tip-alpine3.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build ./cmd/main.go

EXPOSE 3000
CMD ["./main"]


