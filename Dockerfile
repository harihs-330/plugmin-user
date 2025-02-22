# builder image
FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum /build/
COPY . /build/

RUN go mod download

RUN go build -o /build/main /build/main.go

# generate clean, final image for end users
FROM alpine:latest

COPY --from=builder /build/main /app/
COPY --from=builder /build/migrations /app/migrations
COPY --from=builder /build/templates /app/templates

WORKDIR /app

RUN chmod +x /app/main

EXPOSE 8080

# executable
CMD [ "./main"]