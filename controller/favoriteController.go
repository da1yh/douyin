package controller

import (
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; exist {
		c.JSON(http.StatusOK, VideoListResp{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: DemoVideoList,
		})
	}
}
