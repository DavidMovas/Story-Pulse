package contracts

var (
	roles = []string{"admin", "editor", "author", "user", "guest"}
)

type Role string

func ValidateRole(role Role) bool {
	for _, v := range roles {
		if v == string(role) {
			return true
		}
	}
	return false
}
