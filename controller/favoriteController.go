package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction 通过鉴权后，获得用户id，获得点赞类型和对象视频，通过点赞类型判断是点赞还是取消赞，使用service进行后续操作
func FavoriteAction(c *gin.Context) {
	//TODO favorite logic
	//token := c.Query("token")
	//if _, exist := userLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 0,
	//	})
	//} else {
	//	c.JSON(http.StatusOK, Response{
	//		StatusCode: 1,
	//		StatusMsg:  "user not exist",
	//	})
	//}
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
