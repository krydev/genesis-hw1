package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"sync"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetJsonStringFromUrl(url string, ctx *gin.Context) (string, error) {
	resp, err := httpTimeoutClient.Get(url)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, "Unable to provide access external service at the moment")
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to retrieve data")
		return "", err
	}
	return string(body), nil
}


type AutoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	id uint
}

func (a *AutoInc) ID() uint {
	a.Lock()
	defer a.Unlock()

	a.id++
	return a.id
}