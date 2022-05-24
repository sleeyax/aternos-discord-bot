#!/usr/bin/sh
# Installation script for Debian based Linux distributions.
# Tested with Debian 10 and 11.
apt update
apt install -y curl git
sh ./install.sh
