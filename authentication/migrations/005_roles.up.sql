--- Roles in the system ---
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
)