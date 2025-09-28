package request

type StaffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type LoginStaffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
