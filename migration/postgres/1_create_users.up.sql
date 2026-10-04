-- migration/postgres/20250216100653_init.up.sql

CREATE TABLE IF NOT EXISTS users (
                                     id              UUID PRIMARY KEY,
                                     name            TEXT NOT NULL,
                                     email           TEXT NOT NULL UNIQUE,
                                     password_hash   TEXT NOT NULL,
                                     created_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS folders (
                                       id              UUID PRIMARY KEY,
                                       name            TEXT NOT NULL,
                                       owner_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                       parent_id       UUID REFERENCES folders(id) ON DELETE CASCADE,
                                       created_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS files (
                                     id              UUID PRIMARY KEY,
                                     name            TEXT NOT NULL,
                                     path            TEXT NOT NULL,
                                     size            BIGINT NOT NULL,
                                     mime_type       TEXT NOT NULL,
                                     owner_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                     folder_id       UUID REFERENCES folders(id) ON DELETE SET NULL,
                                     created_at      TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_files_owner_id ON files(owner_id);
CREATE INDEX IF NOT EXISTS idx_files_folder_id ON files(folder_id);
CREATE INDEX IF NOT EXISTS idx_folders_owner_id ON folders(owner_id);