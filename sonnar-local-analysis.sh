#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "Error: .env not found in current directory."
  exit 1
fi

# Verify that the variables are defined
if [ -z "$SONAR_HOST_URL" ] || [ -z "$SONAR_TOKEN" ]; then
  echo "Errror: the variables must be defined in .env file."
  exit 1
fi

# Generate test coverage report
echo "Generating test coverage report..."
go test -coverprofile=coverage.out -covermode=atomic ./...

# Generate test execution report
echo "Generating test execution report..."
go test -json ./... > test-report.json

# Run SonarQube analysis
sonar-scanner \
  -Dsonar.projectKey=market-rest-api \
  -Dsonar.sources=. \
  -Dsonar.host.url="$SONAR_HOST_URL" \
  -Dsonar.token="$SONAR_TOKEN"
