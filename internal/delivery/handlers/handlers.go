package handlers

import "github.com/gin-gonic/gin"

type Public interface {
	ManagerLogin(c *gin.Context)

	SpecialistRegister(c *gin.Context)
	SpecialistLogin(c *gin.Context)

	CameraCreate(c *gin.Context)
	CameraDelete(c *gin.Context)

	CaseCreate(c *gin.Context)

	Refresh(c *gin.Context)
}

type Specialists interface {
	GetMe(c *gin.Context)
	UpdateMe(c *gin.Context)

	GetCasesByLevel(c *gin.Context)

	CreateRated(c *gin.Context)
	GetRatedSolved(c *gin.Context)
}
