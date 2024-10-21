FROM golang:latest

RUN mkdir /app

COPY ./src /app

RUN chmod a+x /app/scripts/*.sh

WORKDIR /app

RUN go build -o main .

# CMD ["./main"]