#!/bin/bash

# Ask about database persistence first
read -p "Clear existing database? (y/n) " -n 1 -r
echo
cd docker
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose down -v
else
    docker-compose down
fi

# Start fresh containers
docker-compose up -d

# Wait a bit for database to initialize
echo "Waiting for database to start..."
sleep 5

# Run the application
cd ..
go mod tidy
go run main.go