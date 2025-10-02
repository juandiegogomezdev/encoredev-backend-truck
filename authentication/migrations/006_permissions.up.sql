--- Permissions in the system ---
CREATE TABLE IF NOT EXISTS permissions(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
)