package post05

import (
	"fmt"
	"testing"
)

func TestAddUser(t *testing.T) {

	Host = "localhost"
	User = "postgres"
	Password = "kuan5603"
	Database = "go"

	//AddUser(Userdata{
	//	Username:    "oliviasuper",
	//	FirstName:   "Olivia",
	//	LastName:    "Cheng",
	//	Description: "a daughter",
	//})

	//UpdateUser(Userdata{
	//	Username:    "jasonvitagen4",
	//	FirstName:   "Jayson",
	//	LastName:    "Chong",
	//	Description: "a swimmer",
	//})

	//err := DeleteUser(9)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	users, err := ListUsers()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(users)
}
