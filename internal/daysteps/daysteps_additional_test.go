package daysteps

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
)

var _ interface {
	Print()
	ActionInfo() (string, error)
} = DaySteps{}

func TestDayStepsParseReuse(t *testing.T) {
	personal := personaldata.Personal{Name: "Анна", Weight: 75, Height: 1.75}
	ds := DaySteps{Personal: personal}
	if err := ds.Parse("6000,1h"); err != nil {
		t.Fatal(err)
	}
	before := ds
	for _, input := range []string{
		"", "1000", "1000,1h,", ",1h", "1.5,1h",
		"999999999999999999999999999999,1h", " 1000,1h", "1000 ,1h",
		"0,1h", "-1000,1h", "1000,", "1000,bad", "1000,0s", "1000,-1h",
		"1000,999999999999h", "1000, 1h", "1000,1h ",
	} {
		t.Run(input, func(t *testing.T) {
			if err := ds.Parse(input); err == nil {
				t.Fatalf("Parse(%q) succeeded, want error", input)
			}
			if ds != before {
				t.Fatalf("Parse(%q) changed state: got %+v, want %+v", input, ds, before)
			}
		})
	}
	if err := ds.Parse("+1200,30.5m"); err != nil {
		t.Fatal(err)
	}
	want := DaySteps{Steps: 1200, Duration: 30*time.Minute + 30*time.Second, Personal: personal}
	if ds != want {
		t.Fatalf("Parse() = %+v, want %+v", ds, want)
	}
}

func TestDayStepsParsePreservesConversionErrors(t *testing.T) {
	var ds DaySteps
	err := ds.Parse("bad,1h")
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) || !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("Parse() error = %v, want wrapped strconv syntax error", err)
	}
	err = ds.Parse("1000,bad")
	if err == nil || errors.Unwrap(err) == nil {
		t.Fatalf("Parse() error = %v, want wrapped duration error", err)
	}
}

func TestDayStepsActionInfoFractionalDuration(t *testing.T) {
	ds := DaySteps{
		Steps: 6000, Duration: 45 * time.Minute,
		Personal: personaldata.Personal{Name: "Анна", Weight: 80, Height: 1.6},
	}
	before := ds
	const want = "Количество шагов: 6000.\nДистанция составила 4.32 км.\nВы сожгли 172.80 ккал.\n"
	got, err := ds.ActionInfo()
	if err != nil || got != want {
		t.Fatalf("ActionInfo() = (%q, %v), want (%q, nil)", got, err, want)
	}
	if ds != before {
		t.Fatalf("ActionInfo() changed state: got %+v, want %+v", ds, before)
	}
}
