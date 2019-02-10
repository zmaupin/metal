package migration

var _0 = `
CREATE TABLE IF NOT EXISTS host (
  id SERIAL,
  success BOOLEAN
);
`
