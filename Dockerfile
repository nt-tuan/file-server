FROM golang:latest AS builder
ENV GO111MODULE=on

# Download modules
WORKDIR $GOPATH/src/github.com/thanhtuan260593/file-server
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /file-server .

# Run
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /file-server ./
EXPOSE 5000
ENTRYPOINT ["./file-server"]
