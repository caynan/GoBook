package tempconv

import "fmt"

// Celsius a type to represent Celsius temperatures
type Celsius float64

// Fahrenheit a type to represent Fahrenheit temperatures
type Fahrenheit float64

const (
	// AbsoluteZeroC is the lowest possible value in Celsius
	AbsoluteZeroC Celsius = -273.15
	// FreezingC is the temperature on which water freezes in Celsius
	FreezingC Celsius = 0
	// BoilingC is the temperature on which water boils in Celsius
	BoilingC Celsius = 100
)

// CToF converts temperatures from Celsius to Fahrenheit
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts temperatures from Fahrenheit to Celsius
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }
