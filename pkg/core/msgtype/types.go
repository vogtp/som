package msgtype

const (
	// Event message
	Event = "som.szenario.event"
	// Alert message
	Alert = "som.szenario.alert"
	// Incident message
	Incident = "som.szenario.incident"

	// UserError e.g. no such user
	UserError = "som.user.error"
	// UserRequest get a use from the backend
	UserRequest = "som.user.request"
	// UserList get all users form the backend
	UserList = "som.user.list"
	// UserAdd adds or changes a user
	UserAdd = "som.user.add"
	// UserDelete deletes a use from the backend
	UserDelete = "som.user.delete"
)
