FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /proxy-go

EXPOSE 8080

CMD [ "/proxy-go" ]
