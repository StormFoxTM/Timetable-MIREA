package main
/// Основной модуль в парсере ///

// Импортируем необходимые модули и библиотеки
import (
	"crypto/md5"
	"fmt"
	"io"
	"knocker/pgsql" // Импортируем модуль pgsql из нашего собственного пакета knocker
	"knocker/xls"   // Импортируем модуль xls из нашего собственного пакета knocker

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
)

// Определение структуры Links, хранящей строки с институтом, курсом и степенью
type Links struct {
	institute, course, degree string
}

// Определение глобальных переменных
var (
	// Адрес сайта
	siteURL = "https://www.mirea.ru/schedule/"
	// Словарь ссылок на расписание
	links map[string]Links
	// Строка подключения к БД PostgreSQL
	dsn = "api:api@tcp(timetable.postgres:5432)/TimeTableDB?charset=utf8mb4&parseTime=True&loc=Local"
)

/// Функция main - основная функция ///
func main() {
	log.Println("Парсер успешно запущен!")
	// Запуск функции для проверки здоровья контейнера
	go healthcheack()
	// Бесконечный цикл
	for {
		// Запуск процесса обновления данных
		start_script()

		// Вычисление времени до следующего запуска парсинга (каждый час)
		nextHour := time.Now().Truncate(time.Hour).Add(time.Hour)
		duration := nextHour.Sub(time.Now())
		// Задержка выполнения на время до следующего запуска
		time.Sleep(duration)
	}
}

/// healthcheack - функция для проверки здоровья контейнера ///
func healthcheack() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Fatal(http.ListenAndServe(":9009", nil))
}

/// start_script - функция для обновления данных ///
func start_script() {
	// Вывод в лог сообщения о начале парсинга
	log.Println("Запущен процесс обновления расписания")
	// Удаление каталога "temp" и его содержимого
	check_directory()
	// Парсинг сайта и сохранение полученных файлов в каталог "temp"
	parse_site()
	// Загрузка полученных файлов в базу данных
	loader()
	// Повторное удаление каталога "temp" и его содержимого
	check_directory()
	// Вывод в лог сообщения об окончании парсинга
	log.Println("End parsing")
}

/// check_directory - функция удаления каталога 'temp' и его содержимого, а также создания каталога 'last' ///
func check_directory() {
	os.RemoveAll("temp")
	os.Mkdir("temp", 0777)
	if _, err := os.Stat("last"); os.IsNotExist(err) {
		os.Mkdir("last", 0777)
	}
}

/// loader - функция загрузки файлов в базу данных ///
func loader() error {
	log.Println("Start loader")
	var queue_download sync.WaitGroup
	for link, data := range links {
		queue_download.Add(1)
		go download(link, data, &queue_download)
	}
	queue_download.Wait()
	return nil
}

/// download - функция загрузки файла по ссылке ///
func download(url string, data Links, queue_download *sync.WaitGroup) error {
	// Уменьшение счетчика ожидания на 1 при завершении работы функции
	defer queue_download.Done()
	// Установление соединения с сайтом и получение HTTP-ответа
	resp, err := connect(url, 0)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Ошибка при скачивании файла: " + url)
		return err
	}
	// Определение имени файла
	file := strings.Join(strings.Fields(data.institute), "-") + "_" +
		strings.Join(strings.Fields(data.course), "-") + "_" + data.degree + url[strings.LastIndexAny(url, "."):]
	// Создание файла в каталоге "temp" и запись в него данных из HTTP-ответа
	out, err := os.Create("temp\\" + file) // создание файла
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = io.Copy(out, resp.Body) // запись в созданный файл данных из HTTP-ответа
	out.Close()                      // закрытие файла

	// проверка, является ли файл xls и конвертация его в xlsx
	if url[strings.LastIndexAny(url, ".")+1:] == "xls" {
		err = xls.Convert("temp\\" + file) // конвертация в xlsx
		if err != nil {
			fmt.Println("Convertor", err)
			return err
		}
		file += "x" // добавление расширения xlsx к имени файла
	}

	// проверка, существует ли файл с таким именем в каталоге last и не совпадают ли хэши-суммы
	if !file_comparison(file) {
		// чтение данных из файла
		if err = read_file(file, data); err != nil {
			fmt.Println(err)
			return err
		}
		// перемещение файла из каталога temp в каталог last
		if err = move_file(file); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

/// connect - функция для установления соединения с сайтом и получение HTTP-ответа ///
func connect(url string, count int) (*http.Response, error) {
	if count >= 100 {
		return nil, fmt.Errorf("Не удалось подключиться к сайту")
	}
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		time.Sleep(1)
		resp, err = connect(url, count+1) // рекурсивный вызов, пока ответ не будет успешным
	}
	return resp, err
}

/// parse_site - функция для парсинга сайта ///
func parse_site() error {
	log.Println("Парсинг сайта")
	resp, err := connect(siteURL, 0) // получаем GET запрос
	if err != nil {
		return fmt.Errorf("Ошибка GET запроса: %v", err)
	}
	defer resp.Body.Close()

	if doc, err := goquery.NewDocumentFromReader(resp.Body); err != nil {
		return fmt.Errorf("Ошибка GET запроса: %v", err)
	} else {
		links = make(map[string]Links)
		// поиск всех элементов a и выборка нужных элементов
		doc.Find("a").Each(func(i int, temp_data *goquery.Selection) {
			// выбираем ссылки с классом "uk-text-bold"
			if temp_data.HasClass("uk-text-bold") {
				temp_data.Parent().Parent().Parent().Find("a").Each(func(j int, temp_data_of_link *goquery.Selection) {
					// выбираем ссылки с классом "uk-link-toggle" и получаем значения атрибута "href"
					if band, ok := temp_data_of_link.Attr("href"); ok && temp_data_of_link.HasClass("uk-link-toggle") {
						temp_count := 0
						for _, data := range links {
							// выбираем уникальную ссылку, определяя её по университету, курсу и степени
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

/// file_comparison - функция для сравнения файлов ///
func file_comparison(file string) bool {
	// открываем файлы
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
	// сравниваем суммы md5
	return find_sum(new_file) == find_sum(last_file)
}

/// find_sum - функция для нахождения суммы md5 ///
func find_sum(file *os.File) string {
	file.Seek(0, 0) // Сброс курсора к началу файла
	sum, err := getMD5SumString(file)
	if err != nil {
		panic(err)
	}
	return sum
}

/// getMD5SumString - функция для получения md5-хеша ///
func getMD5SumString(f *os.File) (string, error) {
	fileSum := md5.New()
	_, err := io.Copy(fileSum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", fileSum.Sum(nil)), nil
}

/// read_file - функция для чтения файла xlsx ///
func read_file(file string, data Links) error {
	xlsxFile, err := excelize.OpenFile("temp\\" + file) // открываем файл
	if err != nil {
		return fmt.Errorf("Read file: %s%s%s", file, " - ", err)
	}

	defer xlsxFile.Close() // закрываем файл после чтения

	// получаем список листов в файле и читаем их
	lists := xlsxFile.GetSheetList()
	for _, sheetName := range lists {
		read_list(xlsxFile, sheetName, data)
	}
	return nil
}

/// read_list - функция для чтения данных из списка ///
func read_list(xlsxFile *excelize.File, sheetName string, data Links) {
	find_group := false                                        // флаг, указывающий на то, была ли найдена группа
	rowIndex := 2                                              // начальная строка для поиска
	re := regexp.MustCompile(`\p{L}\p{L}\p{L}\p{L}-\d\d-\d\d`) // регулярное выражение для поиска группы
	if find_group {
		return // если группа уже найдена, то выходим из функции
	}
	count_free := 0                         // счетчик пустых строк
	for columnIndex := 1; ; columnIndex++ { // проходим по столбцам таблицы
		group, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex)+strconv.Itoa(rowIndex))      // получаем значение ячейки
		if status, _ := regexp.MatchString(`\p{L}\p{L}\p{L}\p{L}-\d\d-\d\d`, clear_str(group)); status { // если в ячейке найдена группа
			group = re.FindString(group)          // получаем название группы из строки
			count_day := -1                       // счетчик дней недели
			count_free = 0                        // счетчик пустых строк сбрасываем до нуля
			for row := rowIndex + 2; ; row += 2 { // проходим по строкам таблицы
				number, _ := xlsxFile.GetCellValue(sheetName, cell_name(2)+strconv.Itoa(row)) // получаем номер пары
				if strings.TrimSpace(number) == "1" {                                         // если номер пары равен 1, то увеличиваем счетчик дней на 1
					count_day += 1
				}
				_, err := strconv.Atoi(number) // преобразуем номер пары в число
				if err == nil {                // если преобразование прошло успешно
					// вызываем функцию чтения данных из ячейки с номером пары
					read_data(xlsxFile, sheetName, columnIndex, number, count_day, group, row, row, data)
					read_data(xlsxFile, sheetName, columnIndex, number, count_day, group, row+1, row, data)
				} else { // если номер пары не является числом
					break // прерываем цикл
				}
			}

			count_free = 0   // сбрасываем счетчик пустых строк до нуля
			if !find_group { // если группа еще не была найдена
				find_group = true // устанавливаем флаг в true
			}
		} else if group == "" { // если в ячейке пусто
			count_free += 1 // увеличиваем счетчик пустых строк на 1
		}
		if count_free >= 12 { // если количество пустых строк больше или равно 12
			break // прерываем цикл
		}
	}
}

/// read_data - функция которая считывает данные из excel и отправляет их в СУБД ///
func read_data(xlsxFile *excelize.File, sheetName string, columnIndex int, number string, count_day int,
	group string, elem int, row int, data Links) {
	subject, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex)+strconv.Itoa(elem))
	type_of_subject, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex+1)+strconv.Itoa(elem))
	lecturer, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex+2)+strconv.Itoa(elem))
	auditorium, _ := xlsxFile.GetCellValue(sheetName, cell_name(columnIndex+3)+strconv.Itoa(elem))
	pgsql.MainFunc(group, count_day+1, elem-row+1, clear_str(number), clear_str(subject),
		clear_str(lecturer), clear_str(auditorium), clear_str(type_of_subject), data.institute, data.course)
}

/// clear_str - функция для очистки строки от лишних пробелов ///
func clear_str(str string) string {
	tmp := strings.TrimSpace(strings.ReplaceAll(str, "\n", " "))
	return strings.ReplaceAll(tmp, "\t", " ")
}

/// cell_name - функция, которая считывает строку из ячейки excel ///
func cell_name(number int) string {
	name, _ := excelize.ColumnNumberToName(number)
	return name
}

/// move_file - функция которая перемещает файл из папки temp в папку last ///
func move_file(file string) error {
	os.Remove("last\\" + file)                 // удаление старой версии файла из папки last
	inputFile, err := os.Open("temp\\" + file) //чтение данных из файла в папке temp
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create("last\\" + file) //создание нового файла в папке last
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	_, err = io.Copy(outputFile, inputFile) // запись данных в файл в папке last
	defer inputFile.Close()
	defer outputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	return nil
}
