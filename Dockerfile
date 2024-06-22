from golang:1.21 AS build
WORKDIR /go/src
COPY ./ /go/src/

ENV CGO_ENABLED=0
WORKDIR /go/src/
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o podinate .

FROM ubuntu AS runtime
COPY --from=build /go/src/cli/podinate /usr/bin/
ENTRYPOINT ["podinate"]
