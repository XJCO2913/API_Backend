package sdto

type ParseGPXDataInput struct {
	GPXData []byte
}

type ParseGPXDataOutput struct {
	RouteID int32
}
