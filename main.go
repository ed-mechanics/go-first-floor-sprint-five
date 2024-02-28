package main

import (
	"fmt"
	"time"
)

const (
	MInKm                            = 1000
	MinInHours                       = 60
	LenStep                          = 0.65
	CaloriesMeanSpeedMultiplier      = 18
	CaloriesMeanSpeedShift           = 1.79
	CaloriesWeightMultiplier         = 0.035
	CaloriesSpeedHeightMultiplier    = 0.029
	KmHInMsec                        = 0.278
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

func (t Training) meanSpeed() float64 {
	durationInHours := t.Duration.Hours()
	if durationInHours == 0 {
		return 0
	}
	return t.distance() / durationInHours
}

func (t Training) Calories() float64 {
	return 0
}

func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

type Running struct {
	Training
}

func (r Running) Calories() float64 {
	speed := r.meanSpeed()
	return (CaloriesMeanSpeedMultiplier*speed + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours
}

type Walking struct {
	Training
	Height float64
}

func (w Walking) Calories() float64 {
	speedInMetersPerSecond := w.meanSpeed() * KmHInMsec
	if w.Height == 0 {
		return 0
	}
	return (CaloriesWeightMultiplier*w.Weight + (speedInMetersPerSecond*speedInMetersPerSecond/w.Height)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours
}

type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

func (s Swimming) meanSpeed() float64 {
	if s.Duration.Hours() == 0 {
		return 0
	}
	return float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Hours()
}

func (s Swimming) Calories() float64 {
	speed := s.meanSpeed()
	return (speed + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()
	info := training.TrainingInfo()
	info.Calories = calories
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		info.TrainingType,
		info.Duration.Minutes(),
		info.Distance,
		info.Speed,
		info.Calories)
}

func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
