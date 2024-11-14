CREATE TABLE IF NOT EXISTS session_types
(
  id   SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

INSERT INTO session_types(id, name)
VALUES (1, 'work'), (2, 'break')
