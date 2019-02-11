package migration

var _0 = `
# entity set
CREATE TABLE IF NOT EXISTS host (
  id          SERIAL,
  fqdn        VARCHAR(63) NOT NULL,
  port        SMALLINT UNSIGNED,
  private_key BLOB,
  public_key  BLOB,

  PRIMARY KEY (id),
  INDEX USING BTREE (fqdn)
);

# entity set
CREATE TABLE IF NOT EXISTS user (
  username    VARCHAR(30) UNIQUE NOT NULL,
  first_name  VARCHAR(30),
  last_name   VARCHAR(30),
  admin       BOOLEAN,
  private_key BLOB,
  public_key  BLOB,

  PRIMARY KEY (username)
);

# entity set
CREATE TABLE IF NOT EXISTS command (
  id        SERIAL,
  cmd       TEXT NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  exit_code SMALLINT,

  PRIMARY KEY (id)
);

# multi-valued attribute
CREATE TABLE IF NOT EXISTS command_stdout (
  id        SERIAL,
  timestamp TIMESTAMP NOT NULL,
  line      BLOB NOT NULL,

  FOREIGN KEY (id) REFERENCES command(id) ON DELETE CASCADE,
  PRIMARY KEY (id, timestamp)
);

# multi-valued attribute
CREATE TABLE IF NOT EXISTS command_stderr (
  id        SERIAL,
  timestamp TIMESTAMP NOT NULL,
  line      BLOB NOT NULL,

  FOREIGN KEY (id) REFERENCES command(id) ON DELETE CASCADE,
  PRIMARY KEY (id, timestamp)
);

# relationship set, many-to-one
CREATE TABLE IF NOT EXISTS command_user (
  id       SERIAL,
  username VARCHAR(30) NOT NULL,

  FOREIGN KEY (id)       REFERENCES command(id)    ON DELETE CASCADE,
  FOREIGN KEY (username) REFERENCES user(username) ON DELETE CASCADE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS command_host (
  id      SERIAL,
  host_id BIGINT UNSIGNED NOT NULL,

  FOREIGN KEY (id)      REFERENCES command(id) ON DELETE CASCADE,
  FOREIGN KEY (host_id) REFERENCES host(id)    ON DELETE CASCADE,
  PRIMARY KEY (id)
);
`
