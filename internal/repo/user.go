package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"user/config"
	"user/internal/entities"
)

type User struct {
	db  *sql.DB
	cfg *config.EnvConfig
}

type UserImply interface {
	FetchUserByEmail(context.Context, string) (entities.User, error)
	InsertTokens(context.Context, []entities.Token, int) error
	CreateUser(context.Context, entities.User) error
	ValidateToken(context.Context, entities.Token) (string, time.Time, error)
	GetPermissions(context.Context, string, string) (map[string]map[string]string, error)
	FetchPermissionIDs(context.Context, string, string) (map[string][]string, error)
	FetchPermissionNames(context.Context, []string) (map[string]string, error)
	AddUserPermissions(context.Context, entities.UserProjectPermissions) error
	ListAllPermissions(context.Context) (*entities.Permissions, error)
	ListUser(ctx context.Context, pagination entities.Pagination,
		filter entities.UserFilter) ([]*entities.User, int64, error)
	ResetPassword(ctx context.Context, request entities.ResetPassword) error
	IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)
	RevokeToken(ctx context.Context, tokenID string) error
	ListUserPermissions(context.Context, string, string) ([]entities.PermissionDetail, error)
}

func NewUser(db *sql.DB, cfg *config.EnvConfig) UserImply {
	return &User{
		db:  db,
		cfg: cfg,
	}
}

func (usr *User) ListUser(ctx context.Context, pagination entities.Pagination,
	filter entities.UserFilter) ([]*entities.User, int64, error) {

	var (
		queryBuilder strings.Builder
		args         []interface{}
		total        int64
	)

	queryBuilder.WriteString(listUserQ)

	addFilter(&queryBuilder, &args, "u.userid", filter.ID)
	addFilter(&queryBuilder, &args, "u.name", filter.Name)
	addFilter(&queryBuilder, &args, "u.mailid", filter.Email)
	addFilter(&queryBuilder, &args, "u.organization", filter.Organization)

	if filter.ProjectID != "" {
		addFilter(&queryBuilder, &args, "up.project_id", filter.ProjectID)

		if filter.IsMember == "" {
			filter.IsMember = "true"
		}
		addFilter(&queryBuilder, &args, "up.is_active", filter.IsMember)
	}

	switch filter.IsActive {
	case "all":
		queryBuilder.WriteString(" AND u.is_active IN (true, false)")
	case "false":
		queryBuilder.WriteString(" AND u.is_active = false")
	default:
		queryBuilder.WriteString(" AND u.is_active = true")
	}

	applyPagination(&queryBuilder, &args, pagination)

	query := queryBuilder.String()
	rows, err := usr.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute user query: %v", err)
		return nil, 0, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.UserID, &user.Name, &user.Email,
			&user.Purpose, &user.ISActive, &user.Organization,
			&user.CreatedOn, &user.UpdatedOn,
			&total,
		); err != nil {
			log.Printf("failed to scan user row: %v", err)
			return nil, 0, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error occurred while iterating rows: %v", err)
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return users, total, nil
}

func addFilter(builder *strings.Builder, args *[]interface{}, column string, value string) {
	if value != "" {
		condition := fmt.Sprintf(" AND %s = $%d", column, len(*args)+1)
		builder.WriteString(condition)
		*args = append(*args, value)
	}
}

func applyPagination(builder *strings.Builder, args *[]interface{}, pagination entities.Pagination) {
	limit := pagination.Limit
	offset := (pagination.Page - 1) * limit

	builder.WriteString(fmt.Sprintf(
		" ORDER BY u.%s %s LIMIT $%d OFFSET $%d",
		pagination.Sort, pagination.Order, len(*args)+1, len(*args)+2,
	))

	*args = append(*args, limit, offset)
}
