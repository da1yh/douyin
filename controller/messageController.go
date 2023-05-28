package controller

import (
	"douyin/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ChatResp struct {
	Response
	MessageList []Message `json:"message_list"`
}

var messageIdSequence int64 = 1

//var messageList = make(map[string][]Message) // key is chatKey which represent who are chatting, value is messages

func MessageAction(c *gin.Context) {
	fromUserId := c.GetInt64("id")
	toUserIdStr := c.PostForm("to_user_id")
	if len(toUserIdStr) == 0 {
		toUserIdStr = c.Query("to_user_id")
	}
	actionType := c.PostForm("action_type")
	if len(actionType) == 0 {
		actionType = c.Query("action_type")
	}
	content := c.PostForm("content")
	if len(content) == 0 {
		content = c.Query("content")
	}
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to parse toUserId",
		})
	}
	msi := service.MessageServiceImpl{}
	if actionType == "1" {
		nowTime := time.Now()
		_, err := msi.AddMessageByAll(fromUserId, toUserId, content, nowTime)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: 2, StatusMsg: "failed to send message",
			})
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0, StatusMsg: "send message successfully",
		})
	}
}

func MessageChat(c *gin.Context) {
	myId := c.GetInt64("id")
	yourIdStr := c.Query("to_user_id")
	yourId, err := strconv.ParseInt(yourIdStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "failed to parse user id",
		})
	}
	// preMsgTime is int64
	_ = c.Query("pre_msg_time")
	msi := service.MessageServiceImpl{}
	ids := make([]int64, 0)
	tmpIds, err := msi.FindMessageIdsByFromUserIdAndToUserId(myId, yourId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2, StatusMsg: "failed to find messages between you",
		})
	}
	ids = append(ids, tmpIds...)
	tmpIds, err = msi.FindMessageIdsByFromUserIdAndToUserId(yourId, myId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2, StatusMsg: "failed to find messages between you",
		})
	}
	ids = append(ids, tmpIds...)

	//至此获得双方消息id
	messageList := GetMessageListByIds(ids)
	c.JSON(http.StatusOK, ChatResp{
		Response: Response{
			StatusCode: 0, StatusMsg: "get chat record successfully",
		},
		MessageList: messageList,
	})
}

func GetMessageListByIds(ids []int64) []Message {
	messageList := make([]Message, 0)
	var wg sync.WaitGroup
	wg.Add(len(ids))
	for _, id := range ids {
		go GetMessageById(id, &messageList, &wg)
	}
	wg.Wait()
	return messageList
}

func GetMessageById(id int64, messageList *[]Message, wg *sync.WaitGroup) {
	defer wg.Done()
	msi := service.MessageServiceImpl{}
	message := Message{}
	message.Id = id
	messageDao, err := msi.FindMessageById(id)
	if err != nil {
		log.Println("failed to get message info by id", err)
		return
	}
	message.FromUserId = messageDao.FromUserId
	message.ToUserId = messageDao.ToUserId
	message.Content = messageDao.Content
	message.CreateTime = messageDao.CreateTime.Unix()
	*messageList = append(*messageList, message)
}
