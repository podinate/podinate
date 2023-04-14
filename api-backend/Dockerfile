FROM golang:1.20.3 AS build
WORKDIR /go/src/github.com/johncave/podinate/api-backend
COPY ./api-backend/ ./
COPY ./go.* /go/src/github.com/johncave/podinate/
# RUN ls -lah /go/src/github.com/johncave/podinate
# RUN ls -lah

ENV CGO_ENABLED=0
RUN go get

RUN go build -a -installsuffix cgo -o restapi .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/github.com/johncave/podinate/api-backend/restapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./restapi"]
