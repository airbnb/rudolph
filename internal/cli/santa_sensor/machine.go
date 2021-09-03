package santa_sensor

var selfMachineID string

// GetSelfMachineID functions identically to GetMyMachineUUID but caches its results,
// so multiple successive back-to-back calls are more performant.
func GetSelfMachineID() (string, error) {
	var err error
	if selfMachineID == "" {
		selfMachineID, err = GetMyMachineUUID()
	}
	return selfMachineID, err
}
