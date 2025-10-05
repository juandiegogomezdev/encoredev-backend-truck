CREATE TYPE status_org_membership AS ENUM ('active', 'suspended', 'revoked', 'ended');

CREATE TABLE IF NOT EXISTS org_memberships (
    id UUID PRIMARY KEY,
    org_id UUID NOT NULL REFERENCES organizations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role_id UUID NOT NULL REFERENCES roles(id),
    status status_org_membership NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    finalized_by UUID REFERENCES users(id),
    finalized_at TIMESTAMP WITH TIME ZONE
)