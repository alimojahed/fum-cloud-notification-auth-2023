package services

type RegisterReq struct {
	Email    string
	Password string
}

type LoginReq struct {
	Email    string
	Password string
}
