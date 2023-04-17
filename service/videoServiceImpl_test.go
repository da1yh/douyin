package service

import (
	"douyin/config"
	"douyin/dao"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit2(t *testing.T) {
	dao.InitDb()
}

func TestVideoServiceImpl_AddVideo(t *testing.T) {
	video := dao.Video{
		UserId:      16,
		PublishTime: time.Now(),
		PlayUrl:     config.PlayPathPrefix + "2_bear.mp4",
		CoverUrl:    config.CoverPathPrefix + "2_bear.jpg",
		Title:       "bear",
	}
	vsi := VideoServiceImpl{}
	err := vsi.AddVideo(video)
	assert.Nil(t, err)
}
