package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type PublishListResp struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

func Publish(c *gin.Context) {
	//check if user exists
	token := c.PostForm("token")
	if _, exist := userLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user not exist",
		})
		return
	}
	//get the file and save in the corresponding path
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fileName := filepath.Base(data.Filename)
	user := userLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, fileName)
	saveFilePath := filepath.Join("./public", finalName)
	err = c.SaveUploadedFile(data, saveFilePath)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fileName + " uploaded successfully!",
	})
}

// PublishList all users have same publish_list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, PublishListResp{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideoList,
	})
}
