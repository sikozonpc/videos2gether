FROM golang:1.17-alpine as builder

ENV PORT "8080"
ENV ADDR "0.0.0.0"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080
# CMD ./server
CMD ["sh", "-c", "./server"]