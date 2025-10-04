#!/usr/bin/env python3
"""
Gmail → AI → Telegram agent (Python)
--------------------------------------------------------
Single-file prototype that:
- uses Gmail API (OAuth2) to fetch unread emails
- classifies them via OpenAI (chat completion)
- trashes unimportant emails, notifies you on Telegram about important ones
- marks important emails as read

Requirements & setup (summary):
1) Enable Gmail API in Google Cloud Console and create OAuth 2.0 Client ID (Desktop or Web).
   Download credentials.json and place in the same folder as this program.
2) Configure config.yaml with your API credentials:
      openai_api_key — your OpenAI API key
      telegram_bot_token — your Telegram bot token (BotFather)
      telegram_chat_id — chat ID to receive messages (your user id or group id)
3) Run once to get token.json (the OAuth flow will open a browser). The program will save token.json.
4) python src/main.py

Note: This is a prototype. For production you should:
- persist per-sender rules and allow user feedback
- handle quota/backoff and exponential retries
- secure credentials and token storage
- run as a service (Docker, systemd) and/or use Gmail push notifications

--------------------------------------------------------
"""

import sys
import os
import time
import re
from typing import Optional

# Add src directory to path for imports
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from config.config import load_config, validate_config
from config.prompts import load_prompts
from gmail.client import create_service
from gmail.service import GmailService
from telegram.send import send_message
from classifier.classifying import EmailClassifier, ClassifierConfig, OpenAIConfig, EmailClassificationConfig


def escape_markdown(text: str) -> str:
    """Simple markdown escape for a few characters"""
    return re.sub(r'([_*\[\]])', r'\\\1', text)


def process_inbox(classifier: EmailClassifier, gmail_service: GmailService, config, telegram_token: str, chat_id: str):
    """Process all unread messages in INBOX"""
    # List unread messages in INBOX
    messages = gmail_service.list_unread_messages()
    
    if not messages:
        print("No unread messages.")
        return
    
    print(f"Processing {len(messages)} unread messages...")
    
    for message in messages:
        process_email(classifier, gmail_service, config, telegram_token, chat_id, message)


def process_email(classifier: EmailClassifier, gmail_service: GmailService, config, telegram_token: str, chat_id: str, message: dict):
    """Process a single email"""
    message_id = message['id']
    
    # Get full message details
    msg_full = gmail_service.get_message(message_id)
    if not msg_full:
        print(f"Could not retrieve message {message_id}")
        return
    
    # Extract email details
    snippet = msg_full.get('snippet', '')
    headers = msg_full.get('payload', {}).get('headers', [])
    
    subject = gmail_service.get_header(headers, 'Subject')
    from_addr = gmail_service.get_header(headers, 'From')
    preview = snippet
    
    # Compose text for classifier
    input_text = f"From: {from_addr}\nSubject: {subject}\n\n{preview}"
    
    # Classify email
    important, reason = classifier.classify_email(input_text)
    
    if important:
        # Send Telegram notification, mark as read
        body = config.telegram.important_email_template % (
            escape_markdown(from_addr),
            escape_markdown(subject),
            escape_markdown(preview),
            escape_markdown(reason)
        )
        
        success = send_message(telegram_token, chat_id, body)
        if success:
            # Mark as read to avoid re-processing
            gmail_service.mark_as_read(message_id)
            print(f"Important email processed: {from_addr} - {subject}")
        else:
            print(f"Failed to send Telegram message for: {from_addr} - {subject}")
    else:
        # Trash the message as unimportant
        msg_about_trashed = f"🗑 Trashed message from {from_addr} subject={subject}"
        
        success = send_message(telegram_token, chat_id, msg_about_trashed)
        if success:
            # Trash the message
            if gmail_service.trash_message(message_id):
                print(msg_about_trashed)
            else:
                print(f"Failed to trash message {message_id}")
        else:
            print(f"Failed to send Telegram message about trashed email: {from_addr} - {subject}")


def create_classifier(config, prompts):
    """Create a new classifier instance"""
    openai_config = OpenAIConfig(
        api_key=config.credentials.openai_api_key,
        endpoint=config.openai.endpoint,
        model=config.openai.model,
        max_tokens=config.openai.max_tokens,
        temperature=config.openai.temperature
    )
    
    email_classification_config = EmailClassificationConfig(
        system_message=prompts.email_classification.system_message,
        user_prompt_template=prompts.email_classification.user_prompt_template
    )
    
    classifier_config = ClassifierConfig(
        openai=openai_config,
        email_classification=email_classification_config
    )
    
    return EmailClassifier(classifier_config)


def main():
    """Main application entry point"""
    config_file = "configs/config.yaml"
    
    # Load configuration from YAML
    try:
        config = load_config(config_file)
    except Exception as e:
        print(f"Unable to load config from {config_file}: {e}")
        sys.exit(1)
    
    # Validate configuration
    try:
        validate_config(config)
    except Exception as e:
        print(f"Invalid config: {e}")
        sys.exit(1)
    
    # Load prompts
    try:
        prompts = load_prompts(config.files.prompts_file)
    except Exception as e:
        print(f"Unable to load prompts: {e}")
        sys.exit(1)
    
    # Create Gmail service
    try:
        gmail_service_raw = create_service(config.files.credentials_file, config.files.token_file)
        gmail_service = GmailService(gmail_service_raw)
    except Exception as e:
        print(f"Unable to retrieve Gmail client: {e}")
        sys.exit(1)
    
    print("Agent started — polling Gmail for unread messages...")
    
    # Create classifier
    classifier = create_classifier(config, prompts)
    
    # Main polling loop
    poll_interval = config.polling.interval_seconds
    
    while True:
        try:
            process_inbox(
                classifier, 
                gmail_service, 
                config,
                config.credentials.telegram_bot_token,
                config.credentials.telegram_chat_id
            )
        except KeyboardInterrupt:
            print("\nShutting down...")
            break
        except Exception as e:
            print(f"Error processing inbox: {e}")
        
        time.sleep(poll_interval)


if __name__ == "__main__":
    main()
