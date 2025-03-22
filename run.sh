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

echo "Waiting for database to be ready..."
max_retries=30
count=0
while [ $count -lt $max_retries ]; do
    if docker exec RVParkDB pg_isready -U dbuser > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting for database to start... ($((count+1))/$max_retries)"
    sleep 2
    count=$((count+1))
done

if [ $count -eq $max_retries ]; then
    echo "Error: Database did not start within the expected time"
    exit 1
fi

# Run the application
cd ..
go mod tidy
go run main.go