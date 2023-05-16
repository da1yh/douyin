package controller

import (
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

// FavoriteAction 通过鉴权后，获得用户id，获得点赞类型和对象视频，通过点赞类型判断是点赞还是取消赞，使用service进行后续操作
func FavoriteAction(c *gin.Context) {
	fromUserId, toVideoIdStr, tp := c.GetInt64("id"), c.Query("video_id"), c.Query("action_type")
	toVideoId, _ := strconv.ParseInt(toVideoIdStr, 10, 64)
	fsi := service.FavoriteServiceImpl{}
	if tp == "1" { //这里的1而不是RedisLikeType是因为1是接口决定的
		err := fsi.AddFavoriteByBothId(fromUserId, toVideoId)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1, StatusMsg: "failed to thumb up",
			})
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0, StatusMsg: "thumb up successfully",
			})
		}
	} else {
		err := fsi.DeleteFavoriteByBothId(fromUserId, toVideoId)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1, StatusMsg: "failed to cancel thumb up",
			})
		} else {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0, StatusMsg: "cancel thumb up successfully",
			})
		}
	}
}

// FavoriteList 通过鉴权后，使用service获得点赞视频列表
func FavoriteList(c *gin.Context) {
	myId, userIdStr := c.GetInt64("id"), c.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	fsi := service.FavoriteServiceImpl{}
	videoIds, err := fsi.FindFavoriteVideoIdsByFromUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to get user's favorite videos",
		})
	}

	videoList := GetVideoListByIds(videoIds, myId)
	userResp, err := GetUserRespByBothId(myId, userId)

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to get user info",
		})
	}

	for i := range videoList {
		videoList[i].Author = userResp
		//video.Author = userResp	// 副本
		//(&video).Author = userResp
	}

	c.JSON(http.StatusOK, VideoListResp{
		Response: Response{
			StatusCode: 0, StatusMsg: "get user's favorite video list successfully",
		},
		VideoList: videoList,
	})
}

// GetVideoListByIds 通过video的id列表，获得Video列表的部分信息
func GetVideoListByIds(ids []int64, myId int64) []Video {
	videoList := make([]Video, 0)
	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, id := range ids {
		go GetVideoById(id, &videoList, &wg, myId)
	}
	wg.Wait()
	return videoList
}

// GetVideoById 通过videoId，返回Video的部分信息
func GetVideoById(id int64, videoList *[]Video, wg *sync.WaitGroup, myId int64) {
	defer wg.Done()
	vsi := service.VideoServiceImpl{}
	videoDao, _ := vsi.FindVideoById(id)
	video := Video{}
	video.Id = id
	video.PlayUrl = videoDao.PlayUrl
	video.CoverUrl = videoDao.CoverUrl
	video.Title = videoDao.Title

	var wg2 sync.WaitGroup
	wg2.Add(3)

	go func(video *Video, id int64) {
		defer wg2.Done()
		fsi := service.FavoriteServiceImpl{}
		cnt, _ := fsi.CountFavoritesByToVideoId(id)
		video.FavoriteCount = cnt
	}(&video, id)

	go func(video *Video, id int64) {
		defer wg2.Done()
		csi := service.CommentServiceImpl{}
		cnt, _ := csi.CountCommentsByToVideoId(id)
		video.CommentCount = cnt
	}(&video, id)

	go func(video *Video, id, myId int64) {
		defer wg2.Done()
		fsi := service.FavoriteServiceImpl{}
		ok, _ := fsi.CheckFavoriteByBothId(myId, id)
		video.IsFavorite = ok
	}(&video, id, myId)

	wg2.Wait()
	*videoList = append(*videoList, video)
}
