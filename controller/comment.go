package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentActionResp struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

type CommentListResp struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	if user, exist := userLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResp{
				Response: Response{
					StatusCode: 0,
				},
				Comment: Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "2020-01-02",
				},
			})
		} else {
			c.JSON(http.StatusOK, CommentActionResp{
				Response: Response{
					StatusCode: 0,
				},
			})
		}
	} else {
		c.JSON(http.StatusOK, CommentActionResp{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "user not exist",
			},
		})
	}
}

func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResp{
		Response: Response{
			StatusCode: 0,
		},
		CommentList: DemoComment,
	})
}
