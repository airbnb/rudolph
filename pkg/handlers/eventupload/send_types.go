package eventupload

// ForwardedEventUploadEvent is a single event entry that appends the MachineID with the EventUploadEvent details
// and is send to Firehose
type ForwardedEventUploadEvent struct {
	MachineID string `json:"machine_id"`
	EventUploadEvent
}

func convertRequestEventsToUploadEvents(machineID string, events []EventUploadEvent) []interface{} {
	var forwardedEvents []interface{}

	for _, event := range events {
		forwardedEvents = append(forwardedEvents, ForwardedEventUploadEvent{
			MachineID:        machineID,
			EventUploadEvent: event,
		})
	}

	return forwardedEvents
}
