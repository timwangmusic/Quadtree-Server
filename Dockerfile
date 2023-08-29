FROM golang:1.21-alpine

ENV GO111MODULE=on

COPY . /app
WORKDIR /app

RUN go build -o /quadtree-server

EXPOSE 10086

CMD ["/quadtree-server"]
