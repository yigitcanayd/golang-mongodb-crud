FROM golang:1.14 as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get go.mongodb.org/mongo-driver/mongo \
&& go get go.mongodb.org/mongo-driver/mongo/options \
&& go get go.mongodb.org/mongo-driver/bson \
&& go get go.mongodb.org/mongo-driver/bson/primitive \
&& CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.9.2 as deploy

WORKDIR /root/
COPY --from=build /go/src/app .

ENTRYPOINT ["./main"]