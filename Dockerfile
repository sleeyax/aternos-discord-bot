FROM golang:1.18-alpine

WORKDIR /app

# expose available environment variables
ENV DISCORD_TOKEN=""
ENV ATERNOS_SESSION=""
ENV ATERNOS_SERVER=""
ENV MONGO_DB_URI=""
ENV PROXY=""

# install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy files
COPY . ./

# build binary
RUN go build -o ./bin/aternos-discord-bot ./cmd/main.go

CMD [ "./bin/aternos-discord-bot" ]
