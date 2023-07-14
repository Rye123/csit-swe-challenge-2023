FROM golang:1.20

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
EXPOSE 8080
RUN make build && mv server /usr/local/bin/server

CMD ["server"]
