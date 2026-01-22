package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	appModel "project_m_backend/app/user/model"
	"project_m_backend/apperrors"
	"time"

	"github.com/lib/pq"
)

const(
	queryCreateUser = `INSERT INTO users(first_name, last_name, email, password_hash, position, created_at) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;`
	queryGetUserByID = `SELECT id, first_name, last_name, email, password_hash, position, created_at FROM users WHERE id =$1;`
	queryGetUserByEmail = `SELECT id, first_name, last_name, email, password_hash, position, created_at FROM users WHERE email = $1;`
)

type PostgresUserRepository struct{
	db *sql.DB
	createUserStmt *sql.Stmt
	getUserByIDStmt *sql.Stmt
	getUserByEmail *sql.Stmt
}

func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error){
	createUserStmt, err := db.PrepareContext(context.Background(), queryCreateUser)
	if err != nil{
		return nil, apperrors.NewInternal(err, "error Preparing create user statement")
	}
	getUserByIDStmt, err := db.PrepareContext(context.Background(), queryGetUserByID)
	if err != nil{
		return nil, apperrors.NewInternal(err, "error preparing get user by id")
	}
	getUserByEmailStmt, err := db.PrepareContext(context.Background(), queryGetUserByEmail)
	if err != nil{
		return nil, apperrors.NewInternal(err, "error preparing get user by email")
	}

	return &PostgresUserRepository{
		db: db,
		createUserStmt: createUserStmt,
		getUserByIDStmt: getUserByIDStmt,
		getUserByEmail: getUserByEmailStmt,
	}, nil
}


func (r *PostgresUserRepository) Close() error{
	var errs []error
	if err := r.createUserStmt.Close(); err != nil{
		errs = append(errs, err)
	}
	if err := r.getUserByIDStmt.Close(); err != nil{
		errs = append(errs, err)
	}
	if err := r.getUserByEmail.Close(); err != nil{
		errs = append(errs, err)
	}
	if len(errs) > 0{
		return apperrors.NewInternal(fmt.Errorf("multiple errors: %v", errs), "errors closing statements")
	}
	return nil
}

//create user
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *appModel.AppUser) (int64, *apperrors.AppError){
	createdAt := time.Now()
	
	var id int64
	err := r.createUserStmt.QueryRowContext(
		ctx, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.Position, createdAt,
	).Scan(&id)

	if err != nil{
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505"{
			return 0, apperrors.NewInvalidInput(err, apperrors.CodeConflict, "A user with this email alradyy Exist")
		}
		return 0, apperrors.NewInternal(err, "faild to create user")
	}
	return id,nil
}


func (r *PostgresUserRepository) GetUserById(ctx context.Context, id int64) (*appModel.AppUser, *apperrors.AppError){
	appUser := &appModel.AppUser{}

	err := r.getUserByIDStmt.QueryRowContext(ctx, id).Scan(
		&appUser.ID, &appUser.FirstName, &appUser.LastName, &appUser.Email, &appUser.PasswordHash, &appUser.Position, &appUser.CreatedAt,
	)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, apperrors.NewNotFound(err, "user not found")
		}
		return nil, apperrors.NewInternal(err, "faild to get user by id")
	}

	return appUser, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*appModel.AppUser, *apperrors.AppError){
	appUser := &appModel.AppUser{}

	err := r.getUserByEmail.QueryRowContext(ctx, email).Scan(
		&appUser.ID, &appUser.FirstName, &appUser.LastName, &appUser.Email, &appUser.PasswordHash,
		&appUser.Position, &appUser.CreatedAt,
	)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, apperrors.NewNotFound(err, "user not found")
		}
		return nil, apperrors.NewInternal(err, "faild to get user by email")
	}
	return appUser, nil
}