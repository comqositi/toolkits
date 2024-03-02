package utils

import (
	"math/rand"
	"time"
)

func createId() int64 {
	id := rand.New(rand.NewSource(time.Now().UnixNano() + (int64(rand.Intn(99999))) + (int64(rand.Intn(99999)))))
	return id.Int63()
}
