package sdto

type ParseGPXDataInput struct {
	GPXData []byte
}

type ParseGPXDataOutput struct {
	RouteID int32
}

type ParseLonLatDataInput struct {
	LonLatData [][]string
}

type ParseLonLatDataOutput struct {
	RouteID int32
}
