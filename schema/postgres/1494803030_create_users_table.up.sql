CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
  id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
  email      VARCHAR(500) NOT NULL,
  username   VARCHAR(32)  NOT NULL,
  name       VARCHAR(100) NOT NULL,
  password   TEXT         NOT NULL,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX users_email_uniq_idx
  ON users (email);
CREATE UNIQUE INDEX users_username_uniq_idx
  ON users (username);
