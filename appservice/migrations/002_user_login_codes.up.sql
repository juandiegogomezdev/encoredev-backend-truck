CREATE TABLE IF NOT EXISTS user_login_codes (
    user_id UUID PRIMARY KEY REFERENCES users(id) NOT NULL,
    code char(6) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE
);