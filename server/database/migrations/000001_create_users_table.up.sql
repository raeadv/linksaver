CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users(
   id UUID PRIMARY KEY DEFAULT uuidv7(),
   name VARCHAR (150) NOT NULL UNIQUE,
   username VARCHAR (50) UNIQUE NOT NULL,
   password TEXT NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- add indexes
CREATE INDEX user_id_index ON users (id);
CREATE INDEX user_username_index ON users (username);