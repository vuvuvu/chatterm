#!/bin/bash

# Test script to verify signal handling in the main application

echo "Testing main application signal handling..."

# Start the application in the background
echo "Starting application..."
go run main.go &
APP_PID=$!

# Wait a moment for the application to start
sleep 2

# Send SIGINT to the application
echo "Sending SIGINT to application (PID: $APP_PID)..."
kill -INT $APP_PID

# Wait for the application to exit
wait $APP_PID
EXIT_CODE=$?

echo "Application exited with code: $EXIT_CODE"

# Check if application exited gracefully
if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ Application exited gracefully with signal handling"
else
    echo "✗ Application did not exit gracefully"
fi

echo "Signal handling test completed"
