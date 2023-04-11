""" Screamers Operators Interface Telegram bot ver 0.1
    by Vladimir Monomakh
"""

# Import
from pymongo import MongoClient

import requests

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

from config import TG_TOKEN, MONGO_DB, MONGODB_LINK, PIN_CODE


from keyboards import *

import certifi

ca = certifi.where()
# tele_Bot = telebot.TeleBot(TG_TOKEN, parse_mode=None)

# -------------------------------
# Bot Logic
# -------------------------------

# bot = Bot(TG_TOKEN)

runners_code = {}

pins = {}

def getNumberRequest(numb: int):
    x = requests.get('http://127.0.0.1:8090/numb/?numb=' + numb)
    print(x.status_code)



def on_start(update: Update, context: CallbackContext):
    message = update.message
    if message.chat.id in runners_code.keys():
        message.reply_text(
            'Вы авторизованы! ',
            reply_markup=REPLY_KEYBOARD_MARKUP
        )
    else: 
         message.reply_text(
            'Добро пожаловать! Пожалуйста введите Пин-код. ',
            reply_markup=REPLY_KEYBOARD_MARKUP
        )


def handle_text(update: Update, context: CallbackContext):
    message = update.message
    query = update.callback_query
    text = update.message.text

    if message.chat.id in runners_code.keys():
        if text.isdigit and len(text) == 1:
            runners_code[message.chat.id] += text
        elif text == "DEL":
            if len(runners_code[message.chat.id]) > 0:
                runners_code[message.chat.id] = runners_code[message.chat.id][:-1]
        elif text == "OK":
            if len(runners_code[message.chat.id]) > 0:
                getNumberRequest(runners_code[message.chat.id])
                message.reply_text(
                    'Номер: ' + runners_code[message.chat.id]
                )
                runners_code[message.chat.id] = ""

    else:
        if message.chat.id in pins.keys():
            if text.isdigit and len(text) == 1:
                pins[message.chat.id] += text
            elif text == "DEL":
                if len(pins[message.chat.id]) > 0:
                    pins[message.chat.id] = pins[message.chat.id][:-1]
            elif text == "OK":
                if len(pins[message.chat.id]) > 0:
                    if pins[message.chat.id] == PIN_CODE:
                        runners_code[message.chat.id] = ""
                        message.reply_text(
                            'Вход выполнен! ' 
                        )
                    else:
                        pins[message.chat.id] = ""
                        message.reply_text(
                            'Неверный Пин-Код! Попробуйте ещё раз!' 
                        )
        else:
            if text.isdigit and len(text) == 1:
                pins[message.chat.id] = ""
                pins[message.chat.id] += text
    print(runners_code)



def main():
    updater = Updater(
        token=TG_TOKEN,
        use_context=True,
    )
    dp = updater.dispatcher
    dp.add_handler(CommandHandler('start', on_start))
    #dp.add_handler(CallbackQueryHandler(callback=keyboard_call_handler, pass_chat_data=True))
    dp.add_handler(MessageHandler(Filters.text, handle_text))
    updater.start_polling()
    updater.idle()


if __name__ == '__main__':
    main()
