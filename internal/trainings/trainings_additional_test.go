package trainings

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
} = Training{}

func TestTrainingParseReuse(t *testing.T) {
	personal := personaldata.Personal{Name: "Анна", Weight: 75, Height: 1.75}
	training := Training{Personal: personal}
	if err := training.Parse("6000,Бег,1h"); err != nil {
		t.Fatal(err)
	}
	before := training
	for _, input := range []string{
		"", "1000,Бег", "1000,Бег,1h,", ",Бег,1h", "1.5,Бег,1h",
		"999999999999999999999999999999,Бег,1h", " 1000,Бег,1h", "1000 ,Бег,1h",
		"0,Ходьба,1h", "-1000,Ходьба,1h", "1000,Ходьба,", "1000,Ходьба,bad",
		"1000,Ходьба,0s", "1000,Ходьба,-1h", "1000,Ходьба,999999999999h",
		"1000,Ходьба, 1h", "1000,Ходьба,1h ",
	} {
		t.Run(input, func(t *testing.T) {
			if err := training.Parse(input); err == nil {
				t.Fatalf("Parse(%q) succeeded, want error", input)
			}
			if training != before {
				t.Fatalf("Parse(%q) changed state: got %+v, want %+v", input, training, before)
			}
		})
	}
	if err := training.Parse("+1200,Ходьба,30.5m"); err != nil {
		t.Fatal(err)
	}
	want := Training{Steps: 1200, TrainingType: "Ходьба", Duration: 30*time.Minute + 30*time.Second, Personal: personal}
	if training != want {
		t.Fatalf("Parse() = %+v, want %+v", training, want)
	}
}

func TestTrainingParsePreservesConversionErrors(t *testing.T) {
	var training Training
	err := training.Parse("bad,Бег,1h")
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) || !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("Parse() error = %v, want wrapped strconv syntax error", err)
	}
	err = training.Parse("1000,Бег,bad")
	if err == nil || errors.Unwrap(err) == nil {
		t.Fatalf("Parse() error = %v, want wrapped duration error", err)
	}
}

func TestTrainingActionInfoInvalidData(t *testing.T) {
	for _, trainingType := range []string{"Бег", "Ходьба"} {
		for _, field := range []string{"steps", "weight", "height", "duration"} {
			for _, value := range []int{0, -1} {
				t.Run(trainingType+"/"+field+"/"+strconv.Itoa(value), func(t *testing.T) {
					training := Training{
						Steps: 6000, TrainingType: trainingType, Duration: time.Hour,
						Personal: personaldata.Personal{Name: "Анна", Weight: 75, Height: 1.75},
					}
					wantErr := ""
					switch field {
					case "steps":
						training.Steps = value
						wantErr = "количество шагов должно быть положительным"
					case "weight":
						training.Weight = float64(value)
						wantErr = "вес должен быть положительным"
					case "height":
						training.Height = float64(value)
						wantErr = "рост должен быть положительным"
					case "duration":
						training.Duration = time.Duration(value)
						wantErr = "продолжительность должна быть положительной"
					}
					before := training
					got, err := training.ActionInfo()
					if got != "" || err == nil || err.Error() != wantErr {
						t.Fatalf("ActionInfo() = (%q, %v), want empty report and %q", got, err, wantErr)
					}
					if training != before {
						t.Fatalf("ActionInfo() changed state: got %+v, want %+v", training, before)
					}
				})
			}
		}
	}
}

func TestTrainingUnknownTypeIsRejectedByActionInfo(t *testing.T) {
	for _, trainingType := range []string{"Плавание", "", "бег", " Бег"} {
		t.Run(trainingType, func(t *testing.T) {
			training := Training{Personal: personaldata.Personal{Weight: 75, Height: 1.75}}
			if err := training.Parse("1000," + trainingType + ",1h"); err != nil {
				t.Fatalf("Parse() rejected unknown type: %v", err)
			}
			if training.TrainingType != trainingType {
				t.Fatalf("Parse() type = %q, want %q", training.TrainingType, trainingType)
			}
			got, err := training.ActionInfo()
			if got != "" || err == nil || err.Error() != "неизвестный тип тренировки" {
				t.Fatalf("ActionInfo() = (%q, %v), want unknown type error", got, err)
			}
		})
	}
}

func TestTrainingActionInfoFractionalDuration(t *testing.T) {
	for _, tt := range []struct {
		trainingType string
		want         string
	}{
		{"Бег", "Тип тренировки: Бег\nДлительность: 0.75 ч.\nДистанция: 4.32 км.\nСкорость: 5.76 км/ч\nСожгли калорий: 345.60\n"},
		{"Ходьба", "Тип тренировки: Ходьба\nДлительность: 0.75 ч.\nДистанция: 4.32 км.\nСкорость: 5.76 км/ч\nСожгли калорий: 172.80\n"},
	} {
		t.Run(tt.trainingType, func(t *testing.T) {
			training := Training{
				Steps: 6000, TrainingType: tt.trainingType, Duration: 45 * time.Minute,
				Personal: personaldata.Personal{Name: "Анна", Weight: 80, Height: 1.6},
			}
			before := training
			got, err := training.ActionInfo()
			if err != nil || got != tt.want {
				t.Fatalf("ActionInfo() = (%q, %v), want (%q, nil)", got, err, tt.want)
			}
			if training != before {
				t.Fatalf("ActionInfo() changed state: got %+v, want %+v", training, before)
			}
		})
	}
}
