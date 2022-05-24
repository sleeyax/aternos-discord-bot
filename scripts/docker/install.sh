#!/usr/bin/sh
# Base installation script for Linux using Docker.
# Usage: sudo sh ./install.sh

# setup docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# download aternos-discord-bot image
docker pull sleeyax/aternos-discord-bot

# setup network
docker network create aternos-discord-bot

# post installation
echo ""
echo ""
echo "Almost done!"
echo ""
echo "You should use a DaaS (Database as a Service) provider like https://cloud.mongodb.com if you're serious about security and data persistence."
echo "If you'd rather keep it cheap and simple, execute the following command to setup a local MongoDB database on this VPS (change the username and password environment variables!):"
echo "docker run -d --name mongo --network aternos-discord-bot --restart unless-stopped -p 127.0.0.1:27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=toor mongo:5"
echo "Finally, run the bot with the following command:"
echo "docker run -d --name aternos-discord-bot --network aternos-discord-bot -e DISCORD_TOKEN=\"PASTE YOUR DISCORD TOKEN HERE\" -e MONGO_DB_URI=\"mongodb://root:toor@mongo:27017/\" sleeyax/aternos-discord-bot"
echo "Consult README.md for more details and other installation methods."
