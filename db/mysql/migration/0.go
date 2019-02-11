package migration

var _0 = `
CREATE TABLE IF NOT EXISTS host (
  id SERIAL,
  fqdn VARCHAR(63) NOT NULL,
  private_key BLOB,
  public_key BLOB,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user (
  username VARCHAR(30) NOT NULL UNIQUE,
  first_name VARCHAR(30),
  last_name VARCHAR(30),
  admin BOOLEAN,
  private_key BLOB,
  public_key BLOB,
  PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS command (
  id SERIAL,
  cmd TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS command_stdout (
  id SERIAL,
  line BLOB NOT NULL,
  FOREIGN KEY id REFERENCES command(id) ON DELETE CASCADE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS command_stderr (
  id SERIAL,
  line BLOB NOT NULL,
  FOREIGN KEY id REFERENCES command (id) ON DELETE CASCADE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_command (
  id SERIAL,
  username VARCHAR(30) NOT NULL,
  FOREIGN KEY id REFERENCES command (id) ON DELETE CASCADE,
  FOREIGN KEY username REFERENCES user (username) ON DELETE CASCADE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS host_command (
  id SERIAL,
  host_id BIGINT NOT NULL,
  FOREIGN KEY id REFERENCES command (id) ON DELETE CASCADE,
  FOREIGN KEY host_id REFERENCES host (id) ON DELETE CASCADE,
  PRIMARY_KEY (id)
);
`
