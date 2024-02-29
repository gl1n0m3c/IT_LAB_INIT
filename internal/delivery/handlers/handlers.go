package handlers

import "github.com/gin-gonic/gin"

type Public interface {
	SpecialistRegister(c *gin.Context)
	SpecialistLogin(c *gin.Context)
	Refresh(c *gin.Context)
}
