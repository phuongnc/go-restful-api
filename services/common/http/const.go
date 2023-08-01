package http

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "DELETE"
)

type Code uint32

const (
	InvalidAccount Code = 1000
	SessionExpired Code = 1001
)

var CodeToString = map[Code]string{
	InvalidAccount: "The email or password is incorrect",
	SessionExpired: "Your session has expired, please login again",
}
