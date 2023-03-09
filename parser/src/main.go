package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"knocker/pgsql"
	"knocker/xls"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
	// 	"github.com/xuri/excelize/v2"
)

type Links struct {
	institute, course, degree string
}

var (
	siteURL = "https://www.mirea.ru/schedule/"
	links   map[string]Links
	// days         = [7]string{"Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Субботы", "Воскресенье"}
	// type_of_week = [2]string{"I", "II"}
	dsn = "api:api@tcp(timetable.postgres:5432)/TimeTableDB?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	log.Println("Start parsing")
	check_directory()
	parse_site()
	loader()
	check_directory()
	log.Println("End parsing")
}

func check_directory() {
	os.RemoveAll("temp")
	os.Mkdir("temp", 0777)
	if _, err := os.Stat("last"); os.IsNotExist(err) {
		os.Mkdir("last", 0777)
	}
}

func loader() error {
	log.Println("Start loader")
	var queue_download sync.WaitGroup
	for link, data := range links {
		queue_download.Add(1)
		// go download(link, data, &queue_download)
		download(link, data, &queue_download)
	}
	queue_download.Wait()
	return nil
}

func download(url string, data Links, queue_download *sync.WaitGroup) error {
	defer queue_download.Done()
	resp, err := connect(url, 0)
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	file := strings.Join(strings.Fields(data.institute), "-") + "_" +
		strings.Join(strings.Fields(data.course), "-") + "_" + data.degree + url[strings.LastIndexAny(url, "."):]
	out, err := os.Create("temp\\" + file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if url[strings.LastIndexAny(url, ".")+1:] == "xls" {
		err = xls.Convert("temp\\" + file)
		if err != nil {
			fmt.Println("Convertor", err)
			return err
		}
		file += "x"
	}
	if !file_comparison(file) {
		if err = read_file(file, data); err != nil {
			fmt.Println(err)
			return err
		}
		if err = move_file(file); err != nil {
			log.Println(err)
			return err
		}

	}
	return nil
}

func connect(url string, count int) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		time.Sleep(1)
		resp, err = connect(url, count+1)
	}
	return resp, err
}

func parse_site() error {
	log.Println("Parse site")
	resp, err := connect(siteURL, 0)
	if err != nil {
		return fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if doc, err := goquery.NewDocumentFromReader(resp.Body); err != nil {
		return fmt.Errorf("GET error: %v", err)
	} else {
		links = make(map[string]Links)
		doc.Find("a").Each(func(i int, temp_data *goquery.Selection) {
			if temp_data.HasClass("uk-text-bold") {
				temp_data.Parent().Parent().Parent().Find("a").Each(func(j int, temp_data_of_link *goquery.Selection) {
					if band, ok := temp_data_of_link.Attr("href"); ok && temp_data_of_link.HasClass("uk-link-toggle") {
						temp_count := 0
						for _, data := range links {
							if data.institute == strings.TrimSpace(temp_data.Text()) &&
								data.course == strings.TrimSpace(temp_data_of_link.Text()) && data.degree == strconv.Itoa(temp_count) {
								temp_count += 1
							}
						}
						links[strings.TrimSpace(band)] = Links{
							strings.TrimSpace(temp_data.Text()), strings.TrimSpace(temp_data_of_link.Text()), strconv.Itoa(temp_count),
						}
					}
				})
			}
		})
		return nil
	}
}

func file_comparison(file string) bool {
	new_file, err := os.Open("temp\\" + file)
	if err != nil {
		return false
	}
	last_file, err := os.Open("last\\" + file)
	if err != nil {
		return false
	}
	defer new_file.Close()
	defer last_file.Close()
	return find_sum(new_file) == find_sum(last_file)
}

func find_sum(file *os.File) string {
	file.Seek(0, 0) // Сброс курсора к началу файла
	sum, err := getMD5SumString(file)
	if err != nil {
		panic(err)
	}
	return sum
}

func getMD5SumString(f *os.File) (string, error) {
	fileSum := md5.New()
	_, err := io.Copy(fileSum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", fileSum.Sum(nil)), nil
}

func read_file(file string, data Links) error {
	xlsxFile, err := excelize.OpenFile("temp\\" + file)
	if err != nil {
		return fmt.Errorf("Read file: %s%s%s", file, " - ", err)
	}

	defer xlsxFile.Close()
	lists := xlsxFile.GetSheetList()
	for _, sheetName := range lists {
		read_list(xlsxFile, sheetName, data)
	}
	return nil
}

func read_list(xlsxFile *excelize.File, sheetName string, data Links) {
	find_group := false
	rowIndex := 2
	re := regexp.MustCompile(`\p{L}\p{L}\p{L}\p{L}-\d\d-\d\d`)
	if find_group {
		return
	}
	count_free := 0
	for columnIndex := 1; ; columnIndex++ {
		group, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex)+strconv.Itoa(rowIndex))
		if status, _ := regexp.MatchString(`\p{L}\p{L}\p{L}\p{L}-\d\d-\d\d`, clear_str(group)); status {
			group = re.FindString(group)
			count_day := -1
			count_free = 0
			for row := rowIndex + 2; ; row += 2 {
				number, _ := xlsxFile.GetCellValue(sheetName, cell_name(2)+strconv.Itoa(row))
				if strings.TrimSpace(number) == "1" {
					count_day += 1
				}
				_, err := strconv.Atoi(number)
				if err == nil {
					read_data(xlsxFile, sheetName, columnIndex, number, count_day, group, row, row, data)
					read_data(xlsxFile, sheetName, columnIndex, number, count_day, group, row+1, row, data)

				} else {
					break
				}
			}

			count_free = 0
			if !find_group {
				find_group = true
			}
		} else if group == "" {
			count_free += 1
		}
		if count_free >= 12 {
			break
		}
	}
}

func read_data(xlsxFile *excelize.File, sheetName string, columnIndex int, number string, count_day int,
	group string, elem int, row int, data Links) {
	subject, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex)+strconv.Itoa(elem))
	type_of_subject, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex+1)+strconv.Itoa(elem))
	lecturer, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex+2)+strconv.Itoa(elem))
	auditorium, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex+3)+strconv.Itoa(elem))
	pgsql.MainFunc(group, count_day+1, elem-row+1, clear_str(number), clear_str(subject),
		clear_str(lecturer), clear_str(auditorium), clear_str(type_of_subject), data.institute, data.course)
}

func clear_str(str string) string {
	tmp := strings.TrimSpace(strings.ReplaceAll(str, "\n", " "))
	return strings.ReplaceAll(tmp, "\t", " ")
}

func cell_name(number int) string {
	name, _ := excelize.ColumnNumberToName(number)
	return name
}

func move_file(file string) error {
	os.Remove("last\\" + file)
	inputFile, err := os.Open("temp\\" + file)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create("last\\" + file)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	_, err = io.Copy(outputFile, inputFile)
	defer inputFile.Close()
	defer outputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	return nil
}
