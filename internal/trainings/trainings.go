// Package trainings разбирает записи тренировок и формирует отчёты о них.
package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Training хранит параметры тренировки и персональные данные пользователя.
type Training struct {
	// Steps хранит количество шагов за тренировку.
	Steps int
	// TrainingType хранит тип тренировки: «Бег» или «Ходьба» для расчёта калорий.
	TrainingType string
	// Duration хранит длительность тренировки.
	Duration time.Duration
	// Personal содержит имя, вес в кг и рост в метрах; предоставляет метод Print.
	personaldata.Personal
}

// Parse разбирает строку «шаги,тип,длительность», например «3456,Ходьба,3h00m».
// Возвращает ошибку при числе полей, отличном от трёх, ошибке преобразования
// или неположительных шагах и длительности. Пробелы не удаляются.
// При ошибке объект не меняется; при успехе сохраняются персональные данные.
// Тип тренировки проверяется в ActionInfo.
func (t *Training) Parse(datastring string) (err error) {
	fields := strings.Split(datastring, ",")
	if len(fields) != 3 {
		return errors.New("ожидаются три поля тренировки")
	}
	steps, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("количество шагов: %w", err)
	}
	if steps <= 0 {
		return errors.New("количество шагов должно быть положительным")
	}
	duration, err := time.ParseDuration(fields[2])
	if err != nil {
		return fmt.Errorf("продолжительность: %w", err)
	}
	if duration <= 0 {
		return errors.New("продолжительность должна быть положительной")
	}

	t.Steps = steps
	t.TrainingType = fields[1]
	t.Duration = duration
	return nil
}

// ActionInfo возвращает отчёт с длительностью в часах, дистанцией в км,
// скоростью в км/ч и расходом энергии в ккал. Отчёт заканчивается переводом строки.
// При неизвестном типе тренировки или ошибке расчёта возвращает пустую строку
// и ошибку. Для расчёта шаги, вес, рост и длительность должны быть положительными.
func (t Training) ActionInfo() (string, error) {
	var calories float64
	var err error
	switch t.TrainingType {
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}
	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories), nil
}
