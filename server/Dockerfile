FROM golang:1.24.3
WORKDIR /server-app

RUN go install github.com/air-verse/air@latest

COPY server/go.mod server/go.sum ./

RUN go mod download

COPY server/.air.toml ./.air.toml

COPY server/ ./

EXPOSE ${PORT}

ENV CHOKIDAR_USEPOLLING=true
ENV WATCHPACK_POLLING=true

CMD ["/go/bin/air", "-c", ".air.toml", "--build.poll=true"]
