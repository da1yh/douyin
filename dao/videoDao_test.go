package dao

import (
	"douyin/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit2(t *testing.T) {
	InitDb()
}

func TestAddVideo(t *testing.T) {
	video := Video{
		UserId:      16,
		PublishTime: time.Now(),
		PlayUrl:     config.PlayPathPrefix + "2_bear.mp4",
		CoverUrl:    config.CoverPathPrefix + "2_bear.jpg",
		Title:       "bear",
	}
	err := AddVideo(video)
	assert.Nil(t, err)
}
