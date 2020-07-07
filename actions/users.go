package actions

import "github.com/gookit/color"

// User gitlab user data
type User struct {
	Name     string
	Username string
}

// Print user data
func (u *User) Print() {
	color.LightBlue.Println("%s (%s)", u.Name, u.Username)
}
