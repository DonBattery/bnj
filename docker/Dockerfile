
FROM golang AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bnj .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/bnj .

COPY frontend frontend

COPY bnj_conf.yaml .

CMD ["./bnj", "run"]
