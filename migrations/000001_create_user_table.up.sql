-- Create the ENUM type for user purpose
CREATE TYPE purpose_enum AS ENUM ('student', 'developer', 'other');

-- Create the users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    userid UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Generates a UUID for each user
    name VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL, -- Assuming the hashed password can be stored in 255 characters
    mailid VARCHAR(255) UNIQUE NOT NULL,
    purpose purpose_enum NOT NULL, -- ENUM field for the purpose
    organization VARCHAR(255),

    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically sets the timestamp when a user is created
    updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically sets the timestamp when a user is updated

    is_active BOOLEAN DEFAULT TRUE, -- Field to indicate if the user account is active
    is_deleted BOOLEAN DEFAULT FALSE, -- Field for soft-deletion status
    deleted_on TIMESTAMP -- Nullable field for soft deletion; null if not deleted
);
