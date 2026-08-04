#!/usr/bin/env bash

set -e

echo "Starting Docker..."
docker compose up -d --build

echo "Waiting for API.."
sleep 5

echo "Checking health..."
curl -f http://localhost:8080/api/health

echo
echo "Checking readiness..."
curl -f http://localhost:8080/api/ready

echo
echo "Smoke Test Passed"

docker compose down