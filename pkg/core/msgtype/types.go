package msgtype

const (
	// Event message
	Event = "event"
	// Alert message
	Alert = "alert"
	// Incident message
	Incident = "incident"

	// UserError e.g. no such user
	UserError = "som.user.error"
	// UserRequest get a use from the backend
	UserRequest = "som.user.request"
	// UserList get all users form the backend
	UserList = "som.user.list"
	// UserAdd adds or changes a user
	UserAdd = "som.user.add"
	// UserResponse response for user requests
	UserResponse = "som.user.response"
	// UserDelete deletes a use from the backend
	UserDelete = "som.user.delete"
)
