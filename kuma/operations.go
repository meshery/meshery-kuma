package kuma

// Operation holds the informormation for list of operations
type Operation struct {
	Type       string            `json:"type,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Operations hold a map of Operation objects
type Operations map[string]*Operation

// ApplyOperation applies the operation on kuma
func (h *handler) ApplyOperation() error {
	// Needs implementation
	return nil
}

// ListOperations lists the operations available
func (h *handler) ListOperations() (Operations, error) {
	// Needs implementation
	operations := make(Operations, 0)
	h.config.Operations(&operations)
	return operations, nil
}
