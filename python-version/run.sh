#!/bin/bash

# Gmail AI Telegram Agent - Python version
# Run script

echo "Starting Gmail AI Telegram Agent (Python version)..."

# Check if config files exist
if [ ! -f "configs/config.yaml" ]; then
    echo "Error: configs/config.yaml not found!"
    echo "Please copy configs/config.example.yaml to configs/config.yaml and configure it."
    exit 1
fi

if [ ! -f "configs/prompts.yaml" ]; then
    echo "Error: configs/prompts.yaml not found!"
    echo "Please copy configs/prompts.example.yaml to configs/prompts.yaml"
    exit 1
fi

# Install dependencies if needed
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv venv
fi

echo "Activating virtual environment..."
source venv/bin/activate

echo "Installing dependencies..."
pip install -r requirements.txt

echo "Starting the agent..."
python src/main.py
