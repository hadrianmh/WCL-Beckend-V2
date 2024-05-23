package controllers

import (
	"backend/adapters"
	"backend/services"
	"backend/utils"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BodyRequest struct {
	Action string `json:"action" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Pwd    string `json:"pwd"`
}

func Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "Welcome to WCL microservices."})
}

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

func Auth(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request"})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request"})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequest
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Bad Request"})
		return
	}

	authCheck, err := services.Auth(utils.Ucfirst(utils.StrReplaceAll(BodyReq.Action, "/", "")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    200,
		"message": authCheck})
}
