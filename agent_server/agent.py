from agents.mcp.server import MCPServerStdio
from agents import Agent, Runner
from prompt_request import PromptRequest
import asyncio
from contextlib import AsyncExitStack
from dotenv import load_dotenv
import os
from logger import logger

load_dotenv()


class MyAgent:
    def __init__(self):
        self.logger = logger
        self.mcp_servers = [
            # MCPServerStdio(
            #     name="calculator",
            #     params={"command": "uvx", "args": ["mcp-calculator-demo"]},
            #     cache_tools_list=True,
            # ),
            # MCPServerStdio(
            #     name="filesystem",
            #     params={"command": "uvx", "args": ["mcp-filesystem"]},
            #     cache_tools_list=True,
            # ),
        ]

    async def handle_request(self, request: PromptRequest) -> str:
        self.logger.info(f"Processing request: {request}")

        async with AsyncExitStack() as stack:
            for mcp_server in self.mcp_servers:
                await stack.enter_async_context(mcp_server)

            agent = Agent(
                key=os.getenv("OPENAI_API_KEY"),
                instructions=request.system_prompt,
                mcp_servers=self.mcp_servers,
            )
            response = await Runner.run(agent, request.user_prompt)
            self.logger.info(f"Response: {response}")
            return response


if __name__ == "__main__":
    agent = MyAgent()
    request = PromptRequest(
        system_prompt="You are a helpful assistant.",
        user_prompt="What is 2+2?",
        context_prompt="",
    )
    asyncio.run(agent.handle_request(request))
