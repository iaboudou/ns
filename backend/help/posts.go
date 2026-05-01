package help

import "strings"

func ParseAllowedUsers(input string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, ",")
}
