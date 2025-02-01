FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y

WORKDIR /app

COPY *.go ./

RUN go mod init datasync.go
#RUN go mod tidy
RUN go env -w GO111MODULE=auto
RUN go get github.com/go-sql-driver/mysql
RUN go get -u github.com/gorilla/mux
#RUN go get -u github.com/shopspring/decimal
RUN go get go.mongodb.org/mongo-driver/mongo
RUN go get go.mongodb.org/mongo-driver/bson

RUN go build -o /datasync

EXPOSE 3000
#EXPOSE 3306

CMD [ "/DataSync" ]
