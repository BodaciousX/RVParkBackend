#!/bin/bash

# Check if directory exists and remove it if it does
if [ -d "all_files" ]; then
    echo "all_files directory already exists"
    read -p "Would you like to remove the existing directory first? (y/n) " -n 1 -r
    echo    # Move to a new line
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        echo "Removing existing directory..."
        rm -rf all_files
    else
        echo "Operation cancelled"
        exit 1
    fi
fi

# Create directory and copy files
echo "Creating all_files directory and copying files..."
mkdir all_files && cp api/* docker/* middleware/* payment/* space/* tenant/* user/* main.go run.sh all_files/

# Check if the copy was successful
if [ $? -eq 0 ]; then
    echo "Files copied successfully!"
else
    echo "Error copying files"
    exit 1
fi

# Ask for confirmation before deletion
read -p "Would you like to delete the all_files directory? (y/n) " -n 1 -r
echo    # Move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Deleting all_files directory..."
    rm -rf all_files
    echo "Directory deleted successfully!"
else
    echo "Directory was not deleted"
fi