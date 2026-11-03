// Package actioninfo обрабатывает записи активности и выводит отчёты через общий интерфейс.
package actioninfo

import (
	"fmt"
	"log"
)

// DataParser объединяет разбор записи активности и подготовку отчёта.
// Интерфейс реализуют указатели *trainings.Training и *daysteps.DaySteps.
type DataParser interface {
	// Parse разбирает запись и сохраняет данные для последующего вызова ActionInfo.
	Parse(string) error
	// ActionInfo возвращает отчёт о сохранённой активности или ошибку расчёта.
	ActionInfo() (string, error)
}

// Info последовательно разбирает записи через dp и печатает успешные отчёты.
// Каждую ошибку Parse или ActionInfo логирует и переходит к следующей записи.
// После ошибки Parse отчёт не запрашивается; после ошибки ActionInfo не выводится.
// К успешному отчёту символы не добавляются. Пустой набор ничего не выводит.
func Info(dataset []string, dp DataParser) {
	for _, data := range dataset {
		if err := dp.Parse(data); err != nil {
			log.Println(err)
			continue
		}
		report, err := dp.ActionInfo()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Print(report)
	}
}
