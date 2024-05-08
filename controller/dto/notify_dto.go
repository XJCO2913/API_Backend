package dto

type ShareRouteReq struct {
	ReceiverID string     `json:"receiverId"`
	RouteData  [][]string `json:"routeData"`
}

type OrgResultReq struct {
	ReceiverID string `json:"receiverId"`
	IsAgreed   bool   `json:"isAgreed"`
}
