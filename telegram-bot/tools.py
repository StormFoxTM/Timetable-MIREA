import re
from datetime import date

import requests
from tabulate import tabulate


def get_timetable(key, value, period, view='table'):
    """
    Requests timetable from API and parses the response into messages.
    
    :param key: Query key (`'group'`, `'lecturer'` or `'auditorium'`)
    :param value: Query value (group/lecturer/auditorium name)
    :param period: String with time period for the timetable requested
    :param view: `'table'` or `'list'`, determines output format, defaults to `'table'`
    :return: List of formatted messages each containing a title and timetable for a day
    """
    # Check for valid group / lecturer / auditorium via API request
    query = {key: value}
    # The if statement below is a temporary fix!
    # API currently returns status code 400 for auditorium info requests
    if key != 'auditorium':
        response = requests.get('http://mirea-club.site/api/info', params=query)
        if response.status_code != 200:
            tkey = {'group': 'группу', 'lecturer': 'преподавателя', 'auditorium': 'аудиторию'}.get(key)
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
    response = requests.get('http://mirea-club.site/api/timetable', params=query)

    # Parse response into tables, form messages
    try:
        messages = parse_response(response, key, view=view)
    except:
        return "Получен некорректный ответ от сервера."
    daynames = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота', 'Воскресенье']
    if day is None:
        day = 1
    tmsgs = []
    for msg in messages:
        if msg:
            title = daynames[day - 1]
            if view == 'table':
                tmsgs.append(f"**{title}:**\n```\n{msg}\n```")
            elif view == 'list':
                tmsgs.append(f"**{title}:**\n{msg}")
        day += 1
    return tmsgs


def parse_response(response, mode, view='table'):
    """
    Parses JSON response from API.
    
    :param response: Response object to parse
    :param mode: Query key
    :param view: `'table'` or `'list'`, determines output format, defaults to `'table'`
    :return: List of formatted messages each containing a timetable for a day
    """
    tables = []

    # Filter table headers and json keys needed for parse mode
    headers = ['Пара', 'Предмет', 'Группа', 'Преподаватель', 'Аудитория']
    keys = ['subject_to_number', 'subject_title', 'name_group', 'name_lecturer', 'auditorium']
    index = {'group': 2, 'lecturer': 3, 'auditorium': 4}.get(mode)
    del headers[index]
    del keys[index]

    # Parse json
    for day in response.json()['weeks']:
        table = []
        if day['day'] is None:
            tables.append(None)
            continue
        for subject in day['day']:
            table.append([subject[key] for key in keys])

        # Format parsed data into tables or list
        if view == 'table':
            tables.append(tabulate(table, headers=headers, tablefmt="github"))
        elif view == 'list':
            tables.append('\n'.join([' - '.join([str(x) for x in row]) for row in table]))
    return tables


def parse_msg(text):
    """
    Reads a text message, determines if it is a timetable query.
    If it is, tries to parse time period and query type.
    
    :param text: Message text
    :return: Query key (`'group'`, `'lecturer'` or `'auditorium'`), query value, time period
    """
    text = str.lower(text)
    words = text.split()

    # Detecting if text is a timetable request
    req = None
    req_src = ['расписание', 'пары']
    for substr in req_src:
        if substr in words:
            req = substr
            words.remove(req)
            break

    if not req:
        return None

    # Detecting time period
    period = None
    period_src = ['сегодня', 'завтра', 'неделя', 'неделю', 'нед', 'нед.']
    period_next_src = ['сл', 'сл.', 'след', 'след.', 'следующая', 'следующую']
    for substr in period_src:
        if substr in words:
            period = substr
            words.remove(period)
            break
    if period in period_src[2:]:
        period = 'неделя'
        for substr in period_next_src:
            if substr in words:
                words.remove(substr)
                period = 'след. неделя'
                break
    if not period:
        period = 'сегодня'

    text = ' '.join(words)

    # Detecting group, auditorium or lecturer
    temp = parse_wrapper(text)
    if temp is None:
        return None
    key, value = temp
    return key, value, period


def parse_wrapper(text, mode=None):
    """
    Wrapper for parser functions. Parses query value
    
    :param text: Input string
    :param mode: Query key if parsing for a specific query type, defaults to `None`
    :return: Query key and value if mode is `None`, else only query value
    """
    parsers = [(parse_group, 'group'), (parse_auditorium, 'auditorium'), (parse_lecturer, 'lecturer')]
    str_value = None
    for parser, key in parsers:
        if key == mode or mode is None:
            str_value = parser(text)
            if str_value:
                return str_value if mode else (key, str_value)


def parse_group(text):
    """
    Parses an input string, trying to find a group name.
    
    :param text: Input string
    :return: Formatted string (group name) if found, `None` otherwise
    """
    group = None
    matches = re.findall(r'([а-я]{4})[ -]?(\d{2})[ -]?(\d{2})', text.lower())
    if len(matches):
        group = '-'.join(matches[0])
    return group


def parse_auditorium(text):
    """
    Parses an input string, trying to find an auditorium name.
    
    :param text: Input string
    :return: Formatted string (auditorium name) if found, `None` otherwise
    """
    auditorium = None
    matches = re.findall(r'([а-я]+)[ -]?(\d+)(?:[ -]?([\dа-я]))?', text.lower())
    if len(matches):
        auditorium = '-'.join([x for x in matches[0] if x])
    return auditorium


def parse_lecturer(text):
    """
    Parses an input string, trying to find a lecturer name.
    
    :param text: Input string
    :return: Formatted string (lecturer name) if found, `None` otherwise
    """
    lecturer = None
    matches = re.findall(r'([а-я]{3,}) ?([а-я])?[ .]?([а-я])?[ .]?', text.lower())
    if len(matches):
        ln, fn, mn = matches[0]
        lecturer = f"{ln}{' ' + fn + '.' if fn else ''}{mn + '.' if mn else ''}"
    return lecturer


def get_today():
    """Returns week parity (1-2) and day number (1-7) for current day."""
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
    """Takes week parity and day number for a day, returns these values for next day"""
    day += 1
    if day > 7:
        day -= 7
        week = week % 2 + 1
    return week, day


def next_week(week, day):
    """Takes week parity and day number for a day, returns these values for next week"""
    week = week % 2 + 1
    return week, day
