#!/usr/bin/env bash
echo aa $DOCKERHUB_USER
docker build -t powerusage:latest .
docker tag powerusage:latest $DOCKERHUB_USER/powerusage:latest
docker push $DOCKERHUB_USER/powerusage:latest
