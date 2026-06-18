INSERT INTO user_one_time_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

UPDATE user_one_time_tokens
SET used_at = NOW()
WHERE token_hash = $1
  AND used_at IS NULL
  AND expires_at > NOW()
RETURNING user_id;
