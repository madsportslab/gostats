create table if not exists betausers(
  id INTEGER NOT NULL PRIMARY KEY,
  email VARCHAR NOT NULL UNIQUE,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP
);

create table if not exists users(
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR,
  email VARCHAR NOT NULL UNIQUE,
  mobile VARCHAR,
  password VARCHAR,
  salt VARCHAR,
  icon VARCHAR,
  token VARCHAR,
  city VARCHAR,
  country VARCHAR,
  location VARCHAR,
  defaultLeague INTEGER DEFAULT 0,
  registered BOOL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP
);

create table if not exists sports(
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR,
  description VARCHAR,
  icon VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP
);

create table if not exists positions(
  id INTEGER NOT NULL PRIMARY KEY,
  sport_id INTEGER NOT NULL,
  name VARCHAR,
  short VARCHAR,
  canonical VARCHAR,
  description VARCHAR,
  icon VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(sport_id) REFERENCES sports(id)
);

create table if not exists leagues(
  id INTEGER NOT NULL PRIMARY KEY,
  sport_id INTEGER,
  name VARCHAR NOT NULL UNIQUE,
  canonical VARCHAR UNIQUE,
  icon VARCHAR,
  visible BOOL DEFAULT true,
  official BOOL DEFAULT false,
  metric BOOL DEFAULT false,
  city VARCHAR,
  country VARCHAR,
  location VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(sport_id) REFERENCES sports(id)
);

create table if not exists league_admins(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER,
  user_id INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

create table if not exists seasons(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER,
  games INTEGER,
  periods INTEGER,
  duration INTEGER,
  start DATETIME,
  finish DATETIME,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id)
);

create table if not exists conferences(
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR NOT NULL,
  canonical VARCHAR,
  icon VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP
);

create table if not exists teams(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER,
  name VARCHAR NOT NULL,
  canonical VARCHAR,
  icon VARCHAR,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id)
);

create table if not exists team_admins(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER,
  team_id INTEGER,
  user_id INTEGER,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id),
  FOREIGN KEY(team_id) REFERENCES teams(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

create table if not exists players(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER,
  team_id INTEGER,
  position_id INTEGER,
  first VARCHAR,
  middle VARCHAR,
  last VARCHAR,
  canonical VARCHAR,
  icon VARCHAR,
  height FLOAT,
  weight FLOAT,
  hand INTEGER,
  playerNumber VARCHAR,
  birth DATE,
  handed INTEGER DEFAULT 0,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id),
  FOREIGN KEY(team_id) REFERENCES teams(id),
  FOREIGN KEY(position_id) REFERENCES positions(id)
);

create table if not exists games(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER,
  season_id INTEGER,
  home_id INTEGER,
  away_id INTEGER,
  scheduled DATE,
  completed BOOL DEFAULT false,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id),
  FOREIGN KEY(season_id) REFERENCES seasons(id),
  FOREIGN KEY(home_id) REFERENCES teams(id),
  FOREIGN KEY(away_id) REFERENCES teams(id)
);
