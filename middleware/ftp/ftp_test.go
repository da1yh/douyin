package ftp

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFtp(t *testing.T) {
	InitFtp()
	InitSSH()
	data, _ := os.Open("../../public/2_bear.mp4")
	err := UploadVideo("2_bear.mp4", data)
	assert.True(t, err == nil, "1")
	//assert.True(t, client != nil, "33333")
	err = Screenshot("2_bear.mp4", "2_bear.jpg")
	assert.True(t, err == nil, "2")

	time.Sleep(90 * time.Second)

	data, _ = os.Open("../../public/3_bear.mp4")
	err = UploadVideo("3_bear.mp4", data)
	assert.True(t, err == nil, "3")
	err = Screenshot("3_bear.mp4", "3_bear.jpg")
	assert.True(t, err == nil, "4")
}
