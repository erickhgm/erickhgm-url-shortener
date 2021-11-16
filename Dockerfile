FROM golang:1.17-alpine

# Copy GO App
WORKDIR /go/src/
COPY . app

# Build the Go app
WORKDIR /go/src/app
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Copy binary to small image
FROM alpine:3.13.6
RUN apk add ca-certificates

WORKDIR /home/url_shortener
COPY --from=0 /go/src/app/app .
COPY doc doc
COPY static static

EXPOSE 8090
CMD ["./app"]