package main

/*
All SQL statements in one convenient location.
*/
const (
	SportGet = "SELECT " +
		"id, name " +
		"FROM sports " +
		"WHERE id=?"

	SportGetAll = "SELECT " +
		"id, name " +
		"FROM sports"

	PositionGet = "SELECT " +
		"id, name " +
		"FROM positions " +
		"WHERE id=?"

	PositionGetAll = "SELECT " +
		"id, name " +
		"FROM positions " +
		"WHERE id=?"

	LeagueGet = "SELECT " +
		"id, name, canonical, icon, visible, official, metric, " +
		"city, country, location " +
		"FROM leagues " +
		"WHERE id=?"

	LeagueGetByCanonical = "SELECT " +
		"id, name, canonical, icon, visible, official, metric, " +
		"city, country, location " +
		"FROM leagues " +
		"WHERE canonical=?"

	LeagueGetAll = "SELECT " +
		"id, name, icon, location " +
		"FROM leagues " +
		"ORDER BY name ASC"

	LeagueGetLatest = "SELECT " +
		"id, name, canonical, icon, visible, official, metric, " +
		"city, country, location " +
		"FROM leagues " +
		"WHERE user_id="

	LeagueCreate = "INSERT INTO leagues(" +
		"name, canonical, sport_id) " +
		"VALUES ($1, $2, $3)"

	LeagueUpdate = "UPDATE " +
		"leagues " +
		"SET name=?, canonical=? " +
		"WHERE id=?"

	LeagueAdminCreate = "INSERT INTO league_admins(" +
		"league_id, user_id) " +
		"VALUES ($1, $2) "

	LeagueAdminGet = "SELECT " +
		"league_id, user_id " +
		"FROM league_admins " +
		"WHERE league_id=? and user_id=?"

	LeagueAdminGetAll = "SELECT " +
		"users.id, users.email, league_admins.league_id " +
		"FROM league_admins, users " +
		"WHERE league_id=? and league_admins.user_id=users.id"

	LeagueAdminDelete = "DELETE " +
	  "FROM league_admins " +
		"WHERE user_id=? and league_id=?"

	LeagueAdminCheck = "SELECT " +
		"user_id, league_id " +
		"FROM league_admins " +
		"WHERE league_id=? and user_id=?"

	LeagueIconUpdate = "UPDATE " +
	  "leagues " +
		"SET icon=? " +
		"WHERE id=?"

	SeasonGet = "SELECT " +
		"id, periods, duration, league_id " +
		"FROM seasons " +
		"WHERE id=?"

	SeasonGetAll = "SELECT " +
		"id, periods, duration, league_id " +
		"FROM seasons " +
		"WHERE league_id=?"

	SeasonCreate = "INSERT INTO seasons( " +
		"league_id, periods, duration) " +
		"VALUES ($1, $2, $3)"

	SeasonUpdate = "UPDATE seasons SET " +
		"periods=?, duration=? " +
		"WHERE id=?"

	SeasonGetLatest = "SELECT " +
		"MAX(id), periods, duration, league_id " +
		"FROM seasons " +
		"WHERE league_id=?"

	TeamScheduleGetAll = "SELECT " +
		"games.id, games.home_id, games.away_id, games.league_id, games.season_id, " +
		"games.completed, games.scheduled2, teams.Name " +
		"FROM games, teams " +
		"WHERE games.league_id=? and games.season_id=? and " +
		"(teams.id=games.home_id or teams.id=games.away_id) and " +
		"(games.home_id=? or games.away_id=?) and not teams.id=?" +
		"ORDER BY games.scheduled2 DESC"

	TeamGet = "SELECT " +
		"teams.id, teams.name, teams.canonical, teams.icon, teams.league_id, " +
		"leagues.canonical, leagues.name " +
		"FROM teams, leagues " +
		"WHERE leagues.id=teams.league_id and teams.id=?"

	TeamExists = "SELECT " +
	  "id, name, canonical " +
		"FROM teams " +
		"WHERE league_id=? and canonical=?"

	TeamGetByCanonical = "SELECT " +
		"id, name, canonical, icon, league_id " +
		"FROM teams " +
		"WHERE canonical=?"

	TeamGetAll = "SELECT " +
		"id, name, canonical, icon, league_id " +
		"FROM teams " +
		"WHERE league_id=? " +
		"ORDER BY name ASC"

	TeamGetAllNames = "SELECT " +
		"name, canonical " +
		"FROM teams " +
		"WHERE league_id=? " +
		"ORDER BY name ASC"

	TeamGetAllBut = "SELECT " +
		"id, name, canonical, icon, league_id " +
		"FROM teams " +
		"WHERE league_id=? and not id=?" +
		"ORDER BY name ASC"

	TeamCreate = "INSERT INTO teams(" +
		"name, canonical, icon, league_id) " +
		"VALUES ($1, $2, $3, $4)"

	TeamUpdate = "UPDATE " +
		"teams " +
		"SET name=?, canonical=? " +
		"WHERE id=?"

	TeamIconUpdate = "UPDATE " +
	  "teams " +
		"SET icon=? " +
		"WHERE id=?"

	TeamPlayerGet = "SELECT " +
		"id, first, middle, last, canonical, playerNumber, " +
		"position_id, team_id, league_id " +
		"FROM players " +
		"WHERE team_id=? and canonical=?"

	TeamPlayersGetAll = "SELECT " +
		"id, first, middle, last, canonical, playerNumber, " +
		"position_id, team_id, league_id " +
		"FROM players " +
		"WHERE league_id=? and team_id=?"

	TeamPlayerCreate = "INSERT INTO players(" +
		"first, middle, last, canonical, playerNumber," +
		" position_id, league_id, team_id) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	TeamAdminCreate = "INSERT INTO team_admins(" +
		"league_id, team_id, user_id) " +
		"VALUES ($1, $2, $3) "

	TeamAdminGet = "SELECT " +
		"team_admins.id, team_admins.league_id, team_admins.team_id, " +
		"team_admins.user_id, teams.name " +
		"FROM team_admins, teams " +
		"WHERE team_admins.team_id=? and team_admins.team_id=teams.id"

	TeamAdminGetAll = "SELECT " +
		"team_admins.id, team_admins.league_id, team_admins.team_id, " +
		"team_admins.user_id, teams.name, users.email " +
		"FROM team_admins, teams, users " +
		"WHERE team_admins.team_id=teams.id and team_admins.league_id=? " +
		"and team_admins.user_id=users.id"

	TeamAdminCheck = "SELECT " +
		"user_id, team_id " +
		"FROM team_admins " +
		"WHERE team_id=? and user_id=?"

	PlayerIconUpdate = "UPDATE " +
	  "players " +
		"SET icon=? " +
		"WHERE id=?"

	PlayerGet = "SELECT " +
		"id, first, middle, last, canonical, playerNumber, " +
		"position_id, team_id, league_id " +
		"FROM players " +
		"WHERE id=?"

	PlayerGetAll = "SELECT " +
		"players.id, players.first, players.middle, players.last, " +
		"players.canonical, players.playerNumber, players.position_id, " +
		"players.team_id, players.league_id, teams.name " +
		"FROM players, teams " +
		"WHERE players.team_id=teams.id and players.league_id=? " +
		"ORDER BY players.last ASC"

	PlayerUpdate = "UPDATE players " +
		"set first=?, middle=?, last=?, position_id=?, playerNumber=? " +
		"WHERE id=?"

	PlayerCreate = "INSERT INTO players(" +
		"first, middle, last, canonical, playerNumber," +
		" position_id, league_id) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"

	UserGet = "SELECT " +
		"id, icon, defaultLeague " +
		"FROM users " +
		"WHERE id=?"

	UserGetByEmail = "SELECT " +
		"id, name, email, icon " +
		"FROM users " +
		"WHERE email=?"

	UserCreate = "INSERT INTO users(" +
		"email, password, salt, token) " +
		"VALUES ($1, $2, $3, $4)"

	UserCreateUnregistered = "INSERT INTO users(" +
		"email) " +
		"VALUES ($1)"

	UserUpdate = "UPDATE users(" +
		"name, email, icon, city, country, location) " +
		"VALUES($1, $2, $3, $4, $5, $6) " +
		"WHERE id=?"

	UserUpdateToken = "UPDATE users " +
		"set token=? " +
		"WHERE email=?"

	UserUpdateDefaultLeague = "UPDATE users " +
		"set defaultLeague=? " +
		"WHERE id=?"

	UserLeagueGet = "SELECT " +
		"leagues.id, leagues.name, leagues.canonical, leagues.icon " +
		"FROM leagues, league_admins " +
		"WHERE league_admins.league_id=leagues.id and " +
		"league_admins.league_id=? and league_admins.user_id=?"

	UserLeagueGetAll = "SELECT " +
		"leagues.id, leagues.name, leagues.canonical, leagues.icon " +
		"FROM leagues, league_admins " +
		"WHERE league_admins.league_id=leagues.id and league_admins.user_id=? " +
		"ORDER BY leagues.name ASC"

	UserLeagueFollowGetAll = "SELECT " +
		"leagues.id, leagues.name, leagues.canonical, leagues.icon " +
		"FROM leagues, league_followers " +
		"WHERE league_followers.league_id=leagues.id and league_followers.user_id=?" +
		"ORDER BY leagues.name ASC"

	UserAllLeagueGetAll = "SELECT " +
		"leagues.id, leagues.name, leagues.canonical, leagues.icon " +
		"FROM leagues, league_followers, league_admins " +
		"WHERE league_followers.league_id=leagues.id and league_followers.user_id=? " +
		"or league_admins.league_id=leagues.id and league_admins.user_id=?" +
		"ORDER BY leagues.name ASC"

	UserLeagueGetLatest = "SELECT " +
		"MAX(leagues.id), leagues.name, leagues.canonical, leagues.icon " +
		"FROM leagues, league_admins " +
		"WHERE league_admins.league_id=leagues.id and league_admins.user_id=?"

	UserTeamGet = "SELECT " +
		"teams.id, teams.name, teams.canonical, teams.icon " +
		"FROM teams, team_admins, users " +
		"WHERE team_admins.team_id=teams.id and " +
		"team_admins.team_id=? and team_admins.user_id=?"

	BetaUserCreate = "INSERT INTO betausers(" +
		"email) " +
		"VALUES ($1)"

	ScheduleGet = "SELECT " +
		"id, home_id, away_id, league_id, season_id, completed " +
		"FROM games " +
		"WHERE id=?"

	ScheduleGetAll = "SELECT " +
		"id, scheduled2, home_id, away_id, league_id, season_id, completed " +
		"FROM games " +
		"WHERE league_id=? and season_id=? and scheduled2=? " +
		"ORDER BY scheduled2 ASC"

	ScheduleCreate = "INSERT INTO games(" +
		"home_id, away_id, scheduled2, season_id, league_id)" +
		"VALUES($1, $2, $3, $4, $5)"

	ScheduleUpdate = "UPDATE games(" +
		"scheduled2, home, away) " +
		"VALUES('true', $2, $3) " +
		"WHERE id=?"

	ScheduleFinal = "UPDATE games " +
		"SET completed='true'" +
		"WHERE league_id=? and id=?"

	ScheduleDelete = "DELETE FROM " +
	  "games " +
		"WHERE league_id=? and id=?"

	StandingsGet = "SELECT " +
		"id, team_id, wins, losses, home_wins, home_losses, away_wins, " +
		"away_losses, last_ten_wins, last_ten_losses, streak_wins, streak_losses " +
		"FROM standings " +
		"WHERE league_id=? and season_id=?"

  LeagueFollowerCreate = "INSERT INTO league_followers(" +
		"league_id, user_id)" +
		"VALUES($1, $2)"

	LeagueFollowerGet = "SELECT " +
	  "id, league_id, user_id " +
	  "FROM league_followers " +
		"WHERE league_id=? and user_id=?"

	LeagueFollowerGetAll = "SELECT " +
	  "id, league_id, user_id " +
	  "FROM league_followers " +
		"WHERE league_id=?"

	LeagueFollowerDelete = "DELETE " +
		"FROM league_followers " +
		"WHERE league_id=? and user_id=?"

	ForgotCreate = "INSERT INTO forgot(" +
	  "user_id, token) " +
		"VALUES($1, $2)"

	ForgotGet = "SELECT " +
	  "id, user_id, token " +
		"FROM forgot " +
		"WHERE user_id=?"

	ForgotDelete = "DELETE " +
	  "FROM forgot " +
		"WHERE user_id=? and token=?"

	ForgotUpdate = "UPDATE forgot " +
	  "SET token=? " +
		"WHERE user_id=?"

	ForgotExists = "SELECT " +
	  "id, user_id, token " +
		"FROM forgot " +
		"WHERE token=?"

	UserUpdatePassword = "UPDATE users " +
	  "SET password=? " +
		"WHERE id=?"

)
