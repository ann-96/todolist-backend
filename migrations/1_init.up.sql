CREATE TABLE todos (
  id        SERIAL  PRIMARY KEY,
  task      TEXT    NOT NULL,
  completed BOOLEAN NOT NULL
);
