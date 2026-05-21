# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/vrtraining-server ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/vrtraining-server /app/vrtraining-server
COPY migrations /app/migrations

ENV APP_ENV=production
ENV SERVER_PORT=8080
ENV LOG_LEVEL=info

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/vrtraining-server"]
