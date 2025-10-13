from pydantic import BaseModel


class PromptRequest(BaseModel):
    system_prompt: str
    user_prompt: str
    conversation_id: str
