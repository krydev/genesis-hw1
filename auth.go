package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func SignUp(ctx *gin.Context) {
	user := NewUser()
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var err error
	user.Password, err = HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	if u := GetUserByEmail(user.Email); u != nil {
		ctx.JSON(http.StatusBadRequest, "User already exists with given email")
	}
	err = AddUser(*user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to create a user")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func Login(ctx *gin.Context) {
	var data User
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := GetUserByEmail(data.Email)
	if user == nil || !CheckPasswordHash(data.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized,"Please provide valid login details")
		return
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": token})
}

func CreateToken(userId uint) (string, error) {
	tClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
		Id: strconv.Itoa(int(userId)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tClaims)
	token, err := t.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateToken(ctx *gin.Context) (*jwt.Token, error) {
	tokenString, err := ExtractToken(ctx)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token, nil
}

func ExtractToken(ctx *gin.Context) (string, error) {
	bearToken := ctx.Request.Header.Get("Authorization")
	if bearToken == "" {
		return "", errors.New("missing Authorization header")
	}
	//normally Bearer <token_value>
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return "", errors.New("malformed Authorization header")
}

// validates token and authorizes users
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := ValidateToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusForbidden, err.Error())
			ctx.Abort()
			return
		}

		ctx.Next()

	}
}