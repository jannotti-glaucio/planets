package communication

//Fields standardizing fields
type Fields struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

//Response standardizing responses
type Response struct {
	Status  int         `json:"status,omitempty"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Fields  []Fields    `json:"fields,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

//AddFields add s
func (e *Response) AddFields(field, message string) {
	e.Fields = append(e.Fields, Fields{Field: field, Message: message})
}
