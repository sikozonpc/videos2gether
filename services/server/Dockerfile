FROM golang:1.17-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
## pull in any dependencies
RUN go mod download
## project will now successfully build with the necessary go libraries included.
RUN go build -o main .
## newly created binary executable
CMD ["/app/main"]