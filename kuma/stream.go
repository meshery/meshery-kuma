package kuma

type Event struct {
	Operationid string `json:"operationid,omitempty"`
	EType       int32  `json:"type,string,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Details     string `json:"details,omitempty"`
}

// StreamErr handles the error stream requests
func (h *handler) StreamErr(e *Event, err error) {
	h.log.Err("Sending error event", err.Error())
	e.EType = 2
	*h.channel <- e
}

// StreamInfo handles the info stream requests
func (h *handler) StreamInfo(e *Event) {
	h.log.Info("Sending event")
	e.EType = 0
	*h.channel <- e
}
