CREATE TYPE org_type AS ENUM ('company', 'personal');


CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY,
    owner_id UUID REFERENCES users(id) NOT NULL,
    name VARCHAR(256) NOT NULL,
    type org_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);