"""
Email classification using OpenAI API
"""
import json
import re
from typing import Tuple, Optional
from dataclasses import dataclass
import openai
from openai import OpenAI


@dataclass
class OpenAIConfig:
    """OpenAI API configuration"""
    api_key: str
    endpoint: str
    model: str
    max_tokens: int
    temperature: float


@dataclass
class EmailClassificationConfig:
    """Email classification configuration"""
    system_message: str
    user_prompt_template: str


@dataclass
class ClassifierConfig:
    """Classifier configuration"""
    openai: OpenAIConfig
    email_classification: EmailClassificationConfig


class EmailClassifier:
    """Email classifier using OpenAI API"""
    
    def __init__(self, config: ClassifierConfig):
        self.config = config
        
        # Initialize OpenAI client
        client_config = {
            "api_key": config.openai.api_key,
        }
        
        # Set custom endpoint if provided
        if config.openai.endpoint and config.openai.endpoint != "https://api.openai.com/v1":
            client_config["base_url"] = self._normalize_endpoint(config.openai.endpoint)
        
        self.client = OpenAI(**client_config)
    
    def _normalize_endpoint(self, endpoint: str) -> str:
        """
        Normalize OpenAI endpoint URL.
        Ensures the endpoint ends with /v1 if it doesn't already.
        """
        endpoint = endpoint.rstrip('/')
        
        if endpoint.endswith('/v1'):
            return endpoint
        elif '/v1/' in endpoint:
            # Extract base URL up to /v1
            return endpoint[:endpoint.find('/v1') + 3]
        else:
            return endpoint + '/v1'
    
    def classify_email(self, email_text: str) -> Tuple[bool, str]:
        """
        Classify an email as important or unimportant.
        
        Args:
            email_text: The email content to classify
            
        Returns:
            Tuple of (is_important: bool, explanation: str)
        """
        # Build the prompt
        user_prompt = self.config.email_classification.user_prompt_template % email_text
        
        try:
            # Make API call
            response = self.client.chat.completions.create(
                model=self.config.openai.model,
                messages=[
                    {
                        "role": "system",
                        "content": self.config.email_classification.system_message
                    },
                    {
                        "role": "user", 
                        "content": user_prompt
                    }
                ],
                max_tokens=self.config.openai.max_tokens,
                temperature=self.config.openai.temperature
            )
            
            if not response.choices:
                print("No choices in OpenAI response")
                return False, "no choices"
            
            content = response.choices[0].message.content.strip()
            return self._parse_llm_decision(content)
            
        except Exception as e:
            print(f"OpenAI request failed: {e}")
            return False, f"request failed: {e}"
    
    def _parse_llm_decision(self, content: str) -> Tuple[bool, str]:
        """
        Parse the LLM response to extract importance decision and explanation.
        
        Args:
            content: Raw LLM response content
            
        Returns:
            Tuple of (is_important: bool, explanation: str)
        """
        # Try to extract JSON from the response
        json_match = re.search(r'\{[^}]*\}', content)
        if json_match:
            try:
                json_str = json_match.group()
                decision = json.loads(json_str)
                
                if isinstance(decision, dict):
                    important = decision.get('important', False)
                    explanation = decision.get('explanation', '')
                    return important, explanation
            except json.JSONDecodeError:
                pass
        
        # Fallback: simple heuristic based on keywords
        content_lower = content.lower()
        if any(keyword in content_lower for keyword in ['true', 'important', 'yes']):
            # Take first line as explanation
            explanation = content.split('\n')[0] if content else "classified as important"
            return True, explanation
        
        # Default to unimportant
        explanation = content.split('\n')[0] if content else "classified as unimportant"
        return False, explanation
