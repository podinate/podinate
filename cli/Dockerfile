from golang:1.21 AS build
WORKDIR /go/src
COPY ./cli/ ./cli/
COPY ./lib/ ./lib/
COPY ./go.mod ./
COPY ./go.sum ./

ENV CGO_ENABLED=0
WORKDIR /go/src/cli
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o podinate .

FROM ubuntu AS runtime
COPY --from=build /go/src/cli/podinate /usr/local/bin/
ENTRYPOINT ["podinate"]
