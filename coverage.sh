#!/bin/bash
set -e

go test ./... -coverprofile=reports/coverage.out
total=$(go tool cover -func=reports/coverage.out | grep total: | awk '{print substr($3, 1, length($3)-1)}')

limit=50.0

# Check if the total coverage is less than the limit
if (( $(echo "$total < $limit" | bc -l) )); then
  echo "Build failed: test coverage ($total%) is below the required threshold of $limit%"
  exit 1
else
  echo "Build passed: test coverage ($total%) meets the required threshold of $limit%"
fi
