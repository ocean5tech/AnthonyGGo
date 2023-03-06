FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN go build -o /pricefetcher
EXPOSE 9091

CMD ["/pricefetcher"]