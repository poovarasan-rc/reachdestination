package model

type Debug struct {
	Sts string `json:"sts"`
	Msg string `json:"msg"`
}

type AccountStruct struct {
	Email   string `json:"email,omitempty"`
	AccID   string `json:"accid,omitempty"`
	AccName string `json:"accname,omitempty"`
	Secret  string `json:"appsecret,omitempty"`
}

type Destination struct {
	AccountID  string            `json:"account_id,omitempty"`
	DesID      string            `json:"des_id,omitempty"`
	URL        string            `json:"url,omitempty"`
	HTTPMethod string            `json:"http_method,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type ConstructData struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	JBody   string `json:"jbody"`
	Headers string `json:"headers"`
}
