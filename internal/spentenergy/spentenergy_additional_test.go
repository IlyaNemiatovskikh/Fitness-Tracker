package spentenergy

import (
	"math"
	"testing"
	"time"
)

func TestMeanSpeedNegativeSteps(t *testing.T) {
	if got := MeanSpeed(-1000, 1.75, time.Hour); got != 0 {
		t.Fatalf("MeanSpeed() = %v, want 0", got)
	}
}

func TestCaloriesMissingBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		calculate func(int, float64, float64, time.Duration) (float64, error)
		height    float64
		duration  time.Duration
	}{
		{
			name: "running zero height", calculate: RunningSpentCalories,
			height: 0, duration: time.Hour,
		},
		{
			name: "walking zero duration", calculate: WalkingSpentCalories,
			height: 1.75, duration: 0,
		},
		{
			name: "walking negative duration", calculate: WalkingSpentCalories,
			height: 1.75, duration: -time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.calculate(6000, 75, tt.height, tt.duration)
			if err == nil || got != 0 {
				t.Fatalf("calories = %v, error = %v; want 0 and an error", got, err)
			}
		})
	}
}

func TestCalculationsKeepPrecision(t *testing.T) {
	// 1234 шага по 0.7425 м дают 0.916245 км; 37.5 минуты равны 0.625 часа.
	const tolerance = 1e-10 // Допуск на погрешность float64; меньше округления до сотых.
	const duration = 37*time.Minute + 30*time.Second
	for _, tt := range []struct {
		name string
		got  float64
		want float64
	}{
		{name: "distance", got: Distance(1234, 1.65), want: 0.916245},
		{name: "speed", got: MeanSpeed(1234, 1.65, duration), want: 1.465992},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if math.Abs(tt.got-tt.want) > tolerance {
				t.Fatalf("got %.12f, want %.12f", tt.got, tt.want)
			}
		})
	}
	for _, tt := range []struct {
		name      string
		calculate func(int, float64, float64, time.Duration) (float64, error)
		want      float64
	}{
		{name: "running", calculate: RunningSpentCalories, want: 59.96823525},
		{name: "walking", calculate: WalkingSpentCalories, want: 29.984117625},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.calculate(1234, 65.45, 1.65, duration)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(got-tt.want) > tolerance {
				t.Fatalf("got %.12f, want %.12f", got, tt.want)
			}
		})
	}
}
