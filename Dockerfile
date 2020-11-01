FROM golang:latest AS builder
ENV GO111MODULE=on

# Download modules
WORKDIR $GOPATH/src/github.com/ptcoffee/image-server
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /image-server .

# Run
FROM alpine:latest
RUN apk add --no-cache pngquant
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /image-server ./
EXPOSE 5000
ENTRYPOINT ["./image-server"]
