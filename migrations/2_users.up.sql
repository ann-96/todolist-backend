CREATE TABLE users (
  id           SERIAL  PRIMARY KEY,
  login        TEXT    UNIQUE NOT NULL,
  passwordhash TEXT    NOT NULL
);

ALTER TABLE todos
  ADD userid INTEGER NOT NULL DEFAULT(0);

ALTER TABLE todos
  ADD CONSTRAINT todos_users_fkey
    FOREIGN KEY (userid) REFERENCES users(id);
