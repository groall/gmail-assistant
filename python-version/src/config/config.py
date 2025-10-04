"""
Configuration handling for Gmail AI Telegram Agent
"""
import os
import yaml
from typing import Dict, Any, Optional
from dataclasses import dataclass


@dataclass
class Credentials:
    """API credentials configuration"""
    openai_api_key: str
    telegram_bot_token: str
    telegram_chat_id: str


@dataclass
class Files:
    """File paths configuration"""
    credentials_file: str
    token_file: str
    prompts_file: str


@dataclass
class Polling:
    """Polling settings configuration"""
    interval_seconds: int


@dataclass
class OpenAI:
    """OpenAI API configuration"""
    endpoint: str
    model: str
    max_tokens: int
    temperature: float


@dataclass
class Telegram:
    """Telegram configuration"""
    important_email_template: str


@dataclass
class Config:
    """Main configuration class"""
    credentials: Credentials
    files: Files
    polling: Polling
    openai: OpenAI
    telegram: Telegram


def load_config(filename: str) -> Config:
    """Load configuration from YAML file"""
    try:
        with open(filename, 'r') as f:
            data = yaml.safe_load(f)
    except FileNotFoundError:
        raise FileNotFoundError(f"Config file not found: {filename}")
    except yaml.YAMLError as e:
        raise ValueError(f"Failed to parse config YAML: {e}")

    # Extract credentials
    creds_data = data.get('credentials', {})
    credentials = Credentials(
        openai_api_key=creds_data.get('openai_api_key', ''),
        telegram_bot_token=creds_data.get('telegram_bot_token', ''),
        telegram_chat_id=creds_data.get('telegram_chat_id', '')
    )

    # Extract files
    files_data = data.get('files', {})
    files = Files(
        credentials_file=files_data.get('credentials_file', ''),
        token_file=files_data.get('token_file', ''),
        prompts_file=files_data.get('prompts_file', '')
    )

    # Extract polling
    polling_data = data.get('polling', {})
    polling = Polling(
        interval_seconds=polling_data.get('interval_seconds', 60)
    )

    # Extract OpenAI
    openai_data = data.get('openai', {})
    openai = OpenAI(
        endpoint=openai_data.get('endpoint', ''),
        model=openai_data.get('model', 'gpt-3.5-turbo'),
        max_tokens=openai_data.get('max_tokens', 200),
        temperature=openai_data.get('temperature', 0.0)
    )

    # Extract Telegram
    telegram_data = data.get('telegram', {})
    telegram = Telegram(
        important_email_template=telegram_data.get('important_email_template', '')
    )

    return Config(
        credentials=credentials,
        files=files,
        polling=polling,
        openai=openai,
        telegram=telegram
    )


def validate_config(config: Config) -> None:
    """Validate configuration values"""
    if not config.files.credentials_file:
        raise ValueError("credentials_file is required in config.yaml")

    if not config.files.token_file:
        raise ValueError("token_file is required in config.yaml")

    if not config.files.prompts_file:
        raise ValueError("prompts_file is required in config.yaml")

    if not config.credentials.openai_api_key:
        raise ValueError("openai_api_key is required in config.yaml")

    if not config.credentials.telegram_bot_token:
        raise ValueError("telegram_bot_token is required in config.yaml")

    if not config.credentials.telegram_chat_id:
        raise ValueError("telegram_chat_id is required in config.yaml")

    if config.polling.interval_seconds <= 0:
        raise ValueError("interval_seconds must be greater than 0 in config.yaml")

    if config.openai.max_tokens <= 0:
        raise ValueError("max_tokens must be greater than 0 in config.yaml")

    if not (0 <= config.openai.temperature <= 1):
        raise ValueError("temperature must be between 0 and 1 in config.yaml")

    if not config.openai.endpoint:
        raise ValueError("endpoint is required in config.yaml")

    if not config.telegram.important_email_template:
        raise ValueError("important_email_template is required in config.yaml")
