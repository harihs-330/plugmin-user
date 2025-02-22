-- Revert the token_type column to exclude 'password_reset'
ALTER TABLE tokens 
DROP CONSTRAINT IF EXISTS tokens_token_type_check;

ALTER TABLE tokens 
ADD CONSTRAINT tokens_token_type_check 
CHECK (token_type IN ('access', 'refresh'));
