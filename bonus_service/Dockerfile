# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# ---- run stage ----
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
COPY --from=build /app/api /usr/local/bin/api
USER appuser
EXPOSE 9091
ENTRYPOINT ["api"]