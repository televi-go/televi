package external

import (
	"strconv"
)

type Target interface {
	GetFromUserList(userList map[string]bool) []string
}

type SingleUserTarget struct {
	UserId int64
}

func (singleUserTarget SingleUserTarget) GetFromUserList(userList map[string]bool) []string {
	uString := strconv.FormatInt(singleUserTarget.UserId, 10)
	_, hasUser := userList[uString]
	if !hasUser {
		return nil
	}
	return []string{uString}
}

type AllUsersTarget struct{}

func (allUsers AllUsersTarget) GetFromUserList(userList map[string]bool) []string {
	var result []string
	for s := range userList {
		result = append(result, s)
	}
	return result
}
