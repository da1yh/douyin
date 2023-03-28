package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync/atomic"
)

type ChatResp struct {
	Response
	MessageList []Message `json:"message_list"`
}

var messageIdSequence int64 = 1

var messageList = make(map[string][]Message) // key is chatKey which represent who are chatting, value is messages

func MessageAction(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	if user, exist := userLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))
		atomic.AddInt64(&messageIdSequence, 1)
		curMessage := Message{
			Id:         messageIdSequence,
			FromUserId: user.Id,
			ToUserId:   int64(userIdB),
			Content:    content,
		}
		if messages, exist := messageList[chatKey]; exist {
			messageList[chatKey] = append(messages, curMessage)
		} else {
			messageList[chatKey] = []Message{curMessage}
		}
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

func MessageChat(c *gin.Context) {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	if user, exist := userLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))
		c.JSON(http.StatusOK, ChatResp{
			Response:    Response{StatusCode: 0},
			MessageList: messageList[chatKey],
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user not exist",
		})
	}
}

func genChatKey(a, b int64) string {
	if a < b {
		return fmt.Sprintf("%d+%d", a, b)
	}
	return fmt.Sprintf("%d+%d", b, a)
}
