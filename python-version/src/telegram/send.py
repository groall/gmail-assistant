"""
Telegram message sending functionality
"""
import requests
from typing import Optional


class TelegramClient:
    """Telegram client for sending messages"""
    
    def __init__(self, bot_token: str):
        self.bot_token = bot_token
        self.base_url = f"https://api.telegram.org/bot{bot_token}"
    
    def send_message(self, chat_id: str, text: str, parse_mode: str = "Markdown") -> bool:
        """
        Send a message to a Telegram chat.
        
        Args:
            chat_id: Telegram chat ID
            text: Message text
            parse_mode: Message parse mode (Markdown or HTML)
            
        Returns:
            bool: True if message sent successfully, False otherwise
        """
        url = f"{self.base_url}/sendMessage"
        
        payload = {
            "chat_id": chat_id,
            "text": text,
            "parse_mode": parse_mode
        }
        
        try:
            response = requests.post(url, json=payload, timeout=10)
            response.raise_for_status()
            return True
        except requests.exceptions.RequestException as e:
            print(f"Failed to send Telegram message: {e}")
            return False


def send_message(bot_token: str, chat_id: str, text: str) -> bool:
    """
    Convenience function to send a Telegram message.
    
    Args:
        bot_token: Telegram bot token
        chat_id: Telegram chat ID
        text: Message text
        
    Returns:
        bool: True if message sent successfully, False otherwise
    """
    client = TelegramClient(bot_token)
    return client.send_message(chat_id, text)
