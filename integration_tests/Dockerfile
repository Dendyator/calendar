FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go mod tidy

CMD ["go", "test", "-v", "./integration_tests"]
