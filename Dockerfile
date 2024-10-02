FROM golang:1.21.0

WORKDIR /go/src/app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV APP_ENV=development

COPY .env.dev /go/src/app/.env.dev

WORKDIR /go/src/app/cmd/http

CMD air