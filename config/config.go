package config

// mysql
const MysqlUserName string = "root"
const MysqlPassword string = "123456"
const MysqlIp string = "localhost"
const MysqlPort string = "3306"
const MysqlDbName string = "douyin"

// jwt
const SecretKey = "douyin"

// gin
const ServerAddr = "http://localhost:8080"

// ftp server
const FtpServerAddr = "120.26.161.171"
const FtpServerPort = "21"
const FtpServerUserName = "root"
const FtpServerPassword = "7355608Rushb"
const FtpServerAddrPrefix = "/ftp"

// ssh
const SSHServerUserName = "root"
const SSHServerPassword = "7355608Rushb"
const SSHServerAddr = "120.26.161.171"
const SSHServerPort = "22"

// videoPath
const PlayPathPrefix = "120.26.161.171/video/"
const CoverPathPrefix = "120.26.161.171/image/"

// redis
const RedisServerAddr = "127.0.0.1:6379"
const RedisServerPwd = "douyin"

// rabbitmq
const RabbitMQServerAddr = "amqp://guest:guest@localhost:5672/"
