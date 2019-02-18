package migration

var _0 = `
# entity set
CREATE TABLE IF NOT EXISTS instance (
  # FQDN:PORT
  net VARCHAR(67),

  PRIMARY KEY (net)
);

# entity set
CREATE TABLE IF NOT EXISTS host (
  id          SERIAL,
  fqdn        VARCHAR(63) NOT NULL UNIQUE,
  port        VARCHAR(5),
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

  PRIMARY KEY (username)
);

# entity set
CREATE TABLE IF NOT EXISTS command (
  id        SERIAL,
  cmd       TEXT NOT NULL,
  username  VARCHAR(30) UNIQUE NOT NULL,
  host_id   BIGINT UNSIGNED NOT NULL,
  timestamp BIGINT UNSIGNED NOT NULL,
  stdout    LONGBLOB,
  stderr    LONGBLOB,
  exit_code SMALLINT,

  PRIMARY KEY (id),
  FOREIGN KEY (username) REFERENCES user(username) ON DELETE CASCADE,
  FOREIGN KEY (host_id)  REFERENCES host(id)       ON DELETE CASCADE
);
`
