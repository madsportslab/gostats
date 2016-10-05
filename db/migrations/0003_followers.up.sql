create table if not exists league_followers(
  id INTEGER NOT NULL PRIMARY KEY,
  league_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(league_id) REFERENCES leagues(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

create table if not exists team_followers(
  id INTEGER NOT NULL PRIMARY KEY,
  team_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(team_id) REFERENCES teams(id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);

create table if not exists player_followers(
  id INTEGER NOT NULL PRIMARY KEY,
  player_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  created DATETIME DEFAULT CURRENT_TIMESTAMP,
  modified DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(player_id) REFERENCES players(id),
  FOREIGN KEY(user_id) REFERENCES users(id),
  UNIQUE(player_id, user_id)
);
