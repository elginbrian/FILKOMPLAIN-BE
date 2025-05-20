package request

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateReportRequest struct {
	UserName    string `json:"user_name"`
	Content     string `json:"content"`
	Place       string `json:"place"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Attachment  string `json:"attachment"`
}

type UpdateReportRequest struct {
	UserName    string `json:"user_name"`
	Content     string `json:"content"`
	Place       string `json:"place"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Attachment  string `json:"attachment"`
}

type ResolveReportStatusRequest struct {
	Status string `json:"status"`
}

type ReplyReportRequest struct {
	Reply string `json:"reply"`
}
