"""
Prompts configuration handling for Gmail AI Telegram Agent
"""
import yaml
from typing import Dict, Any
from dataclasses import dataclass


@dataclass
class EmailClassification:
    """Email classification prompts configuration"""
    system_message: str
    user_prompt_template: str


@dataclass
class Prompts:
    """Prompts configuration class"""
    email_classification: EmailClassification


def load_prompts(filename: str) -> Prompts:
    """Load prompts from YAML file"""
    try:
        with open(filename, 'r') as f:
            data = yaml.safe_load(f)
    except FileNotFoundError:
        raise FileNotFoundError(f"Prompts file not found: {filename}")
    except yaml.YAMLError as e:
        raise ValueError(f"Failed to parse prompts YAML: {e}")

    # Extract email classification prompts
    email_class_data = data.get('email_classification', {})
    email_classification = EmailClassification(
        system_message=email_class_data.get('system_message', ''),
        user_prompt_template=email_class_data.get('user_prompt_template', '')
    )

    return Prompts(email_classification=email_classification)
