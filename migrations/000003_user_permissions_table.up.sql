-- Migration up file to create the user_permissions table

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_permissions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    project_id uuid NOT NULL,
    permission_ids uuid[] NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp DEFAULT now(),
    updated_at timestamp DEFAULT now()
);
