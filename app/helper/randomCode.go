package helper

import (
	"math/rand"
)

func Generate4DigitCode() int {
	return rand.Intn(9000) + 1000
}
