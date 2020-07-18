package actions

import "github.com/gookit/color"

// User gitlab user data
type User struct {
	Name     string
	Username string
}

// Print user data
func (u *User) Print() {
	color.Bold.Print(u.Name)
	color.Reset()
	color.OpItalic.Printf(" (%s)\n", u.Username)
}
