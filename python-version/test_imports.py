#!/usr/bin/env python3
"""
Test script to verify all imports work correctly
"""
import sys
import os

# Add src directory to path
sys.path.append(os.path.join(os.path.dirname(__file__), 'src'))

try:
    from config.config import load_config, validate_config
    print("✓ Config module imported successfully")
except ImportError as e:
    print(f"✗ Config module import failed: {e}")

try:
    from config.prompts import load_prompts
    print("✓ Prompts module imported successfully")
except ImportError as e:
    print(f"✗ Prompts module import failed: {e}")

try:
    from gmail.client import create_service
    from gmail.service import GmailService
    print("✓ Gmail modules imported successfully")
except ImportError as e:
    print(f"✗ Gmail modules import failed: {e}")

try:
    from telegram.send import send_message
    print("✓ Telegram module imported successfully")
except ImportError as e:
    print(f"✗ Telegram module import failed: {e}")

try:
    from classifier.classifying import EmailClassifier
    print("✓ Classifier module imported successfully")
except ImportError as e:
    print(f"✗ Classifier module import failed: {e}")

print("\nAll imports completed!")
