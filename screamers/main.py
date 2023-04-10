""" Screamers Operators Interface Telegram bot ver 0.1
    by Vladimir Monomakh
"""

# Import
from pymongo import MongoClient

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

from config import TG_TOKEN
from config import MONGODB_LINK
from config import MONGO_DB

from keyboards import *

import certifi

ca = certifi.where()
# tele_Bot = telebot.TeleBot(TG_TOKEN, parse_mode=None)

# -------------------------------
# Bot Logic
# -------------------------------

# bot = Bot(TG_TOKEN)

# Connect to DataBase
mondb = MongoClient(MONGODB_LINK, tlsCAFile=ca)[MONGO_DB]

cursor = mondb.runners.watch()
document = next(cursor)

print(cursor)


def on_start(update: Update, context: CallbackContext):
    message = update.message

    message.reply_text(
        'Операторская система ввода данных. Вы не авторизованны пожалуйста введите пинкод.',
        reply_markup=REPLY_KEYBOARD_MARKUP
    )

def send_runner_nunmber():
    pass


screamers_id = {}


def handle_text(update: Update, context: CallbackContext):
    message = update.message
    query = update.callback_query
    text = update.message.text

    if text.isdigit and len(text) == 1:
        if message.chat.id in runners_code.keys():
            runners_code[message.chat.id] += text
        else:
            runners_code[message.chat.id] = ""
            runners_code[message.chat.id] += text
    elif text == "DEL":
        if message.chat.id in runners_code.keys():
            if len(runners_code[message.chat.id]) > 0:
                runners_code[message.chat.id] = runners_code[message.chat.id][:-1]
    elif text == "OK":
        if message.chat.id in runners_code.keys():
            if len(runners_code[message.chat.id]) > 0:
                message.reply_text(
                    'Номер: ' + runners_code[message.chat.id]
                )
    print(runners_code)



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

def main():
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


if __name__ == '__main__':
    main()
