package authentication

import "slices"

type AuthLevelOption struct {
	MethodName string
	AuthLevel  string
	Self       bool
}

// Important hold order
var authLevels = []string{"admin", "editor", "author", "user", "guest"}

func ValidateAuthLevel(authLevel string) bool {
	for _, level := range authLevels {
		if authLevel == level {
			return true
		}
	}
	return false
}

func CheckEnoughAuthLevel(requiredLevel, authLevel string) bool {
	reqId := slices.Index(authLevels, requiredLevel)
	authId := slices.Index(authLevels, authLevel)

	if authId == -1 {
		return false
	}

	if authId >= reqId {
		return true
	}

	return false
}
