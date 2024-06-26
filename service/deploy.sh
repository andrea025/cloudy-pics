#!/bin/bash

set -exu
set -o pipefail

# Change directory to the directory of the script
cd "$(dirname "$0")" || exit

# Archive the project
(
  cd ../..
  #tar czf backend.tar wasa-photo/{service/api,service/database_nosql,cmd/webapi,Dockerfile.backend,go.mod,go.sum,service/database}
  #tar czf frontend.tar wasa-photo/{webui/public,webui/src,webui/index.html,webui/package-lock.json,webui/package.json,webui/static-resources.go,webui/vite.config.js,Dockerfile.frontend}
  tar czf wasa-photo.tar wasa-photo
  #ls -lh backend.tar
  #ls -lh frontend.tar
  ls -lh wasa-photo.tar
  #scp -o ConnectTimeout=3 backend.tar learnerlab:.
  #scp -o ConnectTimeout=3 frontend.tar learnerlab:.
  scp -o ConnectTimeout=3 wasa-photo.tar learnerlab:.
  #rm backend.tar
  #rm frontend.tar
  rm wasa-photo.tar
)

# Build Docker image
# ssh learnerlab "rm -rf wasa-photo && tar xzf backend.tar && tar xzf frontend.tar && rm backend.tar && rm frontend.tar"
ssh learnerlab "rm -rf wasa-photo && tar xzf wasa-photo.tar && rm wasa-photo.tar"
ssh learnerlab "cd wasa-photo && sudo docker build -f Dockerfile.backend -t backend:latest ."
ssh learnerlab "cd wasa-photo && sudo docker build -f Dockerfile.frontend -t frontend:latest ."

# Stop remote container if running
ssh learnerlab "sudo docker ps -q -a --filter name=backend | xargs sudo docker rm -f" || true
ssh learnerlab "sudo docker ps -q -a --filter name=frontend | xargs sudo docker rm -f" || true

# Start container in daemon mode
ssh learnerlab "sudo docker run -d --name backend --restart unless-stopped -p 3000:3000 -p 4000:4000 backend:latest"
ssh learnerlab "sudo docker run -d --name frontend --restart unless-stopped -p 80:80 frontend:latest"

# Test the connection
VM_IP="$(grep learnerlab ~/.ssh/config -A10 | grep HostName | awk '{print $2}')"
curl "http://$VM_IP/users" --silent --fail | jq
