
CREATE SCHEMA IF NOT EXISTS hearts;

CREATE TABLE IF NOT EXISTS hearts.event (
  type varchar(20),
  version varchar(20),
  created_on timestamp,
  aggregate_id uuid,
  aggregate_version int,
  payload jsonb,
  PRIMARY KEY (aggregate_id, aggregate_version)
);
