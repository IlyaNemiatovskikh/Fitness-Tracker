// Package daysteps разбирает записи дневных прогулок и формирует отчёты о них.
package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps хранит параметры прогулки и персональные данные пользователя.
type DaySteps struct {
	// Steps хранит количество шагов за прогулку.
	Steps int
	// Duration хранит длительность прогулки.
	Duration time.Duration
	// Personal содержит имя, вес в кг и рост в метрах; предоставляет метод Print.
	personaldata.Personal
}

// Parse разбирает строку «шаги,длительность», например «678,0h50m».
// Возвращает ошибку при числе полей, отличном от двух, ошибке преобразования
// или неположительных шагах и длительности. Пробелы не удаляются.
// При ошибке объект не меняется; при успехе сохраняются персональные данные.
func (ds *DaySteps) Parse(datastring string) (err error) {
	fields := strings.Split(datastring, ",")
	if len(fields) != 2 {
		return errors.New("ожидаются два поля дневной активности")
	}
	steps, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("количество шагов: %w", err)
	}
	if steps <= 0 {
		return errors.New("количество шагов должно быть положительным")
	}
	duration, err := time.ParseDuration(fields[1])
	if err != nil {
		return fmt.Errorf("продолжительность: %w", err)
	}
	if duration <= 0 {
		return errors.New("продолжительность должна быть положительной")
	}

	ds.Steps = steps
	ds.Duration = duration
	return nil
}

// ActionInfo возвращает отчёт с количеством шагов, дистанцией в км
// и расходом энергии при ходьбе в ккал. Отчёт заканчивается переводом строки.
// При неположительных шагах, весе, росте или длительности возвращает
// пустую строку и ошибку.
func (ds DaySteps) ActionInfo() (string, error) {
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories), nil
}
