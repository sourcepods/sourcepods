CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE sessions (
  id       UUID PRIMARY KEY                 DEFAULT gen_random_uuid(),
  expires  TIMESTAMPTZ           NOT NULL   DEFAULT now(),
  owner_id UUID REFERENCES users NOT NULL
);
