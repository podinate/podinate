FROM golang:1.21 AS build
WORKDIR /go/src
COPY ./api-backend/ ./api-backend/
COPY ./go.mod ./
COPY ./go.sum ./

ENV CGO_ENABLED=0
WORKDIR /go/src/api-backend
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o controller .

FROM ubuntu AS runtime
COPY --from=build /go/src/api-backend/controller /usr/bin/controller
EXPOSE 3000/tcp
ENTRYPOINT ["/bin/controller"]
