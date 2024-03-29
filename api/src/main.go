package main
/// Основной модуль в API ///

// Импортируем необходимые модули и библиотеки
import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"knocker/pgsql" // Импортируем модуль pgsql из нашего собственного пакета knocker
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

/// Основная функция main ///
func main() {
	// Создаем новый роутер Gin с настройками по умолчанию
	router := gin.Default()

	router.Use(addHeaders)

	// Добавляем обработчик GET-запроса по пути "/api/timetable"
	router.GET("/api/timetable", getFunction)

	// Добавляем обработчик GET-запроса по пути "/api/info"
	router.GET("/api/info", getInfo)

	// Добавляем обработчик GET-запроса по пути "/api/users"
	router.GET("/api/users", getUsers)

	// Добавляем обработчик POST-запроса по пути "/api/users"
	router.POST("/api/users", postUsers)

	// Запускаем сервер Gin на порту 9888
	err := router.Run(":9888")
	if err != nil {
		// В случае ошибки выводим сообщение об ошибке и завершаем программу
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

// addHeaders - промежуточное ПО для добавления заголовков OPTIONS
func addHeaders(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Access-Control-Allow-Headers", "Content-Type")
	context.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

	// Проверяем, является ли запрос OPTIONS
	if context.Request.Method == "OPTIONS" {
		context.AbortWithStatus(204) // Отвечаем 204 No Content для запросов OPTIONS
		return
	}

	context.Next()
}

/// getUsers - функция для авторизации пользователя ///
func getUsers(context *gin.Context) {
	// Получаем параметры запроса
	data := context.Request.URL.Query()
	// Получаем количество параметров "username" и "password"
	username := data["username"]
	password_enter := data["password"]
	if len(username) == 1 && len(password_enter) == 1 {
		log.Println(username[0], password_enter[0])
		password, role, err := pgsql.GetUserData(username[0])
		if password_enter[0] == password && err == nil {
			context.String(http.StatusOK, string(role))
		} else if password_enter[0] != password{
			context.String(http.StatusOK, "Wrong password")
		} else {
			context.String(http.StatusOK, "User not found")
		}
		return

	}
	context.String(http.StatusOK, "Failure")
}

/// postUsers - функция для регистрации пользователя ///
func postUsers(context *gin.Context) {
	var user User
    context.BindJSON(&user)
	log.Println("Регистрация")
	// Получаем количество параметров "username" и "password"
	username := user.Username
	password := user.Password
	log.Println(username, password)
	if username != "" && password != "" {
		err := pgsql.AddUser(username, password)
		if err == nil {
			context.String(http.StatusOK, "Success")
		} else {
			context.String(http.StatusOK, "Failure")
		}
		return

	}
	context.String(http.StatusOK, "Failure")
}

func CreateUser(c *gin.Context) {
    var user User

    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

    // Вы можете использовать полученные данные пользователя здесь

    c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно создан"})
}


/// getInfo - функция для выдачи информации ///
func getInfo(context *gin.Context) {
	// Получаем параметры запроса
	data := context.Request.URL.Query()

	// Получаем количество параметров "group", "lecturer" и "auditorium"
	group := len(data["group"])
	lecturer := len(data["lecturer"])
	auditorium := len(data["auditorium"])

	// В зависимости от того, сколько параметров указано, выбираем соответствующую функцию
	switch 1 {
	case group:
		// Если не возникло ошибок, возвращаем успех
		if err := check_group(context); err == nil {
			context.String(http.StatusOK, "Success")
		} else if err == errors.New("error connect to db") {
			context.String(http.StatusOK, "DBMS error")
		} else {
			// Если возникли ошибки, возвращаем ошибку"
			context.String(http.StatusOK, "Failure!")
		}
		return

	case lecturer:
		// Если не возникло ошибок, возвращаем успех
		if err := check_lecturer(context); err == nil {
			context.String(http.StatusOK, "Success")
		} else if err == errors.New("error connect to db") {
			context.String(http.StatusOK, "DBMS error")
		} else {
			// Если возникли ошибки, возвращаем ошибку"
			context.String(http.StatusOK, "Failure")
		}
		return

	case auditorium:
		// Если не возникло ошибок, возвращаем успех
		if err := check_auditorium(context); err == nil {
			context.String(http.StatusOK, "Success")
		} else if err == errors.New("error connect to db") {
			context.String(http.StatusOK, "DBMS error")
		} else {
			// Если возникли ошибки, возвращаем ошибку"
			context.String(http.StatusOK, "Failure")
		}
		return

	// Если количество параметров не соответствует ни одному из вариантов, возвращаем ошибку "Ошибка в запросе!"
	default:
		context.String(http.StatusOK, "Request error")
	}
}

/// getFunction - функция для выдачи расписания ///
func getFunction(context *gin.Context) {
	// Получаем параметры запроса
	data := context.Request.URL.Query()

	// Получаем количество параметров "group", "lecturer" и "auditorium"
	group := len(data["group"])
	lecturer := len(data["lecturer"])
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
			context.String(http.StatusOK, "Nothing was found on the request")
		} else {
			// Если возникли ошибки, возвращаем ошибку "Ошибка в запросе!"
			context.String(http.StatusOK, "Request error")
		}

	case lecturer:
		// Получаем данные по лектору и кодируем их в JSON
		dataLectur, err := get_lecturer(context)
		dataJSON, err_json := json.Marshal(dataLectur)

		// Если не возникло ошибок, возвращаем данные по лектору в формате JSON
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataJSON))
		} else if err == pgx.ErrNoRows {
			// Если нет данных, возвращаем ошибку "По запросу ничего не найдено"
			context.String(http.StatusOK, "Nothing was found on the request")
		} else {
			// Если возникли ошибки, возвращаем ошибку "Ошибка в запросе!"
			context.String(http.StatusOK, "Request error")
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
			context.String(http.StatusOK, "Nothing was found on the request")
		} else {
			// Если возникли ошибки, возвращаем ошибку "Ошибка в запросе!"
			context.String(http.StatusOK, "Request error")
		}

	// Если количество параметров не соответствует ни одному из вариантов, возвращаем ошибку "Ошибка в запросе!"
	default:
		context.String(http.StatusOK, "Request error")
	}
}

/// get_group function: получает группу и запрашивает ее расписание из базы данных по заданным дням и неделе. ///
func get_group(context *gin.Context) (pgsql.DataGroupRequests, error) {
	// Получаем параметры запроса из URL.
	data := context.Request.URL.Query()
	// Получаем номер текущей недели и дня недели.
	week_int, day_int := getWeekAndDay(context)
	// Получаем расписание из базы данных на основе параметров запроса и текущей даты.
	return pgsql.GetTimetableGroup(strings.ToUpper(data["group"][0]), week_int, day_int)
}

/// get_lecturer function: получает преподавателя и запрашивает его расписание из базы данных по заданным дням и неделе. ///
func get_lecturer(context *gin.Context) (pgsql.DataLecturRequests, error) {
	// Получаем параметры запроса из URL.
	data := context.Request.URL.Query()
	// Получаем номер текущей недели и дня недели.
	week_int, day_int := getWeekAndDay(context)
	// Получаем расписание из базы данных на основе параметров запроса и текущей даты.
	return pgsql.GetTimetableLectur(strings.Title(data["lecturer"][0]), week_int, day_int)
}

/// get_auditorium function: получает аудиторию и запрашивает ее расписание из базы данных по заданным дням и неделе. ///
func get_auditorium(context *gin.Context) (pgsql.DataAuditoriumRequests, error) {
	// Получаем параметры запроса из URL.
	data := context.Request.URL.Query()
	// Получаем номер текущей недели и дня недели.
	week_int, day_int := getWeekAndDay(context)
	// Получаем расписание из базы данных на основе параметров запроса и текущей даты.
	return pgsql.GetTimetableAuditorium(strings.ToUpper(data["auditorium"][0]), week_int, day_int)
}

/// getWeekAndDay - функция получает контекст запроса и возвращает номер недели и дня в числовом формате ///
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

/// check_group - функция проверки группы на наличие в расписании ///
func check_group(context *gin.Context) error {
	// Получаем параметры запроса
	data := context.Request.URL.Query()
	_, err := pgsql.GetGroupID(strings.TrimSpace(strings.ToUpper(data["group"][0])))

	return err
}

/// check_lecturer - функция проверки преподавателя на присутствие в расписании ///
func check_lecturer(context *gin.Context) error {
	// Получаем параметры запроса
	data := context.Request.URL.Query()
	_, err := pgsql.GetLecturerID(strings.TrimSpace(strings.Title(data["lecturer"][0])))

	return err
}

/// check_auditorium - функция проверки аудитории на наличие в расписании ///
func check_auditorium(context *gin.Context) error {
	// Получаем параметры запроса
	data := context.Request.URL.Query()
	err := pgsql.CheckAuditorium(strings.TrimSpace(strings.ToUpper(data["auditorium"][0])))

	return err
}
