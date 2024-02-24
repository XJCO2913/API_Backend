package errorx

// Error code constants
const (
	ErrInternal = 500
	ErrExternal = 400
)

type ServiceErr struct {
	errCode int
	errMsg  string
	errData map[string]any
}

// Wrap a service layer error
func NewServicerErr(code int, msg string, data map[string]any) *ServiceErr {
	return &ServiceErr{
		errCode: code,
		errMsg:  msg,
		errData: data,
	}
}

func NewInternalErr() *ServiceErr {
	return &ServiceErr{
		errCode: ErrInternal,
		errMsg: "Internal server error",
		errData: make(map[string]any),
	}
}

// Error() implements error interface, return error message
func (s ServiceErr) Error() string {
	return s.errMsg
}

// Code() returns wrapping error code
func (s ServiceErr) Code() int {
	return s.errCode
}

// Get() returns wrapping error data value by key
func (s ServiceErr) Get(key string) any {
	return s.errData[key]
}

// Set() wraps a new key-val in error data
func (s ServiceErr) Set(key string, val any) {
	s.errData[key] = val
}
