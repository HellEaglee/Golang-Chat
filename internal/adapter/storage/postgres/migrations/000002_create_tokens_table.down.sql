DROP INDEX IF EXISTS idx_tokens_user_id;
DROP INDEX IF EXISTS idx_tokens_token;
DROP INDEX IF EXISTS idx_tokens_expires_at;
DROP INDEX IF EXISTS idx_tokens_revoked_at;

DROP TABLE IF EXISTS tokens;