package controller

import (
	"douyin/config"
	"douyin/dao"
	"douyin/middleware/ftp"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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

func GetRespUserByBothId(myId, yourId int64) (User, error) {
	rsi := service.RelationServiceImpl{}
	usi := service.UserServiceImpl{}
	followCount, _ := rsi.CountRelationsByFromUserId(yourId)
	followerCount, _ := rsi.CountRelationsByToUserId(yourId)
	isFollow, _ := rsi.CheckRelationByBothId(myId, yourId)
	usr, _ := usi.FindUserById(yourId)
	user := User{
		Id:            yourId,
		Name:          usr.Name,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}
	return user, nil
}

// PublishList 根据id，查询数据库，获得video
func PublishList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	author, err := GetRespUserByBothId(c.GetInt64("id"), userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "user not exist",
		})
		return
	}
	vsi := service.VideoServiceImpl{}
	videos, err := vsi.FindPublishedVideosByUserId(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2, StatusMsg: "fetch user's published video(s) failed",
		})
		return
	}

	videoList := make([]Video, 0)

	fsi := service.FavoriteServiceImpl{}
	csi := service.CommentServiceImpl{}

	for _, video := range videos {
		var tmpVideo Video
		tmpVideo.Id = video.Id
		tmpVideo.Author = author
		tmpVideo.PlayUrl = video.PlayUrl
		tmpVideo.CoverUrl = video.CoverUrl
		tmpVideo.Title = video.Title
		tmpVideo.FavoriteCount, _ = fsi.CountFavoritesByToVideoId(video.Id)
		tmpVideo.IsFavorite, _ = fsi.CheckFavoriteByBothId(author.Id, video.Id)
		tmpVideo.CommentCount, _ = csi.CountCommentsByToVideoId(video.Id)
		videoList = append(videoList, tmpVideo)
	}

	c.JSON(http.StatusOK, VideoListResp{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
