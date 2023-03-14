package main

// Импортируем необходимые модули и библиотеки
import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"knocker/pgsql" // Импортируем модуль pgsql из нашего собственного пакета knocker

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Основная функция main
func main() {
	// Создаем новый роутер Gin с настройками по умолчанию
	router := gin.Default()

	// Добавляем обработчик GET-запроса по пути "/api/timetable"
	router.GET("/api/timetable", getFunction)

	// Запускаем сервер Gin на порту 9888
	err := router.Run(":9888")
	if err != nil {
		// В случае ошибки выводим сообщение об ошибке и завершаем программу
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

func getFunction(context *gin.Context) {
	// Получаем параметры запроса
	data := context.Request.URL.Query()

	// Получаем количество параметров "group", "lecture" и "auditorium"
	group := len(data["group"])
	lecture := len(data["lecture"])
	auditorium := len(data["auditorium"])

	// В зависимости от того, сколько параметров указано, выбираем соответствующую функцию
	switch 1 {
	case group:
		// Получаем данные по группе и кодируем их в JSON
		dataGroup, err := get_group(context)
		dataGroupJSON, err_json := json.Marshal(dataGroup)

		// Если не возникло ошибок, возвращаем данные по группе в формате JSON
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataGroupJSON))
		} else if err == pgx.ErrNoRows {
			// Если нет данных, возвращаем ошибку "По запросу ничего не найдено"
			context.String(http.StatusBadRequest, "По запросу ничего не найдено")
		} else {
			// Если возникли ошибки, возвращаем ошибку "Ошибка в запросе!"
			context.String(http.StatusBadRequest, "Ошибка в запросе!")
		}

	case lecture:
		// Получаем данные по лектору и кодируем их в JSON
		dataLectur, err := get_lecturer(context)
		dataJSON, err_json := json.Marshal(dataLectur)

		// Если не возникло ошибок, возвращаем данные по лектору в формате JSON
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataJSON))
		} else if err == pgx.ErrNoRows {
			// Если нет данных, возвращаем ошибку "По запросу ничего не найдено"
			context.String(http.StatusBadRequest, "По запросу ничего не найдено")
		} else {
			// Если возникли ошибки, возвращаем ошибку "Ошибка в запросе!"
			context.String(http.StatusBadRequest, "Ошибка в запросе!")
		}

	case auditorium:
		// Получаем данные по аудитории и кодируем их в JSON
		dataAuditorium, err := get_auditorium(context)
		dataJSON, err_json := json.Marshal(dataAuditorium)

		// Если не возникло ошибок, возвращаем данные по аудитории в формате JSON
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataJSON))
		} else if err == pgx.ErrNoRows {
			// Если нет данных, возвращаем ошибку "По запросу ничего не найдено"
			context.String(http.StatusBadRequest, "По запросу ничего не найдено")
		} else {
			// Если возникли ошибки, возвращаем ошибку "Ошибка в запросе!"
			context.String(http.StatusBadRequest, "Ошибка в запросе!")
		}

	// Если количество параметров не соответствует ни одному из вариантов, возвращаем ошибку "Ошибка в запросе!"
	default:
		context.String(http.StatusBadRequest, "Ошибка в запросе!")
	}
}

// get_group function: получает группу и запрашивает ее расписание из базы данных по заданным дням и неделе.
func get_group(context *gin.Context) (pgsql.DataGroupRequests, error) {
	// Получаем параметры запроса из URL.
	data := context.Request.URL.Query()
	// Получаем номер текущей недели и дня недели.
	week_int, day_int := getWeekAndDay(context)
	// Получаем расписание из базы данных на основе параметров запроса и текущей даты.
	return pgsql.GetTimetableGroup(strings.ToUpper(data["group"][0]), week_int, day_int)
}

// get_lecturer function: получает преподавателя и запрашивает его расписание из базы данных по заданным дням и неделе.
func get_lecturer(context *gin.Context) (pgsql.DataLecturRequests, error) {
	// Получаем параметры запроса из URL.
	data := context.Request.URL.Query()
	// Получаем номер текущей недели и дня недели.
	week_int, day_int := getWeekAndDay(context)
	// Получаем расписание из базы данных на основе параметров запроса и текущей даты.
	return pgsql.GetTimetableLectur(strings.Title(data["lecture"][0]), week_int, day_int)
}

// get_auditorium function: получает аудиторию и запрашивает ее расписание из базы данных по заданным дням и неделе.
func get_auditorium(context *gin.Context) (pgsql.DataAuditoriumRequests, error) {
	// Получаем параметры запроса из URL.
	data := context.Request.URL.Query()
	// Получаем номер текущей недели и дня недели.
	week_int, day_int := getWeekAndDay(context)
	// Получаем расписание из базы данных на основе параметров запроса и текущей даты.
	return pgsql.GetTimetableAuditorium(strings.ToUpper(data["auditorium"][0]), week_int, day_int)
}

// Функция получает контекст запроса и возвращает номер недели и дня в числовом формате
func getWeekAndDay(context *gin.Context) (int, int) {
	// Извлекаем данные из URL-запроса
	data := context.Request.URL.Query()
	week, day := "", ""
	// Если параметры "week" и "day" переданы в запросе, сохраняем их
	if len(data["week"]) > 0 {
		week = data["week"][0]
	}
	if len(data["day"]) > 0 {
		day = data["day"][0]
	}

	// Преобразуем строковые значения в числа
	week_int, _ := strconv.Atoi(week)
	day_int, _ := strconv.Atoi(day)

	// Если номер недели или дня не входят в допустимый диапазон, устанавливаем значение 0
	if week_int > 2 || week_int < 1 {
		week_int = 0
	}
	if day_int > 7 || day_int < 1 {
		day_int = 0
	}
	// Возвращаем значения в числовом формате
	return week_int, day_int
}
