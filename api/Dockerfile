# --- Etape 1 : compilation (image Go complete) ---
FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/cron ./cron

# --- Etape 2 : image finale minimale ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=build /out/api  /app/api
COPY --from=build /out/cron /app/cron

EXPOSE 8086
CMD ["/app/api"]
