package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"net/http"
	"time"
)

var httpTimeoutClient = &http.Client{Timeout: 10 * time.Second}

func btcRate(ctx *gin.Context){
	// Getting currency rates
	jsonRates, err := GetJsonStringFromUrl("https://api.privatbank.ua/p24api/pubinfo?exchange&json&coursid=11", ctx)
	if err != nil {
		return
	}
	usdToUahPath := "0.sale"
	btcToUsdPath := "3.sale"
	rates := gjson.GetMany(jsonRates, usdToUahPath, btcToUsdPath)

	// Converting price in usd to uah
	usdToUah := rates[0].Float()
	btcToUsd := rates[1].Float()
	btcToUahRate := btcToUsd * usdToUah
	ctx.JSON(http.StatusOK, gin.H{"btcToUah": fmt.Sprintf("%.2f", btcToUahRate)})
}