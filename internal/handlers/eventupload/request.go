package eventupload

import (
	"encoding/json"
	"log"
	"net/http"

	apirequest "github.com/airbnb/rudolph/pkg/request"
	"github.com/airbnb/rudolph/pkg/response"
	"github.com/aws/aws-lambda-go/events"
)

// Parses the HTTP Request into the appropriate request type, or returns a HTTP Response if something is wrong
func parseRequest(request events.APIGatewayProxyRequest) (machineID string, parsedRequest *EventUploadRequest, errorResponse *events.APIGatewayProxyResponse, err error) {
	if request.Resource != "/eventupload/{machine_id}" || request.HTTPMethod != "POST" {
		// This code is intended to be unreachable, as AWS Lambda will never route to this handler
		// with the wrong method, unless misconfigured.
		log.Printf("ASSERTION FAILED: Reached unreachable route code under /eventupload")
		errorResponse, err = response.APIResponse(http.StatusMethodNotAllowed, nil)
		return
	}

	machineID, errorResponse, err = apirequest.GetMachineID(request)
	if errorResponse != nil || err != nil {
		// Unreachable code; API Gateway should never encounter an error attempting to GetMachineID
		return
	}

	if request.Headers["content-type"] != "application/json" && request.Headers["Content-Type"] != "application/json" {
		errorResponse, err = response.APIResponse(http.StatusUnsupportedMediaType, response.ErrInvalidMediaTypeResponse)
		return
	}

	if len(request.Body) <= 0 {
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
		return
	}

	err = json.Unmarshal([]byte(request.Body), &parsedRequest)
	if err != nil {
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
		return
	}

	// Parse the request
	err = json.Unmarshal([]byte(request.Body), &parsedRequest)
	if err != nil {
		log.Printf("%s\n%s", err.Error(), "request body unmarshal was not successful")
		errorResponse, err = response.APIResponse(http.StatusBadRequest, response.ErrInvalidBodyResponse)
		return
	}

	return
}

// EventUploadRequest encapsulation of an /eventupload POST body sent by a Santa sensor
type EventUploadRequest struct {
	Events []EventUploadEvent `json:"events"`
}

// EventUploadEvent is a single event entry
type EventUploadEvent struct {
	ParentName                   string         `json:"parent_name"`
	FilePath                     string         `json:"file_path"`
	QuarantineTimestamp          int            `json:"quarantine_timestamp"`
	LoggedInUsers                []string       `json:"logged_in_users"`
	SigningChain                 []SigningEntry `json:"signing_chain"`
	SigningIDs                   string         `json:"signing_id"`
	TeamID                       string         `json:"team_id"`
	ParentProcessID              int            `json:"ppid"`
	ExecutingUser                string         `json:"executing_user"`
	FileName                     string         `json:"file_name"`
	ExecutionTime                float64        `json:"execution_time"`
	FileSHA256                   string         `json:"file_sha256"`
	Decision                     string         `json:"decision"`
	ProcessID                    int            `json:"pid"`
	CurrentSesssions             []string       `json:"current_sessions"`
	FileBundleID                 string         `json:"file_bundle_id,omitempty"`
	FileBundlePath               string         `json:"file_bundle_path,omitempty"`
	FileBundleExecutableRelPath  string         `json:"file_bundle_executable_rel_path,omitempty"`
	FileBundleName               string         `json:"file_bundle_name,omitempty"`
	FileBundleVersion            string         `json:"file_bundle_version,omitempty"`
	FileBundleShortVersionString string         `json:"file_bundle_version_string,omitempty"`
	FileBundleHash               string         `json:"file_bundle_hash,omitempty"`
	FileBundleHashMilliseconds   float64        `json:"file_bundle_hash_millis,omitempty"`
	FileBundleBinaryCount        int64          `json:"file_bundle_binary_count,omitempty"`
}

// SigningEntry is optionally present when an event includes a binary that is signed
type SigningEntry struct {
	CertificateName    string `json:"cn"`
	ValidUntil         int    `json:"valid_until"`
	Organization       string `json:"org"`
	ValidFrom          int    `json:"valid_from"`
	OrganizationalUnit string `json:"ou"`
	SHA256             string `json:"sha256"`
}

// EventPayload represents derived metadata for events uploaded with the UploadEvent endpoint.
type EventPayload struct {
	FileSHA  string          `json:"file_sha256"`
	UnixTime float64         `json:"execution_time"`
	Content  json.RawMessage `json:"-"`
}
