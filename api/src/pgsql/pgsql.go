package pgsql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// тип для работы с таблицей преподавателей
type lecturer struct {
	id_lecturer int
	full_name   string
}

// тип для работы с таблицей расписания звонков
type call_schedule struct {
	subject_number int
	time_start     time.Time
	time_end       time.Time
}

// тип для работы с таблицей курсов
type course struct {
	id_of_course int
}

// тип для работы с таблицей степеней
type degree struct {
	id_of_degree    int
	degree_of_study string
}

// тип для работы с таблицей институтов
type institute struct {
	id_of_the_institute   int
	name_of_the_institute string
}

// тип для работы с таблицей учебных групп
type study_group struct {
	id_group     int
	group_name   string
	id_institute int
	id_of_course int
	id_degree    int
}

// тип для работы с таблицей расписания занятий
type timetable struct {
	id_to_group       int
	subject_to_number int
	id_lectur         int
	subject_title     string
	auditorium        string
	day_week          int
	type_of_week      int
}

// тип данных для ответа на запросы по группам
type DataGroupRequests struct {
	Weeks []DataForWeekGroup `json:"weeks"`
}

// тип данных для ответа на запросы по группам за конкретную неделю
type DataForWeekGroup struct {
	Day []DataForDayGroup `json:"day"`
}

// тип данных для ответа на запросы по группам за конкретный день
type DataForDayGroup struct {
	Subject_to_number int    `json:"subject_to_number"`
	Subject_title     string `json:"subject_title"`
	Name_lectur       string `json:"name_lectur"`
	Auditorium        string `json:"auditorium"`
}

// тип данных для ответа на запросы по преподавателям
type DataLecturRequests struct {
	Weeks []DataForWeekLectur `json:"weeks"`
}

// тип данных для ответа на запросы по преподавателям за конкретную неделю
type DataForWeekLectur struct {
	Day []DataForDayLectur `json:"day"`
}

// тип данных для ответа на запросы по преподавателям за конкретный день
type DataForDayLectur struct {
	Subject_to_number int    `json:"subject_to_number"`
	Name_group        string `json:"name_group"`
	Subject_title     string `json:"subject_title"`
	Auditorium        string `json:"auditorium"`
}

// тип данных для ответа на запросы по аудиториям
type DataAuditoriumRequests struct {
	Weeks []DataForWeekAuditorium `json:"weeks"`
}

// тип данных для ответа на запросы по аудиториям за конкретную неделю
type DataForWeekAuditorium struct {
	Day []DataForDayAuditorium `json:"day"`
}

// тип данных для ответа на запросы по аудиториям за конкретный день
type DataForDayAuditorium struct {
	Subject_to_number int    `json:"subject_to_number"`
	Name_group        string `json:"name_group"`
	Subject_title     string `json:"subject_title"`
	Name_lectur       string `json:"name_lectur"`
}

// Функция подключения к БД PostgreSQL
func connectToDB() (*pgxpool.Pool, error) {
	// Строка подключения к БД
	urlDB := "postgres://admin:admin@postgres:5432/TimeTableDB"
	// Подключение к БД
	db, err := pgxpool.New(context.Background(), urlDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Функция GetTimetableGroup возвращает расписание занятий для группы по названию group на указанный тип недели и день недели.
// type_week - тип недели (1 - нечётная, 2 - чётная, 0 - все недели)
// day_week - день недели (1 - понедельник, 2 - вторник, ..., 7 - воскресенье, 0 - все дни)
func GetTimetableGroup(group string, type_week int, day_week int) (DataGroupRequests, error) {
	var dataGroupRequests DataGroupRequests
	// Подключение к БД
	db, err := connectToDB()
	// Получение ID группы
	id_group, group_err := GetGroupID(group)

	start_type_of_week := type_week
	end_type_of_week := type_week

	start_day_of_week := day_week
	end_day_of_week := day_week
	// Обработка случаев, когда тип недели или день недели не указаны
	if type_week == 0 {
		start_type_of_week = 1
		end_type_of_week = 2
	}
	if day_week == 0 {
		start_day_of_week = 1
		end_day_of_week = 7
	}

	if err == nil && group_err == nil {
		defer db.Close()
		count := 0
		count_err := 0

		for type_of_week := start_type_of_week; type_of_week <= end_type_of_week; type_of_week += 1 {
			for day_of_week := start_day_of_week; day_of_week <= end_day_of_week; day_of_week += 1 {
				var dataForWeek DataForWeekGroup
				for subject_to_number := 1; subject_to_number <= 8; subject_to_number += 1 {
					var TimeTable timetable
					// Запрос на получение расписания для указанной группы, номера пары, дня недели и типа недели
					err = db.QueryRow(context.Background(),
						"SELECT * FROM timetable WHERE id_to_group=$1 AND subject_to_number=$2 AND day_week=$3 AND type_of_week=$4;",
						id_group, subject_to_number, day_of_week, type_of_week).Scan(&TimeTable.id_to_group, &TimeTable.subject_to_number,
						&TimeTable.id_lectur, &TimeTable.subject_title, &TimeTable.auditorium, &TimeTable.day_week, &TimeTable.type_of_week)

					if err != nil {
						count_err += 1
					} else {
						var dataForDay DataForDayGroup
						dataForDay.Auditorium = strings.TrimSpace(TimeTable.auditorium)
						dataForDay.Subject_title = strings.TrimSpace(TimeTable.subject_title)
						dataForDay.Subject_to_number = TimeTable.subject_to_number
						dataForDay.Name_lectur, _ = GetLecturerName(TimeTable.id_lectur)
						dataForWeek.Day = append(dataForWeek.Day, dataForDay)
					}
					count += 1
				}
				dataGroupRequests.Weeks = append(dataGroupRequests.Weeks, dataForWeek)
			}
		}
		// Если количество ошибок не равно количеству запросов, то вернуть nil вместо ошибки
		if count_err != count {
			err = nil
		}
	}
	return dataGroupRequests, err
}

// Функция GetTimetableLectur возвращает расписание занятий для преподавателя по имени lecturer на указанный тип недели и день недели.
// type_week - тип недели (1 - нечётная, 2 - чётная, 0 - все недели)
// day_week - день недели (1 - понедельник, 2 - вторник, ..., 7 - воскресенье, 0 - все дни)
func GetTimetableLectur(lecturer string, type_week int, day_week int) (DataLecturRequests, error) {
	var dataRequests DataLecturRequests
	// Подключение к БД
	db, err := connectToDB()
	// Получение id преподавателя
	id_lectur, lectur_err := GetLecturerID(lecturer)

	// Определение типа и диапазона недель, если заданы все недели или все дни
	start_type_of_week := type_week
	end_type_of_week := type_week
	start_day_of_week := day_week
	end_day_of_week := day_week

	// Обработка случаев, когда тип недели или день недели не указаны
	if type_week == 0 {
		start_type_of_week = 1
		end_type_of_week = 2
	}
	if day_week == 0 {
		start_day_of_week = 1
		end_day_of_week = 7
	}

	if err == nil && lectur_err == nil {
		defer db.Close()
		count := 0     // счетчик занятий
		count_err := 0 // счетчик ошибок

		// Обход всех недель и дней
		for type_of_week := start_type_of_week; type_of_week <= end_type_of_week; type_of_week += 1 {
			for day_of_week := start_day_of_week; day_of_week <= end_day_of_week; day_of_week += 1 {
				var dataForWeek DataForWeekLectur
				// Обход всех пар на заданный день и тип недели
				for subject_to_number := 1; subject_to_number <= 8; subject_to_number += 1 {
					var TimeTable timetable
					// Выборка занятий по id преподавателя, номеру пары, дню недели и типу недели
					err = db.QueryRow(context.Background(),
						"SELECT * FROM timetable WHERE id_lectur=$1 AND subject_to_number=$2 AND day_week=$3 AND type_of_week=$4;",
						id_lectur, subject_to_number, day_of_week, type_of_week).Scan(&TimeTable.id_to_group, &TimeTable.subject_to_number,
						&TimeTable.id_lectur, &TimeTable.subject_title, &TimeTable.auditorium, &TimeTable.day_week, &TimeTable.type_of_week)
					if err != nil {
						count_err += 1
					} else {
						var dataForDay DataForDayLectur
						// Заполнение данных для одной пары
						dataForDay.Auditorium = strings.TrimSpace(TimeTable.auditorium)
						dataForDay.Subject_title = strings.TrimSpace(TimeTable.subject_title)
						dataForDay.Subject_to_number = TimeTable.subject_to_number
						dataForDay.Name_group, _ = GetGroupName(TimeTable.id_to_group)
						dataForWeek.Day = append(dataForWeek.Day, dataForDay)
					}
					count += 1
				}
				dataRequests.Weeks = append(dataRequests.Weeks, dataForWeek)
			}
		}
		// Если количество ошибок не равно количеству запросов, то вернуть nil вместо ошибки
		if count_err != count {
			err = nil
		}
	}
	return dataRequests, err
}

// Функция GetTimetableAuditorium возвращает расписание занятий для преподавателя по названию auditorium на указанный тип недели и день недели.
// type_week - тип недели (1 - нечётная, 2 - чётная, 0 - все недели)
// day_week - день недели (1 - понедельник, 2 - вторник, ..., 7 - воскресенье, 0 - все дни)
func GetTimetableAuditorium(auditorium string, type_week int, day_week int) (DataAuditoriumRequests, error) {
	var dataRequests DataAuditoriumRequests
	// Подключение к БД
	db, err := connectToDB()

	// Определение типа и диапазона недель, если заданы все недели или все дни
	start_type_of_week := type_week
	end_type_of_week := type_week

	start_day_of_week := day_week
	end_day_of_week := day_week

	// Обработка случаев, когда тип недели или день недели не указаны
	if type_week == 0 {
		start_type_of_week = 1
		end_type_of_week = 2
	}
	if day_week == 0 {
		start_day_of_week = 1
		end_day_of_week = 7
	}
	if err == nil {
		defer db.Close()
		count := 0
		count_err := 0

		for type_of_week := start_type_of_week; type_of_week <= end_type_of_week; type_of_week += 1 {
			for day_of_week := start_day_of_week; day_of_week <= end_day_of_week; day_of_week += 1 {
				var dataForWeek DataForWeekAuditorium
				for subject_to_number := 1; subject_to_number <= 8; subject_to_number += 1 {
					var TimeTable timetable
					err = db.QueryRow(context.Background(),
						"SELECT * FROM timetable WHERE auditorium LIKE $1 AND subject_to_number=$2 AND day_week=$3 AND type_of_week=$4;",
						"%"+auditorium+"%", subject_to_number, day_of_week, type_of_week).Scan(&TimeTable.id_to_group, &TimeTable.subject_to_number,
						&TimeTable.id_lectur, &TimeTable.subject_title, &TimeTable.auditorium, &TimeTable.day_week, &TimeTable.type_of_week)
					if err != nil {
						count_err += 1
					} else {
						var dataForDay DataForDayAuditorium
						// Заполнение данных для одной пары
						dataForDay.Name_group, _ = GetGroupName(TimeTable.id_to_group)
						dataForDay.Subject_title = strings.TrimSpace(TimeTable.subject_title)
						dataForDay.Subject_to_number = TimeTable.subject_to_number
						dataForDay.Name_lectur, _ = GetLecturerName(TimeTable.id_lectur)
						dataForWeek.Day = append(dataForWeek.Day, dataForDay)
					}
					count += 1
				}
				dataRequests.Weeks = append(dataRequests.Weeks, dataForWeek)
			}
		}
		// Если количество ошибок не равно количеству запросов, то вернуть nil вместо ошибки
		if count_err != count {
			err = nil
		}
	}
	return dataRequests, err
}

// Функция получения имени преподавателя по его идентификатору
func GetLecturerName(lecturer_id int) (string, error) {
	// Открытие соединения с БД
	var lecturer_name string
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		// Выполнение запроса на выборку имени преподавателя из БД
		err = db.QueryRow(context.Background(), "SELECT full_name FROM lecturer WHERE id_lecturer=$1;", lecturer_id).Scan(&lecturer_name)
		if err == nil {
			return strings.TrimSpace(lecturer_name), nil
		}
	}
	return "", err
}

// Функция получения названия группы по ее идентификатору
func GetGroupName(group_id int) (string, error) {
	// Открытие соединения с БД
	var group_name string
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		// Выполнение запроса на выборку названия группы из БД
		err = db.QueryRow(context.Background(), "SELECT group_name FROM study_group WHERE id_group=$1;", group_id).Scan(&group_name)
		if err == nil {
			return strings.TrimSpace(group_name), nil
		}
	}
	return "", err
}

// Функция получения идентификатора преподавателя по его имени
func GetLecturerID(lecturer_name string) (int, error) {
	// Открытие соединения с БД
	var id_lecturer int
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		// Выполнение запроса на выборку идентификатора преподавателя из БД
		err = db.QueryRow(context.Background(), "SELECT id_lecturer FROM lecturer WHERE full_name=$1;", lecturer_name).Scan(&id_lecturer)
		if err == nil {
			return id_lecturer, nil
		}
		return 0, err
	}
	return 0, errors.New("error connect to db")
}

// Функция получения идентификатора группы по ее названию
func GetGroupID(group string) (int, error) {
	// Открытие соединения с БД
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		var group_ID int
		// Выполнение запроса на выборку идентификатора группы из БД
		err = db.QueryRow(context.Background(), "SELECT id_group FROM study_group WHERE group_name=$1;", group).Scan(&group_ID)
		if err == nil {
			return group_ID, nil
		}
		return 0, err
	}
	return 0, errors.New("error connect to db")
}

// Функция проверки аудитории на наличие в расписании
func CheckAuditorium(auditorium string) error {
	// Открытие соединения с БД
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		var group_ID int
		// Выполнение запроса на выборку идентификатора группы из БД
		err = db.QueryRow(context.Background(), "SELECT id_to_group FROM timetable WHERE auditorium=$1;", auditorium).Scan(&group_ID)
		return err
	}
	return errors.New("error connect to db")
}

// Функция получения идентификатора института по его названию
func GetInstituteID(instituteName string) (int, error) {
	// Открытие соединения с БД
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		var id_of_the_institute int
		// Выполнение запроса на выборку идентификатора института из БД
		err_institute := db.QueryRow(context.Background(), "SELECT id_of_the_institute FROM institute WHERE name_of_the_institute=$1;",
			instituteName).Scan(&id_of_the_institute)
		if err_institute == nil {
			return id_of_the_institute, nil
		}
	}
	return 0, err
}

// Функция получает название группы и вырезает из неё 3 символ
func find_letter(group string) int {
	var letter rune
	for i, elem := range group {
		if i == 4 {
			letter = elem
			break
		}
	}
	return find_id_degree(string(letter))
}

// Функция получает 3 символ группы, который определяет степень обучения
func find_id_degree(letter string) int {
	switch letter {
	case "Б":
		return 1
	case "С":
		return 2
	case "М":
		return 3
	default:
		return 4
	}
}
