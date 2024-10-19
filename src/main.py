from aiogram import Bot, Dispatcher
import asyncio
from config import settings
import logging


logging.basicConfig(level=logging.INFO)
async def main():

    bot = Bot(settings.TG_API_KEY)
    dp = Dispatcher()
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())