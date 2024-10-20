CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  hash_password TEXT NOT NULL
);

CREATE TABLE dreams (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  description TEXT NOT NULL,
  explanation TEXT NOT NULL,
  date TEXT NOT NULL
)
