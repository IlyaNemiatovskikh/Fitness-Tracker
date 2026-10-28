# Модуль фитнес-трекера

Модуль разбирает записи о тренировках и дневных прогулках, рассчитывает
дистанцию, среднюю скорость и расход калорий, затем печатает отчёты.

Модуль: `github.com/Yandex-Practicum/tracker`, Go 1.24.1 или новее.
Использован [шаблон задания](https://github.com/Yandex-Practicum/go1fl-sprint5-final/tree/a4ddf000ac5a70eddb0700c14073d0c9eaf30985).


## Запуск

Из корня проекта:

```bash
make bootstrap
make deps
make run
```

`bootstrap` устанавливает Go 1.24.1 в `.tools/go`.
Нужны Bash, curl, tar и sha256sum либо shasum, поддерживаются Linux и macOS
на amd64/arm64. Повторный запуск использует установленный локальный Go.
Если Go версии не ниже 1.24.1 уже установлен, этот шаг можно пропустить.
`make deps` скачивает и проверяет зависимости, сохраняя `go.mod` и `go.sum`.

Makefile выбирает `.tools/go/bin/go`, если он существует, иначе `go` из `PATH`.
Другой компилятор можно указать через `make GO=/absolute/path/to/go check`.
Для прямых команд добавьте локальный Go в `PATH`:

```bash
export PATH="$PWD/.tools/go/bin:$PATH"
go run ./cmd/tracker
```

В `main` уже заданы демонстрационные записи. Программа печатает четыре отчёта
о прогулках и пять отчётов о тренировках.
