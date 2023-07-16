FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o atgo-openmensarobot cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/atgo-openmensarobot ./

ENTRYPOINT ["/app/atgo-openmensarobot"]
