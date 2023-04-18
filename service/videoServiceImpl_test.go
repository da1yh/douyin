package service

import (
	"douyin/dao"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit2(t *testing.T) {
	dao.InitDb()
}

//func TestVideoServiceImpl_AddVideo(t *testing.T) {
//	video := dao.Video{
//		UserId:      16,
//		PublishTime: time.Now(),
//		PlayUrl:     config.PlayPathPrefix + "2_bear.mp4",
//		CoverUrl:    config.CoverPathPrefix + "2_bear.jpg",
//		Title:       "bear",
//	}
//	vsi := VideoServiceImpl{}
//	err := vsi.AddVideo(video)
//	assert.Nil(t, err)
//}

//func TestVideoServiceImpl_FindPublishedVideosByUserId(t *testing.T) {
//	vsi := VideoServiceImpl{}
//	videos, err := vsi.FindPublishedVideosByUserId(16)
//	assert.Nil(t, err)
//	assert.Equal(t, len(videos), 2)
//	for _, v := range videos {
//		assert.Equal(t, v.UserId, int64(16))
//	}
//	videos, err = vsi.FindPublishedVideosByUserId(20)
//	assert.Nil(t, err)
//	assert.Equal(t, len(videos), 0)
//}

func TestVideoServiceImpl_FindVideosByTimeAndNum(t *testing.T) {
	vsi := VideoServiceImpl{}
	videos, err := vsi.FindVideosByTimeAndNum(time.Now(), 20)
	assert.Nil(t, err)
	assert.Equal(t, len(videos), 3)
}
