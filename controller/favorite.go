package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Favorite(c *gin.Context) {
	token := c.Query("token")
	if _, exist := userLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user not exist",
		})
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
