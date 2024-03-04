package handlers

import "github.com/gin-gonic/gin"

type Public interface {
	SpecialistRegister(c *gin.Context)
	SpecialistLogin(c *gin.Context)

	CameraCreate(c *gin.Context)
	CameraDelete(c *gin.Context)

	CaseCreate(c *gin.Context)

	Refresh(c *gin.Context)
}
