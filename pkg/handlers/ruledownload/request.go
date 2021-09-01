package ruledownload

// RuledownloadRequest is the postbody submitted to /ruledownload endpoints
type RuledownloadRequest struct {
	// Cursor is, verbatim, the Cursor that is returned to a sensor in a previous RuledownloadResponse
	// On the very first rule download request in a flight sequence, there will be no cursor provided.
	Cursor *ruledownloadCursor `json:"cursor,omitempty"`
}
