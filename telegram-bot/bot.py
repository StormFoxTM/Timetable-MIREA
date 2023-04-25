import os
from datetime import date
from pprint import pprint

import requests
import telebot as tb
from dotenv import load_dotenv
from tabulate import tabulate

# Loading bot token
load_dotenv()
BOT_TOKEN = os.getenv('BOT_TOKEN')
bot = tb.TeleBot(BOT_TOKEN)
print("Бот успешно запущен!")


def get_timetable(key, value, period):
    # Check for valid group / lecturer / auditorium via API request
    query = {key: value}
    response = requests.get('http://mirea-club.site/api/info', params=query)
    if response.status_code != 200:
        tkey = {'group': 'группу', 'lecturer': 'преподавателя',
                'auditorium': 'аудиторию'}.get(key)
        return f"Не удалось найти {tkey} в расписании."

    # Check for valid time period
    week, day = None, None
    try:
        match period.lower():
            case 'сегодня':
                week, day = get_today()
                query = {key: value, 'week': week, 'day': day}
            case 'завтра':
                week, day = next_day(*get_today())
                query = {key: value, 'week': week, 'day': day}
            case 'неделя':
                week, _ = get_today()
                query = {key: value, 'week': week}
            case 'след. неделя':
                week, _ = next_week(*get_today())
                query = {key: value, 'week': week}
            case _:
                raise ValueError
    except ValueError:
        return "Не удалось определить период, на который запрашивается расписание."
    except not ValueError:
        return "Что-то пошло не так."

    # Get timetable via API request
    response = requests.get(
        'http://mirea-club.site/api/timetable', params=query)
    
    # Parse response into tables, form messages
    messages = parse_response(response, key)
    daynames = ['Понедельник', 'Вторник', 'Среда',
                'Четверг', 'Пятница', 'Суббота', 'Воскресенье']
    if day is None:
        day = 1
    tmsgs = []
    for msg in messages:
        if msg:
            title = daynames[day - 1]
            tmsgs.append(f"**{title}:** \n```\n{msg}\n```")
        day += 1
    return tmsgs


def parse_response(response, mode):
    tables = []
    
    # Filter table headers and json keys needed for parse mode
    headers = ['Пара', 'Предмет', 'Группа', 'Преподаватель', 'Аудитория']
    keys = ['subject_to_number', 'subject_title',
            'name_group', 'name_lecturer', 'auditorium']
    index = {'group': 2, 'lecturer': 3, 'auditorium': 4}.get(mode)
    del headers[index]
    del keys[index]

    # Parse json
    for day in response.json()['weeks']:
        table = []
        if day['day'] is None:
            tables.append(None)
        for subject in day['day']:
            table.append([subject[key] for key in keys])
            
        # Format parsed data into tables
        tables.append(tabulate(table, headers=headers, tablefmt="github"))
    return tables


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
    tkey = {'group': 'группы', 'lecturer': 'преподавателя',
            'auditorium': 'аудитории'}.get(key)
    text = f"Ниже приведено расписание {tkey} {value} на запрошенный период."
    markup = tb.types.ReplyKeyboardRemove()
    bot.send_message(message.chat.id, text,
                     parse_mode='Markdown', reply_markup=markup)
    if type(timetable) == str:
        timetable = [timetable, ]
    for tmsg in timetable:
        bot.send_message(message.chat.id, tmsg, parse_mode='Markdown')


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
    if day > 7:
        day -= 7
        week = week % 2 + 1
    return week, day


def next_week(week, day):
    week = week % 2 + 1
    return week, day


bot.infinity_polling()
