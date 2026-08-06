#!/usr/bin/env bash

set -e

echo "Starting Docker..."
docker compose up -d --build

echo "Waiting for API.."
sleep 20

echo "Checking health..."
curl -f http://localhost:8081/api/health

echo
echo "Checking readiness..."
curl -f http://localhost:8081/api/ready

echo
echo "Smoke Test Passed"

docker compose down