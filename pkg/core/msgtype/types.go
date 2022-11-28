package msgtype

const (
	// Event message
	Event = "event"
	// Alert message
	Alert = "alert"
	// Incident message
	Incident = "incident"

	// UserError e.g. no such user
	UserError = "user.error"
	// UserRequest get a use from the backend
	UserRequest = "user.request"
	// UserList get all users form the backend
	UserList = "user.list"
	// UserAdd adds or changes a user
	UserAdd = "user.add"
	// UserResponse response for user requests
	UserResponse = "user.response"
	// UserDelete deletes a use from the backend
	UserDelete = "user.delete"
)
