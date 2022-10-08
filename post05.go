package post05

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

type Userdata struct {
	Id          int
	Username    string
	FirstName   string
	LastName    string
	Description string
}

var (
	Host     = ""
	Port     = 5432
	User     = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, Database)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userId := -1

	selectStatement := fmt.Sprintf(`select id from users where username = '%s'`, username)

	rows, err := db.Query(selectStatement)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		userId = id
	}
	defer rows.Close()

	return userId
}

func AddUser(user Userdata) int {
	user.Username = strings.ToLower(user.Username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userId := exists(user.Username)

	if userId != -1 {
		fmt.Println("User already exists", user.Username)
		return -1
	}

	insertStatement := `insert into users (username) values ($1)`
	_, err = db.Exec(insertStatement, user.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userId = exists(user.Username)

	if userId == -1 {
		return userId
	}

	insertStatement = `insert into userdata (userid, firstname, lastname, description) values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userId, user.FirstName, user.LastName, user.Description)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return userId
}

func UpdateUser(user Userdata) error {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	userId := exists(user.Username)
	if userId == -1 {
		return errors.New("user does not exist")
	}

	updateStatement := `update userdata set firstname=$1, lastname=$2, description=$3 where userid=$4`
	_, err = db.Exec(updateStatement, user.FirstName, user.LastName, user.Description, userId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(`select username from users where id = '%d'`, id)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with id %d does not exist", id)
	}

	deleteStatement := `delete from userdata where userid = $1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	deleteStatement = `delete from users where id = $1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ListUsers() ([]Userdata, error) {
	users := make([]Userdata, 0, 10)
	db, err := openConnection()
	if err != nil {
		return users, err
	}
	defer db.Close()
	fmt.Println("test")

	queryStatement := `select username, userid, firstname, lastname, description from users inner join userdata on users.id = userdata.userid`
	rows, err := db.Query(queryStatement)
	if err != nil {
		fmt.Println(err)
		return users, err
	}

	for rows.Next() {
		var username string
		var userid int
		var firstname string
		var lastname string
		var description string
		err = rows.Scan(&username, &userid, &firstname, &lastname, &description)
		if err != nil {
			fmt.Println(err)
			return users, err
		}
		user := Userdata{
			Id:          userid,
			Username:    username,
			FirstName:   firstname,
			LastName:    lastname,
			Description: description,
		}
		users = append(users, user)
	}
	defer rows.Close()

	return users, nil
}
