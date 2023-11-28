CREATE TABLE IF NOT EXISTS dummy (
    user_id UUID PRIMARY KEY,
    user_name TEXT NOT NULL,
    created_at timestamptz,
    updated_at timestamptz
)