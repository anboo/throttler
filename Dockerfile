FROM golang:1.19-alpine
WORKDIR /usr/src/app
COPY . .
RUN go mod tidy
RUN go build -o main .
CMD ["/usr/src/app/main"]