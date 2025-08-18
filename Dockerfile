FROM golang:1.25-alpine AS builder
WORKDIR /app-auth

COPY ./AuthService/go.mod ./AuthService/go.sum ./
RUN go mod download

RUN mkdir -p /app-auth/backend
COPY ./AuthService/backend /app-auth/backend

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./backend/cmd/app/main.go --output /app-auth/docs

RUN mkdir -p /app-auth/bin
RUN go build -o /app-auth/bin/auth-app /app-auth/backend/cmd/app

FROM alpine:latest AS runner
WORKDIR /app-auth

ENV CONFIG_PATH=config

COPY --from=builder /app-auth/bin /app-auth/bin

RUN mkdir -p /app-auth/${CONFIG_PATH}
COPY --from=builder /app-auth/backend/internal/config /app-auth/config

CMD ["./bin/auth-app"]