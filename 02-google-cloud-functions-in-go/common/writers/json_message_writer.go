package writers

import "encoding/json"

// JSONMessageWriter struct type
type JSONMessageWriter struct {
	Message string `json:"message"`
}

// NewMessageWriter constructs a new JSONMessageWriter
func NewMessageWriter(message string) *JSONMessageWriter {
	return &JSONMessageWriter{
		Message: message,
	}
}

// JSONString returns a string equivalent
func (jw *JSONMessageWriter) JSONString() (string, error) {
	messageResponse := map[string]interface{}{
		"data": map[string]string{
			"message": jw.Message,
		},
	}

	bytesValue, err := json.Marshal(messageResponse)
	if err != nil {
		return "", err
	}
	return string(bytesValue), nil
}
