"""
Gmail service wrapper for email operations
"""
from typing import List, Dict, Any, Optional
from googleapiclient.discovery import Resource
from googleapiclient.errors import HttpError


class GmailService:
    """Gmail service wrapper class"""
    
    def __init__(self, service: Resource):
        self.service = service
    
    def list_unread_messages(self) -> List[Dict[str, Any]]:
        """List unread messages in inbox"""
        try:
            results = self.service.users().messages().list(
                userId='me', 
                q='is:unread in:inbox'
            ).execute()
            
            messages = results.get('messages', [])
            return messages
        except HttpError as error:
            print(f"Error listing messages: {error}")
            return []
    
    def get_message(self, message_id: str) -> Optional[Dict[str, Any]]:
        """Get full message details"""
        try:
            message = self.service.users().messages().get(
                userId='me', 
                id=message_id, 
                format='full'
            ).execute()
            return message
        except HttpError as error:
            print(f"Error getting message {message_id}: {error}")
            return None
    
    def mark_as_read(self, message_id: str) -> bool:
        """Mark message as read"""
        try:
            self.service.users().messages().modify(
                userId='me',
                id=message_id,
                body={'removeLabelIds': ['UNREAD']}
            ).execute()
            return True
        except HttpError as error:
            print(f"Error marking message {message_id} as read: {error}")
            return False
    
    def trash_message(self, message_id: str) -> bool:
        """Move message to trash"""
        try:
            self.service.users().messages().trash(
                userId='me',
                id=message_id
            ).execute()
            return True
        except HttpError as error:
            print(f"Error trashing message {message_id}: {error}")
            return False
    
    def get_header(self, headers: List[Dict[str, str]], name: str) -> str:
        """Extract header value from message headers"""
        for header in headers:
            if header.get('name', '').lower() == name.lower():
                return header.get('value', '')
        return ''
