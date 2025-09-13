package models

// Command represents a command sent to a client
type Command struct {
	TargetID string `json:"target_id"` // Target client ID ("all" for all)
	Action   string `json:"action"`    // Command to execute (e.g., "TAKE_PHOTO")
	Payload  string `json:"payload"`   // Additional parameters
}

type CommandResponse struct {
	ClientID string      `json:"client_id"`
	Action   string      `json:"action"`  // The action that was executed
	Success  bool        `json:"success"` // Whether the command was successful
	Result   interface{} `json:"result"`  // Result data from the command
	Error    string      `json:"error"`   // Error message if the command failed
}
