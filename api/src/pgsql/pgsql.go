package pgsql

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type lecturer struct {
	id_lecturer int
	full_name   string
}

type call_schedule struct {
	subject_number int
	time_start     time.Time
	time_end       time.Time
}

type course struct {
	id_of_course int
}

type degree struct {
	id_of_degree    int
	degree_of_study string
}

type institute struct {
	id_of_the_institute   int
	name_of_the_institute string
}

type study_group struct {
	id_group     int
	group_name   string
	id_institute int
	id_of_course int
	id_degree    int
}

type timetable struct {
	id_to_group       int
	subject_to_number int
	id_lectur         int
	subject_title     string
	auditorium        string
	day_week          int
	type_of_week      int
}

type DataGroupRequests struct {
	Weeks []DataForWeekGroup `json:"weeks"`
}

type DataForWeekGroup struct {
	Day []DataForDayGroup `json:"day"`
}

type DataForDayGroup struct {
	Subject_to_number int    `json:"subject_to_number"`
	Subject_title     string `json:"subject_title"`
	Name_lectur       string `json:"name_lectur"`
	Auditorium        string `json:"auditorium"`
}

type DataLecturRequests struct {
	Weeks []DataForWeekLectur `json:"weeks"`
}

type DataForWeekLectur struct {
	Day []DataForDayLectur `json:"day"`
}

type DataForDayLectur struct {
	Subject_to_number int    `json:"subject_to_number"`
	Name_group        string `json:"name_group"`
	Subject_title     string `json:"subject_title"`
	Auditorium        string `json:"auditorium"`
}

type DataAuditoriumRequests struct {
	Weeks []DataForWeekAuditorium `json:"weeks"`
}

type DataForWeekAuditorium struct {
	Day []DataForDayAuditorium `json:"day"`
}

type DataForDayAuditorium struct {
	Subject_to_number int    `json:"subject_to_number"`
	Name_group        string `json:"name_group"`
	Subject_title     string `json:"subject_title"`
	Name_lectur       string `json:"name_lectur"`
}

func connectToDB() (*pgxpool.Pool, error) {
	urlDB := "postgres://admin:admin@postgres:5432/TimeTableDB"
	db, err := pgxpool.New(context.Background(), urlDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetTimetableGroup(group string, type_week int, day_week int) (DataGroupRequests, error) {
	var dataGroupRequests DataGroupRequests
	db, err := connectToDB()
	id_group, group_err := GetGroupID(group)

	start_type_of_week := type_week
	end_type_of_week := type_week

	start_day_of_week := day_week
	end_day_of_week := day_week
	if type_week == 0 {
		start_type_of_week = 1
		end_type_of_week = 2
	}
	if day_week == 0 {
		start_day_of_week = 1
		end_day_of_week = 6
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
		if count_err != count {
			err = nil
		}
	}
	return dataGroupRequests, err
}

func GetTimetableLectur(lecturer string, type_week int, day_week int) (DataLecturRequests, error) {
	var dataRequests DataLecturRequests
	db, err := connectToDB()
	id_lectur, lectur_err := GetLecturerID(lecturer)

	start_type_of_week := type_week
	end_type_of_week := type_week

	start_day_of_week := day_week
	end_day_of_week := day_week
	if type_week == 0 {
		start_type_of_week = 1
		end_type_of_week = 2
	}
	if day_week == 0 {
		start_day_of_week = 1
		end_day_of_week = 6
	}
	log.Println(lecturer)
	if err == nil && lectur_err == nil {
		defer db.Close()
		count := 0
		count_err := 0

		for type_of_week := start_type_of_week; type_of_week <= end_type_of_week; type_of_week += 1 {
			for day_of_week := start_day_of_week; day_of_week <= end_day_of_week; day_of_week += 1 {
				var dataForWeek DataForWeekLectur
				for subject_to_number := 1; subject_to_number <= 8; subject_to_number += 1 {
					var TimeTable timetable
					err = db.QueryRow(context.Background(),
						"SELECT * FROM timetable WHERE id_lectur=$1 AND subject_to_number=$2 AND day_week=$3 AND type_of_week=$4;",
						id_lectur, subject_to_number, day_of_week, type_of_week).Scan(&TimeTable.id_to_group, &TimeTable.subject_to_number,
						&TimeTable.id_lectur, &TimeTable.subject_title, &TimeTable.auditorium, &TimeTable.day_week, &TimeTable.type_of_week)
					if err != nil {
						count_err += 1
					} else {
						var dataForDay DataForDayLectur
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
		if count_err != count {
			err = nil
		}
	}
	return dataRequests, err
}

func GetTimetableAuditorium(auditorium string, type_week int, day_week int) (DataAuditoriumRequests, error) {
	var dataRequests DataAuditoriumRequests
	db, err := connectToDB()

	start_type_of_week := type_week
	end_type_of_week := type_week

	start_day_of_week := day_week
	end_day_of_week := day_week
	if type_week == 0 {
		start_type_of_week = 1
		end_type_of_week = 2
	}
	if day_week == 0 {
		start_day_of_week = 1
		end_day_of_week = 6
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
		if count_err != count {
			err = nil
		}
	}
	return dataRequests, err
}

func GetLecturerName(lecturer_id int) (string, error) {
	var lecturer_name string
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		err = db.QueryRow(context.Background(), "SELECT full_name FROM lecturer WHERE id_lecturer=$1;", lecturer_id).Scan(&lecturer_name)
		if err == nil {
			return strings.TrimSpace(lecturer_name), nil
		}
	}
	return "", err
}

func GetGroupName(group_id int) (string, error) {
	var group_name string
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		err = db.QueryRow(context.Background(), "SELECT group_name FROM study_group WHERE id_group=$1;", group_id).Scan(&group_name)
		if err == nil {
			return strings.TrimSpace(group_name), nil
		}
	}
	return "", err
}

func GetLecturerID(lecturer_name string) (int, error) {
	var id_lecturer int
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		err = db.QueryRow(context.Background(), "SELECT id_lecturer FROM lecturer WHERE full_name=$1;", lecturer_name).Scan(&id_lecturer)
		if err == nil {
			return id_lecturer, nil
		}
	}
	return 0, err
}

func GetGroupID(group string) (int, error) {
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		var group_ID int

		err = db.QueryRow(context.Background(), "SELECT id_group FROM study_group WHERE group_name=$1;", group).Scan(&group_ID)
		if err == nil {
			return group_ID, nil
		}
	}
	return 0, nil
}

func GetInstituteID(instituteName string) (int, error) {
	db, err := connectToDB()
	if err == nil {
		defer db.Close()
		var id_of_the_institute int
		err_institute := db.QueryRow(context.Background(), "SELECT id_of_the_institute FROM institute WHERE name_of_the_institute=$1;",
			instituteName).Scan(&id_of_the_institute)
		if err_institute == nil {
			return id_of_the_institute, nil
		}
	}
	return 0, err
}

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
