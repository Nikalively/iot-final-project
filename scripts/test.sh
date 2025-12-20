#!/bin/bash
# Automated test script for all endpoints (fills window with 50 normal POST, adds anomaly, analyzes, checks health)

echo "Starting 50 normal POST requests (rps=10.0) to fill the window..."
for j in {1..5}; do
  for i in {1..10}; do
    curl -X POST http://localhost:8080/metrics -H "Content-Type: application/json" -d "{\"timestamp\":\"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\",\"cpu\":50.0,\"rps\":10.0}"
  done
done

echo "Sending 1 anomaly POST (rps=100.0)..."
curl -X POST http://localhost:8080/metrics -H "Content-Type: application/json" -d "{\"timestamp\":\"$(date -u +'%Y-%m-%dT%H:%M:%SZ')\",\"cpu\":50.0,\"rps\":100.0}"

echo "Analyzing smoothed load and anomaly count..."
curl http://localhost:8080/analyze

echo "Checking health..."
curl http://localhost:8080/health
