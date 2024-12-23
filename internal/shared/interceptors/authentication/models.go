package authentication

type AuthLevelOption struct {
	MethodName string
	AuthLevel  string
}

var authLevels = []string{"admin", "editor", "author", "user", "guest"}

func ValidateAuthLevel(authLevel string) bool {
	for _, level := range authLevels {
		if authLevel == level {
			return true
		}
	}
	return false
}
