package user

import (
	"auth/internal/client/db"
	"auth/internal/model"
	"auth/internal/repository"
	"context"

	converter "auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "user"

	idColumn             = "id"
	nameColumn           = "name"
	phoneNumberColumn    = "email"
	hashedPasswordColumn = "password"
	roleColumn           = "role"
	createdAtColumn      = "created_at"
	updatedAtColumn      = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, phoneNumberColumn, hashedPasswordColumn, roleColumn).
		Values(info.Name, info.PhoneNumber, info.HassedPassword, info.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, phoneNumberColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Info.Name, &user.Info.PhoneNumber, &user.Info.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	if info.Name != "" {
		builder = builder.Set(nameColumn, info.Name)
	}
	if info.PhoneNumber != "" {
		builder = builder.Set(phoneNumberColumn, info.PhoneNumber)
	}
	if info.HassedPassword != "" {
		builder = builder.Set(hashedPasswordColumn, info.HassedPassword)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
