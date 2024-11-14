CREATE TABLE IF NOT EXISTS sessions
(
  id         SERIAL PRIMARY KEY,
  duration   INT NOT NULL,
  type_id    INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  FOREIGN KEY(type_id) REFERENCES session_types(id) ON DELETE SET NULL
);
