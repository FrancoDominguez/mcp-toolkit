from fastapi import FastAPI
from agent import MyAgent
from prompt_request import PromptRequest
import uvicorn
from logger import logger

app = FastAPI()


@app.post("/agent/prompt")
async def process_prompts(request: PromptRequest):
    try:
        agent = MyAgent()
        logger.info(f"Processing request: {request}")
        response = await agent.handle_request(request)
        return {"status": "success", "message": response}
    except Exception as e:
        logger.error(f"Error processing request: {e}")
        return {"status": "error", "message": str(e)}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
