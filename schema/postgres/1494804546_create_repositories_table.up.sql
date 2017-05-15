CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE repositories (
  id             UUID PRIMARY KEY                           DEFAULT gen_random_uuid(),
  name           VARCHAR(100) NOT NULL,
  description    TEXT,
  website        TEXT,
  default_branch TEXT         NOT NULL,
  private        BOOLEAN                                    DEFAULT TRUE,
  bare           BOOLEAN                                    DEFAULT TRUE,
  created_at     TIMESTAMPTZ  NOT NULL                      DEFAULT now(),
  updated_at     TIMESTAMPTZ  NOT NULL                      DEFAULT now(),
  owner_id       UUID REFERENCES users NOT NULL
);
