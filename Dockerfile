FROM golang:1.24.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /website

EXPOSE 8080

CMD ["/website"]