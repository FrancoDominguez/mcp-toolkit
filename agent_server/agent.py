import time
from dotenv import load_dotenv
from logger import logger
from contextlib import AsyncExitStack
from agents.mcp.server import MCPServerStdio
from agents import Agent, Runner, SQLiteSession
from prompt_request import PromptRequest
import os
load_dotenv()

class MyAgent:
    def __init__(self):
        self.logger = logger
        self.mcp_servers = [
            MCPServerStdio(
                name="gmail",
                params={
                    "command": "npx",
                    "args": ["@gongrzhe/server-gmail-autoauth-mcp"]
                }
            )
        ]

    async def handle_request(self, prompt_request: PromptRequest) -> str:
        self.logger.info(f"Processing request...")
        self.logger.info(f"System prompt: '{prompt_request.system_prompt}'")
        self.logger.info(f"User prompt: '{prompt_request.user_prompt}'")

        async with AsyncExitStack() as stack:
            start_time = time.time()
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
            
            end_time = time.time()
            self.logger.info(f"MCP server initialization completed in {end_time - start_time:.2f} seconds")
            
            session = SQLiteSession(prompt_request.conversation_id, os.getenv("SESSION_DB_PATH"))

            agent = Agent(
                name="Jarvis",
                instructions=prompt_request.system_prompt,
                mcp_servers=self.mcp_servers,
            )
            result = await Runner.run(agent, prompt_request.user_prompt, session=session)
            self.logger.info(f"Result: '{result.final_output}'")
            return result.final_output