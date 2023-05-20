package controller

import (
	"douyin/service"
	"douyin/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CommentActionResp struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

type CommentListResp struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

// CommentAction 通过鉴权后，获得相关信息，进行操作
func CommentAction(c *gin.Context) {
	fromUserId := c.GetInt64("id")
	toVideoIdStr := c.PostForm("video_id")
	if len(toVideoIdStr) == 0 {
		toVideoIdStr = c.Query("video_id")
	}
	actionTypeStr := c.PostForm("action_type")
	if len(actionTypeStr) == 0 {
		actionTypeStr = c.Query("action_type")
	}
	toVideoId, _ := strconv.ParseInt(toVideoIdStr, 10, 64)
	nowTime := time.Now()
	userResp, err := GetUserRespByBothId(fromUserId, fromUserId)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResp{
			Response: Response{
				StatusCode: 2, StatusMsg: "failed to get user info of comment",
			},
		})
	}
	csi := service.CommentServiceImpl{}
	// 发表评论
	if actionTypeStr == "1" {
		commentText := c.PostForm("comment_text")
		if len(commentText) == 0 {
			commentText = c.Query("comment_text")
		}
		commentId, err := csi.AddCommentByAll(fromUserId, toVideoId, commentText, nowTime)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResp{
				Response: Response{
					StatusCode: 1, StatusMsg: "failed to add comment",
				},
			})
		}
		c.JSON(http.StatusOK, CommentActionResp{
			Response: Response{
				StatusCode: 0, StatusMsg: "comment successfully",
			},
			Comment: Comment{
				Id:         commentId,
				User:       userResp,
				Content:    commentText,
				CreateDate: nowTime.Format(util.TimeFormat),
			},
		})
	} else { // 删除评论
		commentIdStr := c.PostForm("comment_id")
		if len(commentIdStr) == 0 {
			commentIdStr = c.Query("comment_id")
		}
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 64)
		err = csi.DeleteCommentById(commentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResp{
				Response: Response{
					StatusCode: 3, StatusMsg: "failed to delete comment",
				},
			})
		}
		c.JSON(http.StatusOK, CommentActionResp{
			Response: Response{
				StatusCode: 0, StatusMsg: "delete comment successfully",
			},
		})
	}
}

func CommentList(c *gin.Context) {
	myId := c.GetInt64("id")
	toVideoIdStr := c.Query("video_id")
	toVideo, _ := strconv.ParseInt(toVideoIdStr, 10, 64)
	csi := service.CommentServiceImpl{}
	ids, err := csi.FindCommentIdsByToVideoId(toVideo)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to get comments' id of the video",
		})
	}
	commentList := GetCommentListByIds(ids, myId)
	c.JSON(http.StatusOK, CommentListResp{
		Response: Response{
			StatusCode: 0, StatusMsg: "get comments of the video successfully",
		},
		CommentList: commentList,
	})
}

func GetCommentListByIds(ids []int64, myId int64) []Comment {
	commentList := make([]Comment, 0)
	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, id := range ids {
		go GetCommentById(id, &commentList, &wg, myId)
	}
	wg.Wait()
	return commentList
}

func GetCommentById(id int64, commentList *[]Comment, wg *sync.WaitGroup, myId int64) {
	defer wg.Done()
	csi := service.CommentServiceImpl{}
	comment := Comment{}
	comment.Id = id
	commentDao, err := csi.FindCommentById(id)
	if err != nil {
		log.Println("failed to get comment info by id", err)
		return
	}
	comment.Content = commentDao.Content
	comment.CreateDate = commentDao.CreateDate.Format(util.TimeFormat)
	userId := commentDao.FromUserId
	userResp, err := GetUserRespByBothId(myId, userId)
	if err != nil {
		log.Println("failed to get user response", err)
		return
	}
	comment.User = userResp
	*commentList = append(*commentList, comment)
}
