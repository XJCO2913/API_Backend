package sdto

type Organiser struct {
	UserID         string `json:"userId"`
	Username       string `json:"username"`
	AvatarUrl      string `json:"avatarUrl"`
	MembershipTime int64 `json:"membershipTime"`
	Status         string `json:"status"`
}

type GetAllOrganisersOutput struct {
	Organisers []Organiser
}
