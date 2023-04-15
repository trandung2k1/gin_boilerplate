
FROM golang:1.20.3-bullseye

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /gin

EXPOSE 3000

CMD [ "/gin" ]