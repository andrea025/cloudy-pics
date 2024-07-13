#!/bin/bash

set -exu
set -o pipefail

# Change directory to the directory of the script
cd "$(dirname "$0")" || exit

# Archive the project
(
  cd ../..
  tar czf cloudy-pics.tar cloudy-pics
  ls -lh cloudy-pics.tar
  scp -i ~/.ssh/cloudy-pics-keys.pem -o ConnectTimeout=3 cloudy-pics.tar learnerlab:.
  rm cloudy-pics.tar
)

# Build Docker image
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab "rm -rf cloudy-pics && tar xzf cloudy-pics.tar && rm cloudy-pics.tar"
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab 'cd cloudy-pics && sudo docker build -f Dockerfile.backend -t backend:latest .'
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab 'cd cloudy-pics && sudo docker build --build-arg API_IP=$API_IP -f Dockerfile.frontend -t frontend:latest .'

# Stop remote container if running
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab "sudo docker ps -q -a --filter name=backend | xargs sudo docker rm -f" || true
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab "sudo docker ps -q -a --filter name=frontend | xargs sudo docker rm -f" || true

# Start container in daemon mode
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab "sudo docker run -d --restart unless-stopped --name frontend -p 80:80 frontend:latest"
ssh -i ~/.ssh/cloudy-pics-keys.pem learnerlab "sudo docker run -d --restart unless-stopped --name backend -v /home/ec2-user/.aws:/root/.aws -p 3000:3000 -p 4000:4000 backend:latest"
