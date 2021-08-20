FROM golang:1.16

ENV GO111MODULE=on
ENV GOSUMDB=off

COPY . /go/src/MonGoAuthBE
WORKDIR /go/src/MonGoAuthBE

RUN go build -o MonGoAuthBE
CMD [ "/go/src/MonGoAuthBE/lMonGoAuthBE" ]
