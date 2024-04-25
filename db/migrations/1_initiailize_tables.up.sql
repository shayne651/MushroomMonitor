CREATE TABLE IF NOT EXISTS mushroom(
  mushroom_uuid TEXT PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS stage(
  stage_uuid TEXT PRIMARY KEY,
  mushroom_uuid TEXT,
  name TEXT NOT NULL,
  min_humidity REAL NOT NULL,
  max_humidity REAL NOT NULL,
  min_temp REAL NOT NULL,
  max_temp REAL NOT NULL,
  fea REAL NOT  NULL,
  FOREIGN KEY (mushroom_uuid) REFERENCES mushroom(mushroom_uuid)
);

CREATE TABLE IF NOT EXISTS temp(
  temp_uuid TEXT,
  grow_name TEXT,
  value REAL,
  record_date INTEGER,
  PRIMARY KEY (temp_uuid, grow_name)
);

CREATE TABLE IF NOT EXISTS humidity(
  humidity_uuid TEXT,
  grow_name TEXT,
  value REAL,
  record_date INTEGER,
  PRIMARY KEY (humidity_uuid, grow_name)
);

CREATE TABLE IF NOT EXISTS fea(
  fea_uuid TEXT,
  grow_name TEXT,
  runtime REAL,
  record_date INTEGER,
  PRIMARY KEY (fea_uuid, grow_name)
);

CREATE TABLE IF NOT EXISTS grow (
  name TEXT PRIMARY KEY,
  automation_uuid TEXT,
  mushroom_uuid TEXT,
  stage_uuid TEXT,
  last_stage_update TEXT,
  FOREIGN KEY (mushroom_uuid) REFERENCES mushroom(mushroom_uuid),
  FOREIGN KEY (stage_uuid) REFERENCES stage(stage_uuid)
);
