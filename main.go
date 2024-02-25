package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000  // количество метров в одном километре
	MinInHours = 60    // количество минут в одном часе
	LenStep    = 0.65  // длина одного шага
	CmInM      = 100   // количество сантиметров в одном метре
	KmHInMsec  = 0.278 // коэффициент для перевода км/ч в м/с
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // тип тренировки
	Action       int           // количество повторов(шаги, гребки при плавании)
	LenStep      float64       // длина одного шага или гребка в м
	Duration     time.Duration // продолжительность тренировки
	Weight       float64       // вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость движения.
func (t Training) meanSpeed() float64 {
	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// TrainingInfo возвращает структуру InfoMessage с информацией о тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потраченных килокалорий
const (
	CaloriesMeanSpeedMultiplier      = 18    // множитель средней скорости бега
	CaloriesMeanSpeedShift           = 1.79  // коэффициент изменения средней скорости
	CaloriesWeightMultiplier         = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier    = 0.029 // коэффициент для роста
	SwimmingLenStep                  = 1.38  // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1   // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2     // множитель веса пользователя
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возвращает количество потраченных калорий при беге.
func (r Running) Calories() float64 {
	speed := r.meanSpeed()
	return ((CaloriesMeanSpeedMultiplier*speed + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours)
}

// Walking структура, описывающая тренировку Ходьба.
type Walking struct {
	Training
	Height float64
}

// Calories возвращает количество потраченных калорий при ходьбе.
func (w Walking) Calories() float64 {
	speedMPerSec := w.meanSpeed() * KmHInMsec
	return ((CaloriesWeightMultiplier*w.Weight + (math.Pow(speedMPerSec, 2)/w.Height)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours)
}

// Swimming структура, описывающая тренировку Плавание.
type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

// meanSpeed переопределение метода для плавания.
func (s Swimming) meanSpeed() float64 {
	return float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Hours()
}

// Calories возвращает количество калорий, потраченных при плавании.
func (s Swimming) Calories() float64 {
	speed := s.meanSpeed()
	return ((speed + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours())
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	return fmt.Sprint(training.TrainingInfo())
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
		CountPool:  40,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       10000,
			LenStep:      LenStep,
			Duration:     2 * time.Hour,
			Weight:       70,
		},
		Height: 175,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       3000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       70,
		},
	}

	fmt.Println(ReadData(running))
}
