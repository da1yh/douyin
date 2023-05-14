package util

import (
	"math/rand"
	"time"
)

// RandomDuration 生成1天到3天之间的随机时间段
func RandomDuration() time.Duration {
	rand.Seed(time.Now().Unix())
	duration := rand.Intn(48) + 24
	interval := time.Duration(duration) * time.Hour
	return interval
}
