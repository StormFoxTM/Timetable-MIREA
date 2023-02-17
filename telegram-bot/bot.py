import telebot

# This is how it should be, but setting env variables in Windows is pain
# import os
# BOT_TOKEN = os.environ.get('BOT_TOKEN')

# This is how it absolutely shouldn't be, but let's at least test it
BOT_TOKEN = "<token>"

bot = telebot.TeleBot(BOT_TOKEN)


def get_timetable(group, period):
    # This is a placeholder function for now
    # It will call the API for timetable
    return f"<Расписание группы *{group.upper()}* на *{period.lower()}*>"


@bot.message_handler(commands=['start', 'help'])
def send_welcome(message):
    text = """\
Привет! Я бот, присылающий расписание МИРЭА.
Список команд:
- `/start`, `/help`: выводят это сообщение
- `/timetable`: присылает расписание
"""
    bot.send_message(message.chat.id, text, parse_mode='Markdown')


@bot.message_handler(commands=['timetable'])
def group_handler(message):
    text = "Напишите группу, расписание которой хотите получить."
    sent_msg = bot.send_message(message.chat.id, text)
    bot.register_next_step_handler(sent_msg, period_handler)


def period_handler(message):
    group = message.text
    text = """\
Напишите, расписание на какой период вы хотите получить.
Варианты: *сегодня*, *завтра*, *неделя*, *месяц*.
"""
    sent_msg = bot.send_message(message.chat.id, text, parse_mode='Markdown')
    bot.register_next_step_handler(sent_msg, fetch_timetable, group)


def fetch_timetable(message, group):
    period = message.text
    timetable = get_timetable(group, period)
    text = ("Я *типа* спарсил группу и период, "
            "*типа* отправил запрос к API, "
            "*типа* принял ответ с данными "
            "и сейчас *типа* отправлю вам.")
    bot.send_message(message.chat.id, text, parse_mode='Markdown')
    bot.send_message(message.chat.id, timetable, parse_mode='Markdown')


@bot.message_handler(func=lambda msg: True)
def echo_all(message):
    bot.reply_to(message, message.text)


bot.infinity_polling()
