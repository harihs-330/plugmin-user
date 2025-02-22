package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"user/internal/entities"
	"user/utils/querylib"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// fetch user details based on input request mail
func (usr *User) FetchUserByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User

	query := `
        SELECT userid, name, password, mailid, purpose, organization, created_on, updated_on, is_deleted, deleted_on
        FROM users
        WHERE mailid = $1 AND is_deleted = FALSE;
    `

	row := usr.db.QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.UserID,
		&user.Name,
		&user.HashedPassword,
		&user.Email,
		&user.Purpose,
		&user.Organization,
		&user.CreatedOn,
		&user.UpdatedOn,
		&user.IsDeleted,
		&user.DeletedOn,
	)
	if err != nil {
		log.Printf("FetchUserByEmail: failed to scan user with email %s: %v", email, err)

		return entities.User{}, err
	}

	return user, nil
}

// insert the generated tokens to database with corresponding user id
func (usr *User) InsertTokens(ctx context.Context, tokens []entities.Token, columnCount int) error {
	if len(tokens) == 0 {
		return nil
	}

	// Generate the placeholders for the query
	placeholders := querylib.GeneratePlaceholders(len(tokens), columnCount)

	// Format the query with placeholders
	//nolint
	query := fmt.Sprintf(`
        INSERT INTO tokens (
            id, user_id, token, token_type, created_on, expires_on
        ) VALUES %s`, placeholders)

	// Collect arguments for the query
	valueArgs := make([]interface{}, 0, len(tokens)*columnCount)
	for _, token := range tokens {
		valueArgs = append(valueArgs, token.ID, token.UserID,
			token.Token, token.TokenType, token.CreatedOn, token.ExpiresOn,
		)
	}

	// Execute the query with the collected arguments
	_, err := usr.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		log.Printf("InsertTokens: failed to insert tokens: %v", err)
		return fmt.Errorf("failed to insert tokens: %w", err)
	}

	return nil
}

// CreateUser inserts a new user into the database with the details provided in the request.
func (usr *User) CreateUser(ctx context.Context, request entities.User) error {
	// Define the SQL query to insert a new user into the users table.
	query := `
INSERT INTO users 
(name, password, mailid, purpose, organization, created_on, updated_on)
VALUES ($1, $2, $3, $4, $5, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC');
`

	// Execute the query with the provided user details.
	_, err := usr.db.ExecContext(ctx, query,
		request.Name,           // User's name
		request.HashedPassword, // User's hashed password
		request.Email,          // User's email address
		request.Purpose,        // User's purpose (e.g., student, developer)
		request.Organization,   // User's organization (optional)
	)

	// Return nil if the insertion is successful.
	return err
}

func (usr *User) ValidateToken(ctx context.Context, token entities.Token) (string, time.Time, error) {
	var (
		userID    string
		expiresOn time.Time
	)

	query := `
	SELECT user_id, expires_on
	FROM tokens
	WHERE token = $1 
	AND token_type = $2
`

	err := usr.db.QueryRowContext(
		ctx, query, token.Token, token.TokenType,
	).Scan(&userID, &expiresOn)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return "", time.Time{}, err
	}

	return userID, expiresOn, nil
}

func (usr *User) GetPermissions(ctx context.Context, userID, projectID string) (map[string]map[string]string, error) {
	permissionsMap := make(map[string]map[string]string)

	// Step 1: Fetch permission IDs
	permissionsByProject, err := usr.FetchPermissionIDs(ctx, userID, projectID)
	if err != nil {
		log.Printf("Error fetching permission IDs for userID %s and projectID %s: %v", userID, projectID, err)
		return nil, err
	}

	// Step 2: Fetch permission names in bulk for each project
	for projID, permissionIDs := range permissionsByProject {
		if len(permissionIDs) == 0 {
			continue
		}
		permissions, err := usr.FetchPermissionNames(ctx, permissionIDs)
		if err != nil {
			log.Printf("Error fetching permission names for projectID %s: %v", projID, err)
			return nil, err
		}

		permissionsMap[projID] = permissions
	}

	return permissionsMap, nil
}

// Helper function to fetch permission IDs
func (usr *User) FetchPermissionIDs(ctx context.Context, userID, projectID string) (map[string][]string, error) {
	permissionsByProject := make(map[string][]string)
	var query string
	var rows *sql.Rows
	var err error

	if projectID != "" {
		// Fetch permission IDs based on userID and projectID
		query = "SELECT permission_ids FROM user_permissions WHERE user_id = $1 AND project_id = $2"
		rows, err = usr.db.QueryContext(ctx, query, userID, projectID)
		if err != nil {
			log.Printf("Error executing query for project-specific permissions: %v", err)
			return nil, err // Handle query execution error
		}
	} else {
		// Fetch permission IDs based on userID across all projects
		query = "SELECT project_id, permission_ids FROM user_permissions WHERE user_id = $1"
		rows, err = usr.db.QueryContext(ctx, query, userID)
		if err != nil {
			log.Printf("Error executing query for all-projects permissions: %v", err)
			return nil, err // Handle query execution error
		}
	}

	defer rows.Close()

	for rows.Next() {
		var permissionIDArray []string
		var projID string

		if projectID != "" {
			projID = projectID
			if err := rows.Scan(pq.Array(&permissionIDArray)); err != nil {
				log.Printf("Error scanning rows for projectID %s: %v", projectID, err)
				return nil, err
			}
		} else {
			if err := rows.Scan(&projID, pq.Array(&permissionIDArray)); err != nil {
				log.Printf("Error scanning rows for userID %s: %v", userID, err)
				return nil, err
			}
		}
		permissionsByProject[projID] = append(permissionsByProject[projID], permissionIDArray...)
	}
	if rows.Err() != nil {
		log.Printf("Error encountered while iterating over rows: %v", rows.Err())
		return nil, rows.Err() // Handle errors encountered during row iteration
	}

	return permissionsByProject, nil
}

// Helper function to fetch permission names
func (usr *User) FetchPermissionNames(ctx context.Context, permissionIDs []string) (map[string]string, error) {
	permQuery := "SELECT id, name FROM permissions WHERE id = ANY($1)"
	permRows, err := usr.db.QueryContext(ctx, permQuery, pq.Array(permissionIDs))
	if err != nil {
		log.Printf("Error executing permission query for permission IDs %v: %v", permissionIDs, err)
		return nil, err
	}
	defer permRows.Close()

	permissions := make(map[string]string)
	for permRows.Next() {
		var permID, permName string
		if err := permRows.Scan(&permID, &permName); err != nil {
			log.Printf("Error scanning permission row with permission ID: %v", err)
			return nil, err
		}
		// Store permission ID and name in the map
		permissions[permID] = permName
	}

	// Check if there were any errors during the iteration
	if err := permRows.Err(); err != nil {
		log.Printf("Error encountered during permission row iteration: %v", err)
		return nil, err
	}

	return permissions, nil
}

func (usr *User) AddUserPermissions(ctx context.Context, request entities.UserProjectPermissions) error {
	// Begin a transaction
	tx, err := usr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Step 1: Update existing record
	updateQuery := `
	  UPDATE user_permissions 
	  SET permission_ids = $3, 
	      is_active = $4
	  WHERE user_id = $1 AND project_id = $2
	`
	_, err = tx.ExecContext(ctx, updateQuery, request.UserID, request.ProjectID, pq.Array(request.PermissionIDs), true)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Step 2: Insert if no rows were updated
	insertQuery := `
	  INSERT INTO user_permissions (user_id, project_id, permission_ids, is_active) 
	  SELECT $1, $2, $3, $4 
	  WHERE NOT EXISTS (
	    SELECT 1 FROM user_permissions WHERE user_id = $1 AND project_id = $2
	  )
	`
	_, err = tx.ExecContext(ctx, insertQuery, request.UserID, request.ProjectID, pq.Array(request.PermissionIDs), true)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (usr *User) ListAllPermissions(_ context.Context) (*entities.Permissions, error) {
	query := `
		SELECT p.id, p.name, p.description
		FROM permissions p
		WHERE p.is_deleted = false
	`

	permissions := make(map[uuid.UUID]entities.PermissionDetail)

	rows, err := usr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and populate the map
	for rows.Next() {
		var permissionID uuid.UUID
		var permissionName, permissionDescription string

		if err := rows.Scan(&permissionID, &permissionName, &permissionDescription); err != nil {
			return nil, err
		}

		// Add permission detail with name and description
		permissions[permissionID] = entities.PermissionDetail{
			Name:        permissionName,
			Description: permissionDescription,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &entities.Permissions{PermissionsMap: permissions}, nil
}

func (usr *User) ResetPassword(ctx context.Context, request entities.ResetPassword) error {
	query := `UPDATE users SET password = $1 WHERE userid = $2`

	_, err := usr.db.ExecContext(ctx, query, request.HashedPassword, request.UserID)
	if err != nil {
		return fmt.Errorf("failed to update password for user_id %s: %w", request.UserID, err)
	}

	return nil
}

func (usr *User) RevokeToken(ctx context.Context, tokenID string) error {
	query := `UPDATE tokens SET is_revoked = TRUE WHERE id = $1`

	_, err := usr.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke token %s: %w", tokenID, err)
	}

	return nil
}

func (usr *User) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	query := `SELECT is_revoked FROM tokens WHERE id = $1`

	var isRevoked bool
	err := usr.db.QueryRowContext(ctx, query, tokenID).Scan(&isRevoked)
	if err != nil {
		return false, fmt.Errorf("failed to check revocation status for token %s: %w", tokenID, err)
	}

	return isRevoked, nil
}

func (usr *User) ListUserPermissions(
	ctx context.Context,
	userID string,
	projectID string,
) ([]entities.PermissionDetail, error) {

	query := `
		SELECT p.id AS permission_id, p.name AS permission_name
		FROM user_permissions u
		JOIN permissions p ON p.id = ANY(u.permission_ids)
		WHERE u.user_id = $1
		AND u.project_id = $2
	`

	// Initialize a slice to store the permissions
	var permissions []entities.PermissionDetail

	// Execute the query with userID and projectID
	rows, err := usr.db.QueryContext(ctx, query, userID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and populate the slice
	for rows.Next() {
		var (
			permissionID   string
			permissionName string
		)

		// Scan the permission ID and permission name from the query result
		if err := rows.Scan(&permissionID, &permissionName); err != nil {
			return nil, err
		}

		// Append the permission detail to the slice
		permissions = append(permissions, entities.PermissionDetail{
			ID:   permissionID,
			Name: permissionName,
		})
	}

	// Check if there was an error during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of permissions
	return permissions, nil
}
