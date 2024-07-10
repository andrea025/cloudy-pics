#!/bin/bash

set -exu
set -o pipefail

# Change directory to the directory of the script
cd "$(dirname "$0")" || exit

# Archive the project
(
  cd ../..
  #tar czf backend.tar cloudy-pics/{service/api,service/database_nosql,cmd/webapi,Dockerfile.backend,go.mod,go.sum,service/database}
  #tar czf frontend.tar cloudy-pics/{webui/public,webui/src,webui/index.html,webui/package-lock.json,webui/package.json,webui/static-resources.go,webui/vite.config.js,Dockerfile.frontend}
  tar czf cloudy-pics.tar cloudy-pics
  #ls -lh backend.tar
  #ls -lh frontend.tar
  ls -lh cloudy-pics.tar
  #scp -o ConnectTimeout=3 backend.tar learnerlab:.
  #scp -o ConnectTimeout=3 frontend.tar learnerlab:.
  scp -i ~/.ssh/labuser.pem -o ConnectTimeout=3 cloudy-pics.tar learnerlab:.
  #rm backend.tar
  #rm frontend.tar
  rm cloudy-pics.tar
)

# Build Docker image
# ssh learnerlab "rm -rf cloudy-pics && tar xzf backend.tar && tar xzf frontend.tar && rm backend.tar && rm frontend.tar"
ssh -i ~/.ssh/labuser.pem learnerlab "rm -rf cloudy-pics && tar xzf cloudy-pics.tar && rm cloudy-pics.tar"
ssh -i ~/.ssh/labuser.pem learnerlab "cd cloudy-pics && sudo docker build -f Dockerfile.backend -t backend:latest ."
ssh -i ~/.ssh/labuser.pem learnerlab "cd cloudy-pics && sudo docker build --build-arg API_IP=18.232.58.105 -f Dockerfile.frontend -t frontend:latest ."

# Stop remote container if running
ssh -i ~/.ssh/labuser.pem learnerlab "sudo docker ps -q -a --filter name=backend | xargs sudo docker rm -f" || true
ssh -i ~/.ssh/labuser.pem learnerlab "sudo docker ps -q -a --filter name=frontend | xargs sudo docker rm -f" || true

# Start container in daemon mode
ssh -i ~/.ssh/labuser.pem learnerlab "sudo docker run -d --name frontend -p 80:80 frontend:latest"
ssh -i ~/.ssh/labuser.pem learnerlab "sudo docker run -d --name backend -v /home/ec2-user/.aws:/root/.aws -p 3000:3000 -p 4000:4000 backend:latest"

# Test the connection
VM_IP="$(grep learnerlab ~/.ssh/config -A10 | grep HostName | awk '{print $2}')"
curl "http://$VM_IP/users" --silent --fail | jq
