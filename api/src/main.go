package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/timetable", func(context *gin.Context) {
		get_function(context)
	})
	router.POST("/timetable", func(context *gin.Context) {
		post_function(context)
	})
	err := router.Run(":9888")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
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

func get_function(context *gin.Context) {
	data := context.Request.URL.Query()
	group := len(data["group"])
	lecturer := len(data["lecturer"])
	auditorium := len(data["auditorium"])

	switch 1 {
	case group:
		get_group(context)
		context.String(http.StatusOK, "group")

	case lecturer:
		get_lecturer(context)
		context.String(http.StatusOK, "lecturer")

	case auditorium:
		get_auditorium(context)
		context.String(http.StatusOK, "auditorium")

	default:
		context.String(http.StatusBadRequest, "Ошибка в запросе!")
	}

}

func get_group(context *gin.Context) {
	data := context.Request.URL.Query()
	group, week, day := "", "", ""
	if len(data["group"]) > 0 {
		group = data["group"][0]
	}
	if len(data["week"]) > 0 {
		week = data["week"][0]
	}
	if len(data["day"]) > 0 {
		day = data["day"][0]
	}
	log.Println(group, week, day)
}

func get_lecturer(context *gin.Context) {

}

func get_auditorium(context *gin.Context) {

}
