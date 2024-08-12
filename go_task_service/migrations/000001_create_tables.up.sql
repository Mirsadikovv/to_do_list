CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done', 'caceled');

CREATE SEQUENCE task_external_id_seq START WITH 1;

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    external_id VARCHAR NOT NULL,
    title VARCHAR(35),
    task_status task_status DEFAULT 'todo',
    task_description VARCHAR(300),
    deadline TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
