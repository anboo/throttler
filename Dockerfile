FROM golang:1.19-alpine
WORKDIR /usr/src/app
COPY . .
RUN go mod tidy