package response

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type LoginData struct {
	Token string `json:"token"`
}

type ReportData struct {
	ID          uint   `json:"id"`
	UserName    string `json:"user_name"`
	Content     string `json:"content"`
	Place       string `json:"place"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Attachment  string `json:"attachment"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ReportListData struct {
	Reports []ReportData `json:"reports"`
}
