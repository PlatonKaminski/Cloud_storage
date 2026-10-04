-- migration/postgres/20250216100653_init.down.sql

DROP INDEX IF EXISTS idx_folders_owner_id;
DROP INDEX IF EXISTS idx_files_folder_id;
DROP INDEX IF EXISTS idx_files_owner_id;
DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS users;