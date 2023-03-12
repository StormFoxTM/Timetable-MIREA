package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"knocker/pgsql"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	router := gin.Default()
	router.GET("/api/timetable", func(context *gin.Context) {
		get_function(context)
	})
	router.POST("/api/timetable", func(context *gin.Context) {
		post_function(context)
	})
	err := router.Run(":9888")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}

func get_function(context *gin.Context) {
	data := context.Request.URL.Query()
	group := len(data["group"])
	lectur := len(data["lectur"])
	auditorium := len(data["auditorium"])

	switch 1 {
	case group:
		dataGroup, err := get_group(context)
		dataGroupJSON, err_json := json.Marshal(dataGroup)
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataGroupJSON))
		} else if err == pgx.ErrNoRows {
			context.String(http.StatusBadRequest, "По запросу ничего не найдено")
		} else {
			context.String(http.StatusBadRequest, "Ошибка в запросе!")
		}

	case lectur:
		dataLectur, err := get_lecturer(context)
		dataJSON, err_json := json.Marshal(dataLectur)
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataJSON))
		} else if err == pgx.ErrNoRows {
			context.String(http.StatusBadRequest, "По запросу ничего не найдено")
		} else {
			context.String(http.StatusBadRequest, "Ошибка в запросе!")
		}

	case auditorium:
		dataAuditorium, err := get_auditorium(context)
		dataJSON, err_json := json.Marshal(dataAuditorium)
		if err == nil && err_json == nil {
			context.String(http.StatusOK, string(dataJSON))
		} else if err == pgx.ErrNoRows {
			context.String(http.StatusBadRequest, "По запросу ничего не найдено")
		} else {
			context.String(http.StatusBadRequest, "Ошибка в запросе!")
		}

	default:
		context.String(http.StatusBadRequest, "Ошибка в запросе!")
	}

}

func get_group(context *gin.Context) (pgsql.DataGroupRequests, error) {
	data := context.Request.URL.Query()
	week_int, day_int := getWeekAndDay(context)
	return pgsql.GetTimetableGroup(strings.ToUpper(data["group"][0]), week_int, day_int)
}

func get_lecturer(context *gin.Context) (pgsql.DataLecturRequests, error) {
	data := context.Request.URL.Query()
	week_int, day_int := getWeekAndDay(context)
	return pgsql.GetTimetableLectur(strings.Title(data["lectur"][0]), week_int, day_int)
}

func get_auditorium(context *gin.Context) (pgsql.DataAuditoriumRequests, error) {
	data := context.Request.URL.Query()
	week_int, day_int := getWeekAndDay(context)
	return pgsql.GetTimetableAuditorium(strings.ToUpper(data["auditorium"][0]), week_int, day_int)
}

func getWeekAndDay(context *gin.Context) (int, int) {
	data := context.Request.URL.Query()
	week, day := "", ""
	if len(data["week"]) > 0 {
		week = data["week"][0]
	}
	if len(data["day"]) > 0 {
		day = data["day"][0]
	}
	week_int, _ := strconv.Atoi(week)
	day_int, _ := strconv.Atoi(day)
	return week_int, day_int
}

func post_function(context *gin.Context) {
	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		// Handle error
	}
	data := string(jsonData)
	log.Println(data)

	// group := len(data["group"])
	// lecturer := len(data["lecturer"])
	// auditorium := len(data["auditorium"])
	// subject := len(data["subject"])
	// number := len(data["number"])
	// type_of_subject := len(data["type_of_subject"])
	// log.Println(group, subject, lecturer, auditorium, number, type_of_subject)

	context.String(http.StatusOK, "Принял")
}
