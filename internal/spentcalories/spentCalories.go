package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
)

var (
	// parseData := strings.Split(data, ',') возвращает менее 3 строк
	err1 = errors.New("parseData := strings.Split(data, ',') возвращает менее 3 строк")
	// строка имеет недопустимый символ ':'
	err2 = errors.New("строка имеет недопустимый символ ':'")
	// значение до разделителя ',' отсутствуют
	err3 = errors.New("значение до разделителя ',' отсутствуют")
)

// Функция parseTraining парсит строку на 3 элемента по ","
func parseTraining(data string) (int, string, time.Duration, error) {
	parseData := strings.Split(data, ",")

	separator := strings.Contains(data, ":")
	if separator { // separator != 0
		return 0, "", 0, fmt.Errorf("ошибка выполнения parseTraining: %v", err2)
	}
	if len(parseData) != 3 {
		return 0, "", 0, fmt.Errorf("ошибка выполнения parseTraining: %v", err1)
	}
	if parseData[0] == "" {
		return 0, "", 0, fmt.Errorf("ошибка выполнения parseTraining: %v", err3)
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, "", 0, err
	}
	temeWorkout, err := time.ParseDuration(parseData[2])
	if err != nil {
		return 0, "", 0, err
	}
	return steps, parseData[1] /*вид активности*/, temeWorkout, err
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	distance := float64(steps) * lenStep / mInKm
	return distance
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps)
	averageSpid := distance / duration.Hours()
	return averageSpid
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, typeActivity, duration, err := parseTraining(data)
	if err != nil {
		return ""
	}
	var treningInfo string
	switch typeActivity {
	case "Ходьба":
		treningInfo = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			typeActivity, duration.Hours(), distance(steps), meanSpeed(steps, duration), WalkingSpentCalories(steps, weight, height, duration))
	case "Бег":
		treningInfo = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			typeActivity, duration.Hours(), distance(steps), meanSpeed(steps, duration), RunningSpentCalories(steps, weight, duration))
	default:
		treningInfo = "неизвестный тип тренировки"
	}

	return treningInfo
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	averageSpid := meanSpeed(steps, duration)
	runningCalories := ((runningCaloriesMeanSpeedMultiplier * averageSpid) - runningCaloriesMeanSpeedShift) * weight
	return runningCalories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// ваш код здесь
	averageSpid := meanSpeed(steps, duration)
	walkingCalories := ((walkingCaloriesWeightMultiplier * weight) + (averageSpid*averageSpid/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
	return walkingCalories
}
