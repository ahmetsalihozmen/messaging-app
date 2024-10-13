FROM golang:alpine AS builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev 

ENV CGO_ENABLED=1
WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o /go/bin/app ./cmd/myapp

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]

EXPOSE 8080
