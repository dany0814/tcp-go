FROM    golang:1.20.2-alpine3.17 AS dependencies
WORKDIR /opt/app-code
COPY    ./go.mod .
RUN     go mod download

FROM    golang:1.20.2-alpine3.17 AS builder
COPY    ./ /opt/app-code/
COPY    --from=dependencies /opt/app-code/ /opt/app-code/
WORKDIR /opt/app-code/
RUN     go build -v -o ./dist/app ./server/main.go

FROM    alpine:3.17 AS dist
COPY    --from=builder /opt/app-code/dist/. /opt/app-code
WORKDIR /opt/app-code
CMD     ["./app"]