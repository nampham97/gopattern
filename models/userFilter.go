// file: models/user_filter.go
package models

type UserFilter struct {
	Search string
	Role   string
	Page   int
	Limit  int
}
