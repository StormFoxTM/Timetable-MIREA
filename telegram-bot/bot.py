import os

import telebot as tb
from dotenv import load_dotenv
from tools import *

load_dotenv()
BOT_TOKEN = os.getenv('BOT_TOKEN')
bot = tb.TeleBot(BOT_TOKEN)


@bot.message_handler(commands=['start', 'help'])
def send_welcome(message):
    text = """\
Привет! Я бот, присылающий расписание МИРЭА.

Список команд:
- `/start`, `/help`: выводят это сообщение
- `/timetable`, `/group`: присылает расписание группы
- `/lecturer`: присылает расписание преподавателя
- `/auditorium`: присылает расписание по аудитории

Также можно запросить расписание одним сообщением, например:
- "Расписание ИКБО-02-20" (период по умолчанию - сегодня)
- "Пары А-10 завтра"
- "Расписание Волков след. неделя"
"""
    bot.send_message(message.chat.id, text, parse_mode='Markdown')


@bot.message_handler(commands=['timetable', 'group'])
def group_handler(message):
    text = "Напишите группу, расписание которой хотите получить."
    sent_msg = bot.send_message(message.chat.id, text)
    key = 'group'
    bot.register_next_step_handler(sent_msg, period_handler, key)


@bot.message_handler(commands=['lecturer'])
def lecturer_handler(message):
    text = "Напишите преподавателя, чье расписание хотите получить."
    sent_msg = bot.send_message(message.chat.id, text)
    key = 'lecturer'
    bot.register_next_step_handler(sent_msg, period_handler, key)


@bot.message_handler(commands=['auditorium'])
def auditorium_handler(message):
    text = "Напишите аудиторию, расписание по которой хотите получить."
    sent_msg = bot.send_message(message.chat.id, text)
    key = 'auditorium'
    bot.register_next_step_handler(sent_msg, period_handler, key)


@bot.message_handler(func=lambda msg: True)
def text_handler(message):
    params = parse_msg(message.text)
    if params is None:
        text = "Не удалось определить запрос."
        bot.reply_to(message, text)
        return
    key, value, period = params
    timetable = get_timetable(key, value, period, view='list')
    tkey = {'group': 'группы', 'lecturer': 'преподавателя', 'auditorium': 'аудитории'}.get(key)
    tperiod = ((period[:-1] + 'ю') if period[-2:] == 'ля' else period).lower()
    text = f"Ниже приведено расписание {tkey} {value.title()} на {tperiod}."
    bot.send_message(message.chat.id, text, parse_mode='Markdown')
    if type(timetable) == str:
        timetable = [timetable, ]
    for tmsg in timetable:
        bot.send_message(message.chat.id, tmsg, parse_mode='Markdown')


def period_handler(message, key):
    value = message.text
    text = "Выберите, расписание на какой период вы хотите получить."
    markup = tb.types.ReplyKeyboardMarkup(row_width=2, one_time_keyboard=True)
    periods = ['Сегодня', 'Завтра', 'Неделя', 'След. неделя']
    btns = [tb.types.KeyboardButton(period) for period in periods]
    markup.add(*btns)
    sent_msg = bot.send_message(message.chat.id, text, parse_mode='Markdown', reply_markup=markup)
    bot.register_next_step_handler(sent_msg, fetch_timetable, key, value)


def fetch_timetable(message, key, value):
    period = message.text
    value = parse_wrapper(value, mode=key)
    markup = tb.types.ReplyKeyboardRemove()
    if value is None:
        tkey = {'group': 'группу', 'lecturer': 'преподавателя', 'auditorium': 'аудиторию'}.get(key)
        text = f"Не удалось распознать {tkey}."
        bot.send_message(message.chat.id, text, parse_mode='Markdown', reply_markup=markup)
        return
    timetable = get_timetable(key, value, period, view='list')
    tkey = {'group': 'группы', 'lecturer': 'преподавателя', 'auditorium': 'аудитории'}.get(key)
    tperiod = ((period[:-1] + 'ю') if period[-2:] == 'ля' else period).lower()
    text = f"Ниже приведено расписание {tkey} {value.title()} на {tperiod}."
    bot.send_message(message.chat.id, text, parse_mode='Markdown', reply_markup=markup)
    if type(timetable) == str:
        timetable = [timetable, ]
    for tmsg in timetable:
        bot.send_message(message.chat.id, tmsg, parse_mode='Markdown')


if __name__ == '__main__':
    print("Бот успешно запущен!")
    bot.infinity_polling()
