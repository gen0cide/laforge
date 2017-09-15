package competition

type UserList map[string][]User

type User struct {
	Firstname  string `json:"first_name"`
	Lastname   string `json:"last_name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Title      string `json:"title"`
	Password   string `json:"password"`
	OU         string `json:"ou"`
}
