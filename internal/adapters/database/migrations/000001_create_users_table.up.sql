CREATE TABLE IF NOT EXISTS users (
    id         CHAR(36)     NOT NULL,
    name       VARCHAR(100) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    CONSTRAINT users_id_unique UNIQUE (id)
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
