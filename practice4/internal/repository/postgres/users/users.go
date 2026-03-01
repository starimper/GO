package users

import (
	"errors"
	"time"

	"practice_4/internal/repository/postgres"
	"practice_4/pkg/modules"
)

type Repository struct {
	db      *postgres.Dialect
	timeout time.Duration
}

func NewUserRepository(db *postgres.Dialect) *Repository {
	return &Repository{
		db:      db,
		timeout: 5 * time.Second,
	}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT * FROM users")
	return users, err
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *Repository) CreateUser(user *modules.User) (int, error) {
	var id int
	query := `
	INSERT INTO users (name, email, age)
	VALUES ($1,$2,$3)
	RETURNING id`

	err := r.db.DB.QueryRow(query,
		user.Name,
		user.Email,
		user.Age,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateUser(user *modules.User) error {
	result, err := r.db.DB.Exec(`
	UPDATE users
	SET name=$1, email=$2, age=$3
	WHERE id=$4`,
		user.Name,
		user.Email,
		user.Age,
		user.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user does not exist")
	}

	return nil
}

func (r *Repository) DeleteUser(id int) (int64, error) {
	result, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
