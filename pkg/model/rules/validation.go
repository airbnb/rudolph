package rules

import "regexp"

var sha256Regexp = regexp.MustCompile(`^[a-f0-9]{64}$`)

var teamIDRegexp = regexp.MustCompile(`^([A-Z0-9]{1,10})$`)

var signingIDRegexp = regexp.MustCompile(`^([A-Z0-9]{1,10}|platform)(:[\w\-\.]+)$`)

func ValidSha256(sha256 string) bool {
	return sha256Regexp.MatchString(sha256)
}

func ValidTeamID(teamID string) bool {
	return teamIDRegexp.MatchString(teamID)
}

func ValidSigningID(signingID string) bool {
	return signingIDRegexp.MatchString(signingID)
}
