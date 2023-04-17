package ftp

import (
	"douyin/config"
	"github.com/jlaffaye/ftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"time"
)

var c *ftp.ServerConn
var client *ssh.Client

// var session *ssh.Session
var err error

func InitFtp() {
	//ftp本身就是长连接
	c, err = ftp.Dial(config.FtpServerAddr+":"+config.FtpServerPort, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println(err)
		return
	}
	err = c.Login(config.FtpServerUserName, config.FtpServerPassword)
	if err != nil {
		log.Println(err)
		return
	}
}

func InitSSH() {
	sshConfig := &ssh.ClientConfig{
		User: config.SSHServerUserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SSHServerPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	client, err = ssh.Dial("tcp", config.SSHServerAddr+":"+config.SSHServerPort, sshConfig)
	if err != nil {
		log.Println(err)
		return
	}
	go keepAlive()
}

func keepAlive() {
	for {
		select {
		case <-time.After(10 * time.Second):
			session, err := client.NewSession()
			if err != nil {
				log.Println(err)
			}
			session.Close()
		}
	}
}

func UploadVideo(fileName string, data io.Reader) error {
	err = c.Stor("video/"+fileName, data)
	if err != nil {
		log.Println("error: ", err.Error())
		return err
	}
	return nil
}

// Screenshot use ffmpeg
func Screenshot(videoFileName, imageFileName string) error {
	//是否要判断视频是否存在
	session, err := client.NewSession()
	if err != nil {
		log.Println(err)
		return err
	}
	defer session.Close()
	cmd := "/usr/local/ffmpeg/bin/ffmpeg -i " + config.FtpServerAddrPrefix + "/video/" + videoFileName + " -ss 00:00:01.000 -vframes 1 " + config.FtpServerAddrPrefix + "/image/" + imageFileName
	err = session.Run(cmd)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
