// Package spentenergy рассчитывает дистанцию, среднюю скорость и расход калорий.
package spentenergy

import (
	"errors"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// WalkingSpentCalories возвращает расход энергии при ходьбе в ккал.
// Вес задаётся в кг, рост в метрах. Все параметры должны быть положительными;
// иначе функция возвращает ноль и ошибку.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

// RunningSpentCalories возвращает расход энергии при беге в ккал.
// Вес задаётся в кг, рост в метрах. Все параметры должны быть положительными;
// иначе функция возвращает ноль и ошибку.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}
	return weight * MeanSpeed(steps, height, duration) * duration.Minutes() / minInH, nil
}

// MeanSpeed возвращает среднюю скорость в км/ч для роста в метрах.
// При отрицательных шагах или неположительной длительности возвращает ноль.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps < 0 || duration <= 0 {
		return 0
	}
	return Distance(steps, height) / duration.Hours()
}

// Distance возвращает дистанцию в километрах для роста в метрах.
// Входные значения не проверяются; ноль шагов даёт нулевую дистанцию.
func Distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient / mInKm
}
