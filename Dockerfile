FROM golang:1.19-alpine

WORKDIR /
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build -o /go-app

CMD [ "/go-app" ]