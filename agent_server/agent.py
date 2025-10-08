import sys
import asyncio
from dotenv import load_dotenv
from logger import logger
from contextlib import AsyncExitStack
from agents.mcp.server import MCPServerStdio
from agents import Agent, Runner, SQLiteSession
from prompt_request import PromptRequest

load_dotenv()

class MyAgent:
    def __init__(self):
        self.logger = logger
        self.mcp_servers = [
            MCPServerStdio(
                name="gmail-mcp",
                params={
                    "url": "http://localhost:8000/sse"
                }
            )
        ]

    async def handle_request(self, prompt_request: PromptRequest) -> str:
        self.logger.info(f"Processing request...")
        self.logger.info(f"System prompt: '{prompt_request.system_prompt}'")
        self.logger.info(f"User prompt: '{prompt_request.user_prompt}'")

        async with AsyncExitStack() as stack:
            to_remove = []
            for server in list(
                self.mcp_servers
            ):  # big opportunity for speed improvements here using asyncio
                try:
                    await stack.enter_async_context(server)
                except Exception as e:
                    self.logger.error(f"Removing mcp server from agent config, error: {e}")
                    to_remove.append(server)
            for s in to_remove:
                self.mcp_servers.remove(s)
            
            session = SQLiteSession("my_session")

            agent = Agent(
                name="Jarvis",
                instructions=prompt_request.system_prompt,
                mcp_servers=self.mcp_servers,
            )
            result = await Runner.run(agent, prompt_request.user_prompt, session=session)
            self.logger.info(f"Result: '{result.final_output}'")
            return result.final_output

