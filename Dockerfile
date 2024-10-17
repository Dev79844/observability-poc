FROM golang:1.23.1
WORKDIR /app
RUN export GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 9000
CMD ["./main"]