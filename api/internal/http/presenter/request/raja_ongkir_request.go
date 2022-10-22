package request

type GetDeliveryRequest struct {
	Origin      string `form:"origin"`
	Destination string `form:"destination"`
	Weight      int    `form:"weight"`
	Courier     string `form:"courier"`
}
