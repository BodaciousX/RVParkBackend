#!/bin/bash

# Define the target directory
TARGET_DIR="ze_backend_files"

# Remove the directory if it exists without prompting
if [ -d "$TARGET_DIR" ]; then
    echo "Removing existing $TARGET_DIR directory..."
    rm -rf "$TARGET_DIR"
fi

# Create directory
echo "Creating $TARGET_DIR directory..."
mkdir -p "$TARGET_DIR"

# Copy main files from root directory
echo "Copying root files..."
cp main.go "$TARGET_DIR/" 2>/dev/null
cp run.sh "$TARGET_DIR/" 2>/dev/null

# Function to copy files with prefix
copy_with_prefix() {
    local dir=$1
    local prefix=$2
    
    echo "Copying files from $dir..."
    
    # Check if directory exists before trying to copy
    if [ -d "$dir" ]; then
        for file in "$dir"/*.go; do
            if [ -f "$file" ]; then
                # Extract just the filename without path
                filename=$(basename "$file")
                # Copy to target with prefix
                cp "$file" "$TARGET_DIR/${prefix}${filename}" 2>/dev/null
            fi
        done
    else
        echo "Warning: Directory $dir not found, skipping..."
    fi
}

# Copy files from docker directory (non-Go files)
if [ -d "docker" ]; then
    echo "Copying files from docker..."
    for file in docker/*; do
        if [ -f "$file" ]; then
            filename=$(basename "$file")
            cp "$file" "$TARGET_DIR/d_${filename}" 2>/dev/null
        fi
    done
else
    echo "Warning: Docker directory not found, skipping..."
fi

# Copy Go files from each specified directory with appropriate prefixes
copy_with_prefix "api" "api_"
copy_with_prefix "course" "course_"
copy_with_prefix "middleware" "middleware_"
copy_with_prefix "payment" "payment_"
copy_with_prefix "space" "space_"
copy_with_prefix "tenant" "tenant_"
copy_with_prefix "user" "user_"

echo "Files copied to $TARGET_DIR successfully!"