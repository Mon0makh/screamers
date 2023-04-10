from telegram import ReplyKeyboardMarkup

custom_keyboard_login = [['1', '2', '3'],['4', '5', '6'], ['7', '8', '9'], ['OK', '0', 'DEL']]
REPLY_KEYBOARD_MARKUP = ReplyKeyboardMarkup(custom_keyboard_login, resize_keyboard=True)
