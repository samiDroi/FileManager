-- +goose Up

CREATE TABLE IF NOT EXISTS public."sessions"
(
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES public."user" (id) ON DELETE CASCADE,
    token_hash VARCHAR(64) UNIQUE NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE,
    ip_address VARCHAR(45),
    user_agent TEXT
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON public."sessions" (user_id);

-- +goose Down

DROP TABLE IF EXISTS public."sessions";
