FROM golang:1.16

ENV GO111MODULE=on
ENV GOSUMDB=off

COPY . /go/src/login_server
WORKDIR /go/src/login_server

RUN go build -o login_server
CMD [ "/go/src/login_server/login_server" ]