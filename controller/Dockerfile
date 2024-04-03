FROM golang:1.21 AS build
WORKDIR /go/src
COPY ./controller/ ./controller/
COPY ./go.mod ./
COPY ./go.sum ./

ENV CGO_ENABLED=0
WORKDIR /go/src/controller
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o controller .

FROM ubuntu AS runtime
COPY --from=build /go/src/controller/controller /usr/bin/controller
EXPOSE 3000/tcp
ENTRYPOINT ["/bin/controller"]
