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
  key_type    VARCHAR(30),

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
  username  VARCHAR(30) UNIQUE NOT NULL,
  host_id   BIGINT UNSIGNED NOT NULL,
  timestamp BIGINT UNSIGNED NOT NULL,
  exit_code SMALLINT,

  PRIMARY KEY (id),
  FOREIGN KEY (username) REFERENCES user(username) ON DELETE CASCADE,
  FOREIGN KEY (host_id)  REFERENCES host(id)       ON DELETE CASCADE
);

# multi-valued attribute
CREATE TABLE IF NOT EXISTS command_stdout (
  id        BIGINT UNSIGNED NOT NULL,
  timestamp BIGINT UNSIGNED NOT NULL,
  line      BLOB NOT NULL,

  FOREIGN KEY (id) REFERENCES command(id) ON DELETE CASCADE,
  PRIMARY KEY (id, timestamp)
);

# multi-valued attribute
CREATE TABLE IF NOT EXISTS command_stderr (
  id        BIGINT UNSIGNED NOT NULL,
  timestamp BIGINT UNSIGNED NOT NULL,
  line      BLOB NOT NULL,

  FOREIGN KEY (id) REFERENCES command(id) ON DELETE CASCADE,
  PRIMARY KEY (id, timestamp)
);
`
