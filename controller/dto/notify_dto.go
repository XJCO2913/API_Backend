package dto

type ShareRouteReq struct {
	ReceiverID string     `json:"receiverid"`
	RouteData  [][]string `json:"routeData"`
}
