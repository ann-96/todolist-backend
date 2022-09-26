ALTER TABLE todos 
  DROP CONSTRAINT todos_users_fkey;

ALTER TABLE todos
  DROP COLUMN userid;

DROP TABLE users;
