"""
Gmail OAuth2 client handling
"""
import json
import os
from typing import Optional
from google.auth.transport.requests import Request
from google.oauth2.credentials import Credentials
from google_auth_oauthlib.flow import InstalledAppFlow
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError


# Gmail API scopes
SCOPES = ['https://www.googleapis.com/auth/gmail.modify']


def get_credentials(credentials_file: str, token_file: str) -> Credentials:
    """
    Get Gmail API credentials, handling OAuth2 flow if needed.
    """
    creds = None
    
    # Load existing token if available
    if os.path.exists(token_file):
        try:
            creds = Credentials.from_authorized_user_file(token_file, SCOPES)
        except Exception as e:
            print(f"Error loading existing token: {e}")
            creds = None
    
    # If there are no (valid) credentials available, let the user log in
    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            try:
                creds.refresh(Request())
            except Exception as e:
                print(f"Error refreshing token: {e}")
                creds = None
        
        if not creds:
            try:
                flow = InstalledAppFlow.from_client_secrets_file(
                    credentials_file, SCOPES)
                creds = flow.run_local_server(port=0)
            except Exception as e:
                raise Exception(f"Failed to get credentials: {e}")
        
        # Save the credentials for the next run
        try:
            with open(token_file, 'w') as token:
                token.write(creds.to_json())
            print(f"Credentials saved to: {token_file}")
        except Exception as e:
            print(f"Warning: Could not save credentials: {e}")
    
    return creds


def create_service(credentials_file: str, token_file: str):
    """
    Create Gmail service instance.
    """
    try:
        creds = get_credentials(credentials_file, token_file)
        service = build('gmail', 'v1', credentials=creds)
        return service
    except Exception as e:
        raise Exception(f"Unable to create Gmail service: {e}")
