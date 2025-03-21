package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
	err1       = errors.New("parseData := strings.Split(data, ',') возвращает менее 2 строк")
	err2       = errors.New("строка имеет недопустимый символ ':'")
	err3       = errors.New("значение до разделителя ',' отсутствуют")
)

// parsePackage парсит входящие строки по ","
func parsePackage(data string) (int, time.Duration, error) {
	parseData := strings.Split(data, ",")
	//fmt.Println(len(parseData))
	separator := strings.Contains(data, ":")
	if separator { // separator != 0
		return 0, 0, fmt.Errorf("ошибка выполнения parsePackage: %v", err2)
	}
	if parseData[0] == "" {
		return 0, 0, fmt.Errorf("ошибка выполнения parsePackage: %v", err3)
	}
	if len(parseData) != 2 {
		return 0, 0, fmt.Errorf("ошибка выполнения parsePackage: %v", err1)
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, errors.New("кол-во шагов <= 0")
	}
	timeWorkout, err := time.ParseDuration(parseData[1])
	//fmt.Println(len(parseData))
	//fmt.Println(timeWorkout)
	if err != nil {
		return 0, 0, err
	}
	return steps, timeWorkout, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {

	steps, timeWorkout, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		fmt.Println(errors.New("функция parsePackage вернула 0 шагов"))
		return ""
	}

	distance := StepLength * float64(steps) / 1000
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, timeWorkout)
	info := fmt.Sprintf("\nКоличество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)

	return info
}
