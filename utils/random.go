package utils

import (
	"math/rand"
	"time"
)

var (
	rander = rand.New(rand.NewSource(time.Now().Unix()))
)

//RandomDigit random a number string, length is count
func RandomDigit(count int) string {
	var byts []byte = make([]byte, count)

	for i := 0; i < count; i++ {
		n := rander.Intn(10)
		byts[i] = byte('0' + n)
	}
	return string(byts)
}
