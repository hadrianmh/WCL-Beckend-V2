package controllers

import (
	"backend/adapters"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	_, err := adapters.NewSql()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":    500,
			"message": err})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "Successfully connected and pinged the database!"})
}
