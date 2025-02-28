FROM golang:1.23.5-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

WORKDIR /app/cmd/backend

RUN go build -o main .

EXPOSE 8080

CMD ["sh", "-c", "sleep 5 && /app/cmd/backend/main"]
