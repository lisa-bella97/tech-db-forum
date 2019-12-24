#!/bin/bash

docker stop $(docker ps -a -q)
docker rm $(docker ps -a -q)
easyjson -all -pkg app/models/
docker build -t tech-db-forum -f Dockerfile .
docker run -p 5000:5000 -t tech-db-forum
