import os
from datetime import date

import requests
import telebot as tb
from dotenv import load_dotenv

load_dotenv()

BOT_TOKEN = os.getenv('BOT_TOKEN')

bot = tb.TeleBot(BOT_TOKEN)

print("Бот успешно запущен!")


def get_timetable(key, value, period):
    # Check for valid time period
    try:
        match period.lower():
            case 'сегодня':
                week, day = get_today()
            case 'завтра':
                week, day = next_day(*get_today())
            case 'неделя', 'след. неделя':
                raise NotImplementedError
            case _:
                raise ValueError
    except NotImplementedError:
        return "Запрос расписания на несколько дней еще не реализован."
    except ValueError:
        return "Не удалось определить период, на который запрашивается расписание."
    
    # Check for valid group / lecturer / auditorium via API request
    # query = {key: value}
    # response = requests.get(url, params=query)
    # ...

    # Get timetable via API request
    query = {key: value, 'week': str(week), 'day': str(day)}
    # response = requests.get(url, params=query)
    # ...

    # This response is a placeholder
    return f"К API отправляется get-запрос с параметрами:\n`{str(query)}`"


@bot.message_handler(commands=['start', 'help'])
def send_welcome(message):
    text = """\
Привет! Я бот, присылающий расписание МИРЭА.
Список команд:
- `/start`, `/help`: выводят это сообщение
- `/timetable`, `/group`: присылает расписание группы
- `/lecturer`: присылает расписание преподавателя
- `/auditorium`: присылает расписание по аудитории
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


def period_handler(message, key):
    value = message.text
    text = "Выберите, расписание на какой период вы хотите получить."
    markup = tb.types.ReplyKeyboardMarkup(row_width=2, one_time_keyboard=True)
    periods = ['Сегодня', 'Завтра', 'Неделя', 'След. неделя']
    btns = [tb.types.KeyboardButton(period) for period in periods]
    markup.add(*btns)
    sent_msg = bot.send_message(
        message.chat.id, text, parse_mode='Markdown', reply_markup=markup)
    bot.register_next_step_handler(sent_msg, fetch_timetable, key, value)


def fetch_timetable(message, key, value):
    period = message.text
    timetable = get_timetable(key, value, period)
    text = ("На данном этапе API еще не развернут, "
            "поэтому ниже будут выведены действия бота.")
    markup = tb.types.ReplyKeyboardRemove()
    bot.send_message(message.chat.id, text,
                     parse_mode='Markdown', reply_markup=markup)
    bot.send_message(message.chat.id, timetable, parse_mode='Markdown')


def get_today():
    today = date.today()

    # Set first day of semester to Sep 1 or Feb 9
    firstday = (9, 1) if today.month >= 9 else (2, 9)
    firstday = date(today.year, *firstday)

    # Set firstday to Monday of first week
    if firstday.weekday() == 6:
        firstday = firstday.replace(day=firstday.day+1)
    else:
        firstday = firstday.replace(day=firstday.day-firstday.weekday())

    # Get week parity (1 if odd, 2 if even)
    weekparity = (((today - firstday).days // 7) % 2) + 1

    return weekparity, today.isoweekday()


def next_day(week, day):
    day += 1
    if day == 7:
        day = 1
        week = (week + 1) % 2
    return week, day


bot.infinity_polling()
