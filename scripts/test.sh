#!/bin/bash
# Simple test script for endpoints
curl -X POST http://localhost/metrics -H "Content-Type: application/json" -d '{"timestamp":"2025-12-17T00:00:00Z","cpu":50,"rps":10}'
curl http://localhost/analyze
curl http://localhost/health
