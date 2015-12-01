package appinfo

import "strings"

func NormalizeAppID(candidateID string) string {
	normalizedAppID := strings.ToLower(NormalizeAppIDWithCaseSensitive(candidateID))
	return normalizedAppID
}

func NormalizeAppIDWithCaseSensitive(candidateID string) string {
	normalizedAppID := strings.Replace(candidateID, "_", "-", -1)
	return normalizedAppID
}
