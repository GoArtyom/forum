FROM golang:1.20.1-alpine3.16 as build
RUN apk add build-base 
WORKDIR /app 
COPY . .
RUN go build -o forum ./cmd/main.go
FROM alpine:3.16
WORKDIR /app
COPY --from=build /app /app
CMD ["./forum"]