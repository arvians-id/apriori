package rest

import (
	"encoding/json"
	"fmt"
	"github.com/arvians-id/apriori/http/controller/rest/request"
	"github.com/arvians-id/apriori/http/controller/rest/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

type RajaOngkirController struct {
}

func NewRajaOngkirController() *RajaOngkirController {
	return &RajaOngkirController{}
}

func (controller *RajaOngkirController) Route(router *gin.Engine) *gin.Engine {
	router.GET("/api/raja-ongkir/:place", controller.FindAll)
	router.POST("/api/raja-ongkir/cost", controller.GetCost)
	return router
}

func (controller *RajaOngkirController) FindAll(c *gin.Context) {
	placeParam := c.Param("place")
	if placeParam == "province" {
		placeParam = "province"
	} else if placeParam == "city" {
		placeParam = "city?province=" + c.Query("province")
	}

	url := "https://api.rajaongkir.com/starter/" + placeParam
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var rajaOngkirModel interface{}
	err := json.NewDecoder(res.Body).Decode(&rajaOngkirModel)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rajaOngkirModel,
	})
}

func (controller *RajaOngkirController) GetCost(c *gin.Context) {
	var requestDelivery request.GetDeliveryRequest
	err := c.ShouldBind(&requestDelivery)
	if err != nil {
		response.ReturnErrorBadRequest(c, err, nil)
		return
	}

	payload := fmt.Sprintf(
		"origin=%v&destination=%v&weight=%v&courier=%v",
		requestDelivery.Origin,
		requestDelivery.Destination,
		requestDelivery.Weight,
		requestDelivery.Courier,
	)
	data := strings.NewReader(payload)
	req, _ := http.NewRequest("POST", "https://api.rajaongkir.com/starter/cost", data)
	req.Header.Add("key", os.Getenv("RAJA_ONGKIR_SECRET_KEY"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var rajaOngkirModel interface{}
	err = json.NewDecoder(res.Body).Decode(&rajaOngkirModel)
	if err != nil {
		response.ReturnErrorInternalServerError(c, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": rajaOngkirModel,
	})
}
