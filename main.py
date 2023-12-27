import asyncio
import os
from datetime import datetime

import aiohttp
from dotenv import load_dotenv
from fastapi import FastAPI
from fastapi_utils.tasks import repeat_every
from motor import motor_asyncio as motor

load_dotenv()

app = FastAPI()
database_client = motor.AsyncIOMotorClient(os.environ.get("MONGO"))
database = database_client["tensyo-web"]
address = "play.tensyoserver.net"


async def get_status(session: aiohttp.ClientSession, address: str, port: str, name: str) -> dict:
    async with session.get(f"https://api.mcstatus.io/v2/status/java/{address}:{port}") as res:
        data = await res.json()

    num = 0
    online = 0

    past_data = []
    async for item in database.status.find({"name": str(name.lower())}, {"_id": False}):
        num += 1

        if num > 43:
            break

        if item["online"]:
            online += 1

        past_data.append(item)

    if data["online"]:
        online += 1

        return {
            "name": name,
            "percent": int(online / 44 * 100),
            "online": True,
            "port": port,
            "players": {
                "online": data["players"]["online"],
                "max": data["players"]["max"],
                "list": data["players"]["list"],
            },
            "past_data": past_data,
        }
    else:
        return {
            "name": name,
            "percent": int(online / 44 * 100),
            "online": False,
            "port": port,
            "past_data": past_data,
        }


@app.get("/v1/status")
async def status() -> None:
    async with aiohttp.ClientSession() as session:
        tasks = [
            asyncio.ensure_future(get_status(session, address, "25565", "Proxy")),
            asyncio.ensure_future(get_status(session, address, "25566", "Lobby")),
            asyncio.ensure_future(get_status(session, address, "25570", "Survival")),
            asyncio.ensure_future(get_status(session, address, "25567", "MiniGame")),
            asyncio.ensure_future(get_status(session, address, "25568", "PVE")),
            asyncio.ensure_future(get_status(session, address, "25559", "Build")),
            asyncio.ensure_future(get_status(session, address, "25569", "Railway")),
        ]

        result = await asyncio.gather(*tasks)

    num = 0
    info = ""
    online = 0

    for item in result:
        num += 1

        if item["online"]:
            online += 1

    if num == online:
        info = "online_all"
    elif num != online:
        info = "partial_down"
    elif online == 0:
        info = "all_down"

    return {"status": 200, "info": info, "data": result}


@app.on_event("startup")
@repeat_every(seconds=60)
async def startup_event():
    async with aiohttp.ClientSession() as session:
        tasks = [
            asyncio.ensure_future(get_status(session, address, "25565", "Proxy")),
            asyncio.ensure_future(get_status(session, address, "25566", "Lobby")),
            asyncio.ensure_future(get_status(session, address, "25570", "Survival")),
            asyncio.ensure_future(get_status(session, address, "25567", "MiniGame")),
            asyncio.ensure_future(get_status(session, address, "25568", "PVE")),
            asyncio.ensure_future(get_status(session, address, "25559", "Build")),
            asyncio.ensure_future(get_status(session, address, "25569", "Railway")),
        ]

        result = await asyncio.gather(*tasks)

    for item in result:
        await database.status.insert_one(
            {"name": str(item.get("name")).lower(), "online": item["online"], "time": datetime.now()}
        )

    await asyncio.sleep(72)
