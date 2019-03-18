FROM golang:1.9.2 as build
WORKDIR /go/src/github.com/iron-io/ironcli
COPY . .
RUN go get -u github.com/golang/dep/...
RUN dep ensure
RUN go build

FROM iron/busybox

COPY --from=build /go/src/github.com/iron-io/ironcli/ironcli /usr/local/bin/iron

ENTRYPOINT ["/usr/local/bin/iron"]
