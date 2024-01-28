FROM golang:1.21

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /go-todo ./cmd/main.go

EXPOSE 3000

CMD [ "/go-todo" ]
