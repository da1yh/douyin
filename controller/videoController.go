package controller

import (
	"douyin/config"
	"douyin/dao"
	"douyin/middleware/ftp"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type VideoListResp struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
}

// PublishAction 通过鉴权后，获得token的id，name和password，获得context的视频数据和标题，封装一个新的视频文件名，将它上传到ftp服务器，
// 如果成功，利用ssh调用ffmpeg命令执行截图，保存到相应路径下，如果成功，将数据写入video数据库
func PublishAction(c *gin.Context) {
	userId, userName, _ := c.GetInt64("id"), c.GetString("name"), c.GetString("password")
	data, err := c.FormFile("data")
	if err != nil {
		log.Println("error: ", err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "cannot get video you upload",
		})
		return
	}
	title := c.PostForm("title")
	file, err := data.Open()
	if err != nil {
		log.Println("error: ", err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 2, StatusMsg: "open video failed",
		})
		return
	}
	videoFileName := userName + "-" + title + ".mp4"
	imageFileName := userName + "-" + title + ".jpg"
	err = ftp.UploadVideo(videoFileName, file)
	if err != nil {
		log.Println("error: ", err.Error())
		c.JSON(http.StatusOK, Response{
			StatusCode: 3, StatusMsg: "upload video failed",
		})
		return
	}
	err = ftp.Screenshot(videoFileName, imageFileName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 4, StatusMsg: "get cover of video failed",
		})
		return
	}
	newVideo := dao.Video{
		UserId:      userId,
		PublishTime: time.Now(),
		PlayUrl:     config.PlayPathPrefix + videoFileName,
		CoverUrl:    config.CoverPathPrefix + imageFileName,
		Title:       title,
	}
	vsi := service.VideoServiceImpl{}
	if err = vsi.AddVideo(newVideo); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 5, StatusMsg: "fail to upload video",
		})
		return
	}
	log.Println("video upload successfull")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0, StatusMsg: "video upload successfully",
	})
}

// PublishList all users have same publish_list
func PublishList(c *gin.Context) {
	//c.JSON(http.StatusOK, VideoListResp{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideoList,
	//})
}
