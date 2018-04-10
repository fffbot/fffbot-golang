#!/bin/bash -e

docker run -d --restart unless-stopped --name fffbot --env-file .grawenv fffbot/fffbot:1.0
