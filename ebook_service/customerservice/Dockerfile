FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y

WORKDIR /app

COPY *.go ./

RUN go mod init Customer.go
#RUN go mod tidy
RUN go env -w GO111MODULE=auto
RUN go get github.com/go-sql-driver/mysql
RUN go get -u github.com/gorilla/mux
RUN go get github.com/confluentinc/confluent-kafka-go/kafka
#RUN go get -u github.com/shopspring/decimal

RUN go build -o /Customer

EXPOSE 3000
#EXPOSE 3306

CMD [ "/Customer" ]
