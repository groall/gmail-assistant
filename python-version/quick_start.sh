#!/bin/bash

# Gmail AI Agent - Quick Start Script for Colleagues
# This script sets up everything needed to run the Gmail AI Agent

echo "🚀 Gmail AI Agent - Quick Start"
echo "================================"

# Check Python version
python_version=$(python3 --version 2>&1 | cut -d' ' -f2 | cut -d'.' -f1,2)
required_version="3.8"

if [ "$(printf '%s\n' "$required_version" "$python_version" | sort -V | head -n1)" != "$required_version" ]; then
    echo "❌ Error: Python 3.8+ required, found $python_version"
    exit 1
fi

echo "✅ Python version: $python_version"

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "📦 Creating virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source venv/bin/activate

# Install dependencies
echo "📥 Installing dependencies..."
pip install --upgrade pip
pip install -r requirements.txt

# Check if config files exist
if [ ! -f "configs/config.yaml" ]; then
    echo "⚙️  Setting up configuration..."
    if [ -f "configs/config.example.yaml" ]; then
        cp configs/config.example.yaml configs/config.yaml
        echo "📝 Created configs/config.yaml from example"
        echo "⚠️  Please edit configs/config.yaml with your credentials!"
    else
        echo "❌ Error: config.example.yaml not found!"
        exit 1
    fi
fi

if [ ! -f "configs/prompts.yaml" ]; then
    echo "⚙️  Setting up prompts..."
    if [ -f "configs/prompts.example.yaml" ]; then
        cp configs/prompts.example.yaml configs/prompts.yaml
        echo "📝 Created configs/prompts.yaml from example"
    else
        echo "❌ Error: prompts.example.yaml not found!"
        exit 1
    fi
fi

# Test imports
echo "🧪 Testing imports..."
python test_imports.py

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Setup complete!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Edit configs/config.yaml with your API credentials"
    echo "2. Add your Gmail credentials JSON to configs/"
    echo "3. Run: python src/main.py"
    echo ""
    echo "📖 For detailed setup instructions, see README.md"
else
    echo "❌ Import test failed. Please check your Python environment."
    exit 1
fi

