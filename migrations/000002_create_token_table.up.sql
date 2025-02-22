CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Generates a UUID for each token
    user_id UUID REFERENCES users(userid) ON DELETE CASCADE, -- Foreign key to the users table
    token TEXT NOT NULL, -- Stores the token string
    token_type VARCHAR(20) NOT NULL CHECK (token_type IN ('access', 'refresh')), -- Token type enum
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically sets the timestamp when a token is created
    expires_on TIMESTAMP NOT NULL, -- Explicitly set expiration time
    is_revoked BOOLEAN DEFAULT FALSE -- Indicates whether the token is revoked, defaults to FALSE
);
