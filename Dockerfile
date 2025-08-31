FROM golang:1.25.0-alpine3.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o sequence-technical-test cmd/api/main.go

FROM alpine:3.22 as runner

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chown -R appuser:appgroup /app

USER appuser

COPY --from=builder --chown=appuser:appgroup /app/sequence-technical-test .

CMD ["./sequence-technical-test"]