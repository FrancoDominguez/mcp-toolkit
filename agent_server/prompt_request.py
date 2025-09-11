from pydantic import BaseModel


class PromptRequest(BaseModel):
    system_prompt: str
    user_prompt: str
    context_db_url: str
