package controllers

import (
	"backend/adapters"
	"backend/services"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BodyRequest struct {
	Action        string `json:"action" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password"`
	ResfreshToken string `json:"refresh_token"`
}

func Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"response": gin.H{
			"message": "Welcome to WCL microservices."}})
}

func Ping(ctx *gin.Context) {
	ping, err := adapters.NewSql()
	if err != nil {
		ctx.JSON(500, gin.H{
			"code":   500,
			"status": "error",
			"response": gin.H{
				"message": err.Error()}})
		return
	}

	defer ping.Connection.Close()

	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"response": gin.H{
			"message": "Pinged the database!"}})
}

func Auth(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"status": "error",
			"response": gin.H{
				"message": "Bad Request"}})
		return
	}

	if len(body) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"status": "error",
			"response": gin.H{
				"message": "Bad Request"}})
		return
	}

	// Reset the request body so it can be read again before ShouldBindJSON
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	var BodyReq BodyRequest
	if err := ctx.ShouldBindJSON(&BodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"status": "error",
			"response": gin.H{
				"message": err.Error()}})
		return
	}

	authCheck, err := services.Auth(BodyReq.Action, BodyReq.Email, BodyReq.Password, BodyReq.ResfreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":   http.StatusBadRequest,
			"status": "error",
			"response": gin.H{
				"message": err.Error()}})
		return
	}

	ctx.JSON(200, gin.H{
		"code":     200,
		"status":   "success",
		"response": authCheck})
}

func Dashboard(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":   200,
		"status": "success",
		"response": gin.H{
			"message": "Welcome to WCL Dashboard."}})
}
