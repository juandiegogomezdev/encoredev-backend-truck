CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID REFERENCES roles(id) NOT NULL,
    permissions_id UUID REFERENCES permissions(id) NOT NULL,
    PRIMARY KEY (role_id, permissions_id)
)