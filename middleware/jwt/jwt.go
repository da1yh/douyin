package jwt

import (
	"douyin/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type MyCustomClaims struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(id int64, name, password string) string {
	claims := MyCustomClaims{
		Id:       id,
		Name:     name,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			Audience:  name,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			Id:        strconv.FormatInt(id, 10),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "admin",
			NotBefore: time.Now().Unix(),
			Subject:   "login",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Println("error: ", err.Error())
		return ""
	}
	return token
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if len(tokenString) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 1, StatusMsg: "unauthorized",
			})
		}
		myCustomClaims, err := ParseToken(tokenString)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 2, StatusMsg: "unauthorized, login again",
			})
		}
		c.Set("id", myCustomClaims.Id)
		c.Set("name", myCustomClaims.Name)
		c.Set("password", myCustomClaims.Password)
		c.Next()
	}
}

func AuthPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.PostForm("token")
		if len(tokenString) == 0 {
			tokenString = c.Query("token")
		}
		if len(tokenString) == 0 {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 1, StatusMsg: "unauthorized",
			})
		}
		myCustomClaims, err := ParseToken(tokenString)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: 2, StatusMsg: "unauthorized, login again",
			})
		}
		c.Set("id", myCustomClaims.Id)
		c.Set("name", myCustomClaims.Name)
		c.Set("password", myCustomClaims.Password)
		c.Next()
	}
}

func ParseToken(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
