# Builder stage
FROM golang:1.17-alpine as builder

ENV PORT "8080"
ENV ADDR "0.0.0.0"
ENV REDIS_ADDR "localhost:6379"
ENV REDIS_PW ""
ENV API_KEY ""

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o serverbin .

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/serverbin .

EXPOSE 8080
CMD ./serverbin