""" Screamers Operators Interface Telegram bot ver 0.1
    by Vladimir Monomakh
"""

# Import
from pymongo import MongoClient



import uvicorn
import logging


from telegram import Bot
from telegram import Update

from telegram import ReplyKeyboardRemove

from telegram.ext import CallbackQueryHandler
from telegram.ext import Updater
from telegram.ext import CommandHandler
from telegram.ext import MessageHandler
from telegram.ext import Filters
from telegram.ext import CallbackContext

from keyboards import *
from config import TG_TOKEN, MONGODB_LINK, MONGO_DB, SERVER_HOST, SERVER_PORT

import certifi



ca = certifi.where()
# tele_Bot = telebot.TeleBot(TG_TOKEN, parse_mode=None)

# -------------------------------
# Bot Logic
# -------------------------------

bot = Bot(TG_TOKEN)

# Connect to DataBase
mondb = MongoClient(MONGODB_LINK)[MONGO_DB]

# cursor = mondb.runners.watch()
# document = next(cursor)

# print(cursor)


screamers_id = []

async def send_runner_number(name: str):
    print(name)
    bot.send_message(chat_id=screamers_id[0], text=name)
    

resume_token = None
pipeline = [{'$match': {'operationType': 'update'}}]
with mondb.runners.watch() as stream:
    while stream.alive:
        change = stream.try_next()
        # Note that the ChangeStream's resume token may be updated
        # even when no changes are returned.
        if change is not None:
            print(change)
            runner = mondb.runners.find_one({'_id' : change['documentKey']['_id']})
            send_runner_number(runner['Name'])
            continue






def on_start(update: Update, context: CallbackContext):
    message = update.message
    screamers_id.append[message.chat.id]
    message.reply_text(
        'Вы авторизованы! ',
        reply_markup=REPLY_KEYBOARD_MARKUP
    )






def handle_text(update: Update, context: CallbackContext):
    message = update.message
    query = update.callback_query
    text = update.message.text

    # if text.isdigit and len(text) == 1:
    #     if message.chat.id in runners_code.keys():
    #         runners_code[message.chat.id] += text
    #     else:
    #         runners_code[message.chat.id] = ""
    #         runners_code[message.chat.id] += text
    # elif text == "DEL":
    #     if message.chat.id in runners_code.keys():
    #         if len(runners_code[message.chat.id]) > 0:
    #             runners_code[message.chat.id] = runners_code[message.chat.id][:-1]
    # elif text == "OK":
    #     if message.chat.id in runners_code.keys():
    #         if len(runners_code[message.chat.id]) > 0:
    #             message.reply_text(
    #                 'Номер: ' + runners_code[message.chat.id]
    #             )
    # print(runners_code)





# Telegram inline menu buttons handler
def keyboard_call_handler(update: Update, context: CallbackContext):
    query = update.callback_query
    data = query.data

    # if data == CALLBACK_MM:
    #     query.edit_message_text(
    #         text="Основное меню: ",
    #         reply_markup=get_main_menu()
    #     )
    # elif data == CALLBACK_MM_HUB:
    #     pass


async def main():
    updater = Updater(
        token=TG_TOKEN,
        use_context=True,
    )
    dp = updater.dispatcher
    dp.add_handler(CommandHandler('start', on_start))
    dp.add_handler(CallbackQueryHandler(callback=keyboard_call_handler, pass_chat_data=True))
    dp.add_handler(MessageHandler(Filters.text, handle_text))
    updater.start_polling()
    updater.idle()


if __name__ == "__main__":
    main()