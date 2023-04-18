package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInit2(t *testing.T) {
	InitDb()
}

//func TestAddVideo(t *testing.T) {
//	video := Video{
//		UserId:      16,
//		PublishTime: time.Now(),
//		PlayUrl:     config.PlayPathPrefix + "2_bear.mp4",
//		CoverUrl:    config.CoverPathPrefix + "2_bear.jpg",
//		Title:       "bear",
//	}
//	err := AddVideo(video)
//	assert.Nil(t, err)
//}

//func TestFindPublishedVideosByUserId(t *testing.T) {
//	videos, err := FindPublishedVideosByUserId(16)
//	assert.Nil(t, err)
//	assert.Equal(t, len(videos), 2)
//	for _, v := range videos {
//		assert.Equal(t, v.UserId, int64(16))
//	}
//
//	videos, err = FindPublishedVideosByUserId(20)
//	assert.Nil(t, err)
//	assert.Equal(t, len(videos), 0)
//
//}

func TestFindVideosByTimeAndNum(t *testing.T) {
	videos, err := FindVideosByTimeAndNum(time.Now(), 20)
	assert.Nil(t, err)
	assert.Equal(t, len(videos), 3)
}
