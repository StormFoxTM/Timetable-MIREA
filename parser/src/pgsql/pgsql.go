package pgsql

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
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

// основная функция, которая вызывается для добавления определённых данных в СУБД
func MainFunc(group string, day int, type_of_week int, number string,
	subject string, lecturer string, auditorium string, type_of_subject string, institute string, course string) {
	urlDB := "postgres://admin:admin@postgres:5432/TimeTableDB" // адрес СУБД
	db, err := pgx.Connect(context.Background(), urlDB)         // подключение к СУБД
	if err != nil {
		fmt.Errorf("Error to connect db", err)
	} else {
		defer db.Close(context.Background())  // закрытие соединения с СУБД
		int_course, _ := strconv.Atoi(course) // превращения курса в число
		addTimetable(db, group, day, type_of_week, number, subject, lecturer, auditorium, type_of_subject, institute, int_course+1)
	}
}

// добавление расписания в СУБД
func addTimetable(db *pgx.Conn, group string, day_week int, type_of_week int, subject_to_number string,
	subject_title string, name_lecturer string, auditorium string, type_of_subject string, institute string, int_course int) {
	var id_lecturer int
	var StudyGroup study_group
	id_lecturer = addLecturer(db, name_lecturer)

	err := db.QueryRow(context.Background(), "SELECT * FROM study_group WHERE group_name=$1;", group).Scan(&StudyGroup.id_group,
		&StudyGroup.group_name, &StudyGroup.id_institute, &StudyGroup.id_of_course, &StudyGroup.id_degree)
	if err == pgx.ErrNoRows || GetInstituteName(db, StudyGroup.id_institute) != institute || int_course != StudyGroup.id_of_course {
		StudyGroup.id_group = addGroup(db, group, institute, int_course, err, StudyGroup.id_group)
	}

	if StudyGroup.id_group != 0 {
		err_duplicate := db.QueryRow(context.Background(),
			"SELECT id_to_group FROM timetable WHERE id_to_group=$1 AND subject_to_number=$2 AND day_week=$3 AND type_of_week=$4;",
			StudyGroup.id_group, subject_to_number, day_week, type_of_week).Scan(&StudyGroup.id_group)
		if err_duplicate == pgx.ErrNoRows {
			if id_lecturer == -1 {
				db.QueryRow(context.Background(),
					"INSERT INTO timetable (id_to_group, subject_to_number, subject_title, auditorium, day_week, type_of_week) "+
						"VALUES ($1, $2, $3, $4, $5, $6);", StudyGroup.id_group, subject_to_number, subject_title, auditorium, day_week, type_of_week)
			} else {
				db.QueryRow(context.Background(),
					"INSERT INTO timetable (id_to_group, subject_to_number, id_lectur, subject_title, auditorium, day_week, type_of_week) "+
						"VALUES ($1, $2, $3, $4, $5, $6, $7);", StudyGroup.id_group, subject_to_number, id_lecturer, subject_title, auditorium,
					day_week, type_of_week)
			}
		} else {
			if id_lecturer == -1 {
				db.QueryRow(context.Background(), "UPDATE timetable SET subject_title=$1, auditorium=$2 "+
					"WHERE id_to_group=$3 AND subject_to_number=$4 AND day_week=$5 AND type_of_week=$6;", subject_title, auditorium,
					StudyGroup.id_group, subject_to_number, day_week, type_of_week)
			} else {
				db.QueryRow(context.Background(), "UPDATE timetable SET id_lectur=$1, subject_title=$2, auditorium=$3 "+
					"WHERE id_to_group=$4 AND subject_to_number=$5 AND day_week=$6 AND type_of_week=$7;", id_lecturer, subject_title,
					auditorium, StudyGroup.id_group, subject_to_number, day_week, type_of_week)
			}
		}
	}
}

// добавление преподавателя в СУБД
func addLecturer(db *pgx.Conn, lecturers string) int {
	if lecturers != "" {
		var Lecturer lecturer
		re := regexp.MustCompile(`[ЁА-Я][ёа-я]+\s[ЁёА-я].[ЁёА-я].`)
		lecturer_name := re.FindAllStringSubmatch(lecturers, -1)
		if len(lecturer_name) > 0 {
			for i := 0; i < len(lecturer_name); i += 1 {
				err := db.QueryRow(context.Background(), "SELECT id_lecturer FROM lecturer WHERE full_name=$1;", lecturer_name[i][0]).Scan(
					&Lecturer.id_lecturer)
				if err == pgx.ErrNoRows {
					db.QueryRow(context.Background(), "INSERT INTO lecturer (full_name) VALUES ($1) RETURNING id_lecturer;",
						lecturer_name[i][0]).Scan(&Lecturer.id_lecturer, &Lecturer.full_name)
				}
			}
			return Lecturer.id_lecturer
		}
	}
	return -1
}

// добавление группы в СУБД
func addGroup(db *pgx.Conn, group string, name_institute string, course int, lastErr error, id_group int) int {
	var Study_Group study_group
	var Institute institute

	err := db.QueryRow(context.Background(), "SELECT * FROM institute WHERE name_of_the_institute=$1;",
		name_institute).Scan(&Institute.name_of_the_institute)
	if err == pgx.ErrNoRows {
		Institute.id_of_the_institute = addInstitute(db, name_institute)
	}
	if lastErr == pgx.ErrNoRows {
		db.QueryRow(context.Background(),
			"INSERT INTO study_group (group_name, id_institute, id_of_course, id_degree) VALUES ($1, $2, $3, $4) RETURNING id_group;",
			group, Institute.id_of_the_institute, course, find_letter(group)).Scan(&Study_Group.id_group)
	} else {
		db.QueryRow(context.Background(),
			"UPDATE study_group SET id_institute=$1, id_of_course=$2 WHERE id_group=$3 RETURNING id_group;",
			Institute.id_of_the_institute, course, id_group).Scan(&Study_Group.id_group)
	}
	return Study_Group.id_group
}

// добавление института в СУБД
func addInstitute(db *pgx.Conn, instituteName string) int {
	var id_of_the_institute int
	db.QueryRow(context.Background(), "INSERT INTO institute (name_of_the_institute) VALUES ($1) RETURNING id_of_the_institute;",
		instituteName).Scan(&id_of_the_institute)
	return id_of_the_institute
}

// получение название института по его ID
func GetInstituteName(db *pgx.Conn, instituteID int) string {
	var name_of_the_institute string
	err_institute := db.QueryRow(context.Background(), "SELECT name_of_the_institute FROM institute WHERE id_of_the_institute=$1;",
		instituteID).Scan(&name_of_the_institute)
	if err_institute == nil {
		return name_of_the_institute
	}
	return ""
}

// поиск третий буквы в названии группы
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

// определение степени обучения, которая зашифрована в 3 букве
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
