package repository

import (
	"database/sql"
	"project/guidemysteps/src/internal/core/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Get(email string) (*models.User, error) {
	/* EXAMPLES
		rows, err := db.Query(`SELECT "Name", "Roll_Number" FROM "Students"`)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
	    var name string
	    var roll_number int

	    err = rows.Scan(&name, &roll_number)
	    CheckError(err)

	    fmt.Println(name, roll_number)
	}*/
	user := models.User{UserName: "german@gmail.com"}
	user = models.NewUser(user)
	return &user, nil
}

func (u *userRepository) Create(user models.User) error {
	// dynamic
	/*insertDynStmt := `insert into "Students"("Name", "Roll_Number") values($1, $2)`
	  _, e = db.Exec(insertDynStmt, "Jack", 21)*/
	return nil
}

func (u *userRepository) Update(user models.User) error {
	/*// update
	updateStmt := `update "Students" set "Name"=$1, "Roll_Number"=$2 where "id"=$3`
	_, e := db.Exec(updateStmt, "Rachel", 24, 8)*/
	return nil
}
