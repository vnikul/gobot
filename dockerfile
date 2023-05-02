FROM golang:1.20

WORKDIR /app

COPY . .

RUN go build -o main .

ENV PORT=8080

EXPOSE 8080

CMD ["./main"]
