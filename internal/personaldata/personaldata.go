// Package personaldata хранит и выводит персональные данные пользователя.
package personaldata

import "fmt"

// Personal содержит имя, вес и рост пользователя.
type Personal struct {
	// Name хранит имя пользователя.
	Name string
	// Weight хранит вес в килограммах.
	Weight float64
	// Height хранит рост в метрах.
	Height float64
}

// Print выводит имя, вес в кг и рост в метрах в стандартный поток вывода.
// Вес и рост печатаются с двумя знаками после точки, каждое поле на своей строке.
func (p Personal) Print() {
	fmt.Printf("Имя: %s\nВес: %.2f\nРост: %.2f\n", p.Name, p.Weight, p.Height)
}
