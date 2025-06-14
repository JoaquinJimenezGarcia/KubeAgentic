FROM golang:1.24.2 as builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o kube-agent

# Final lightweight image
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/kube-agent .

EXPOSE 8080

CMD ["./kube-agent"]
