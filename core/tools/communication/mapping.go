package communication

type mapItem struct {
	Message string
	Code    int
}

//ResponseMapping responsable to mapping responses
type ResponseMapping struct {
	Mapping map[string]mapItem
}

var singletonResponseMapping *ResponseMapping = nil

//New ...
func New() ResponseMapping {
	if singletonResponseMapping == nil {
		mapping := ResponseMapping{}
		mapping.populate()
		singletonResponseMapping = &mapping
	}
	return *singletonResponseMapping
}

//Response responsible to make response
func (response *ResponseMapping) Response(status int, identifier string) Response {
	return Response{
		Status:  status,
		Code:    response.Mapping[identifier].Code,
		Message: response.Mapping[identifier].Message,
	}
}

//Fields responsible to make fields
func (response *ResponseMapping) Fields(field string, identifier string) Fields {
	return Fields{
		Field:   field,
		Code:    response.Mapping[identifier].Code,
		Message: response.Mapping[identifier].Message,
	}
}

func (response *ResponseMapping) populate() {
	data := make(map[string]mapItem)

	data["already_exists"] = mapItem{Message: "Already exists", Code: 100000}
	data["validate_required"] = mapItem{Message: "Required", Code: 100001}
	data["validate_invalid"] = mapItem{Message: "Invalid", Code: 100002}
	data["validate_email"] = mapItem{Message: "Invalid e-mail", Code: 100003}
	data["validate_date"] = mapItem{Message: "Invalid date", Code: 100004}
	data["validate_password_length"] = mapItem{Message: "Password must be between 6 and 40 characters.", Code: 100005}
	data["success"] = mapItem{Message: "Success", Code: 100006}
	data["not_found"] = mapItem{Message: "Not found", Code: 100007}
	data["error"] = mapItem{Message: "Error", Code: 100008}
	data["success_create"] = mapItem{Message: "Record successfully created", Code: 100009}
	data["success_update"] = mapItem{Message: "Record successfully updated", Code: 100010}
	data["success_delete"] = mapItem{Message: "Record successfully deleted", Code: 100011}
	data["error_create"] = mapItem{Message: "Unable to create record", Code: 100012}
	data["error_update"] = mapItem{Message: "Unable to update record", Code: 100013}
	data["error_delete"] = mapItem{Message: "Unable to delete record", Code: 100014}
	data["error_list"] = mapItem{Message: "Unable to list record", Code: 100015}
	data["authenticate_failed"] = mapItem{Message: "Authenticate failed", Code: 100016}
	data["authenticate_success"] = mapItem{Message: "Authenticate success", Code: 100017}
	data["validate_failed"] = mapItem{Message: "Validation failed", Code: 100018}
	data["end_point_not_found"] = mapItem{Message: "Endpoint not found", Code: 100019}
	data["unexpected"] = mapItem{Message: "Unexpected", Code: 100020}

	response.Mapping = data
}
