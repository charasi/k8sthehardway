FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y

WORKDIR /app

COPY *.go ./

RUN go mod init BookStoreMobileBFF.go
#RUN go mod tidy
RUN go env -w GO111MODULE=auto
RUN go get github.com/go-sql-driver/mysql
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/shopspring/decimal
RUN go get github.com/golang-jwt/jwt

RUN go build -o /BookStoreMobileBFF

EXPOSE 8080
#EXPOSE 3306

CMD [ "/BookStoreMobileBFF" ]
