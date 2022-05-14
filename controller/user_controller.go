package controller

import (
	"apriori/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService service.UserService
}

func (controller *UserController) FindAll(c *gin.Context) {

}

func (controller *UserController) FindById(c *gin.Context) {

}

func (controller *UserController) Save(c *gin.Context) {

}

func (controller *UserController) Update(c *gin.Context) {

}

func (controller *UserController) Delete(c *gin.Context) {

}
