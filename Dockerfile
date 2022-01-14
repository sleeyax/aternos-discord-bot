FROM golang:1.17-alpine

WORKDIR /app

# expose available environment varuables
ENV DISCORD_TOKEN=""
ENV ATERNOS_SESSION=""
ENV ATERNOS_SERVER=""

# install dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy files
COPY . ./

# build binary
RUN go build -o ./bin/aternos-discord-bot ./cmd/main.go

CMD [ "./bin/aternos-discord-bot" ]
