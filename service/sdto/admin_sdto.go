package sdto

type AdminAuthenticateInput struct {
	Name     string
	Password string
}

type AdminAuthenticateOuput struct {
	Token   string
	AdminId int32
	Name    string
}
