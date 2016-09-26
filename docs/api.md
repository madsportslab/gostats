# Auth
---

| API | |
| --- | --- |
| URL | PUT /auth |
| Description | Sign in, redirect to /home if cookie exists |
| Auth | No |
| Request Parameters | email, password |
| Errors | 401, 409 |

| API | |
| --- | --- |
| URL | DELETE /auth |
| Description | Sign out, remove cookie |
| Auth | No |
| Request Parameters | |
| Errors | 401, 409 |

# Users
---

| API | |
| --- | --- |
| URL | POST /users |
| Description | Create new user |
| Auth | No |
| Request Parameters | email, password |
| Errors | 401, 409 |

| API | |
| --- | --- |
| URL | PUT /users/:user |
| Description | Update user |
| Auth | Yes |
| Request Parameters | name, email, password, mobile, icon, city, country, location, defaultLeague |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /users/:user |
| Description | Get user |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /users |
| Description | Get all users |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

# Sports
---

| API | |
| --- | --- |
| URL | GET /sports |
| Description | Get all sports |
| Auth | No |
| Request Parameters ||
| Errors | 401, 409 |


## Positions

| API | |
| --- | --- |
| URL | GET /sports/:sport/positions |
| Description | Get all positions for sport |
| Auth | No |
| Request Parameters | |
| Errors | 401, 404, 409 |

# Leagues
---

| API | |
| --- | --- |
| URL | GET /leagues |
| Description | Get all leagues |
| Auth | No |
| Request Parameters | city, country, sport, query |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | POST /leagues |
| Description | Create new league |
| Auth | Yes |
| Request Parameters | name, sport |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /leagues/:league |
| Description | Get league |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | PUT /leagues/:league |
| Description | Update league |
| Auth | Yes |
| Request Parameters | name, sport, description, icon, visible, official, metric, city, country, location |
| Errors | 401, 404, 409 |

## Seasons

| API | |
| --- | --- |
| URL | GET /leagues/:league/seasons |
| Description | Get all seasons for a league |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | POST /leagues/:league/seasons |
| Description | Create season for a league |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /leagues/:league/seasons/:season |
| Description | Get season |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | PUT /leagues/:league/seasons/:season |
| Description | Update season |
| Auth | Yes |
| Request Parameters | periods, duration, fouls, start, end |
| Errors | 401, 404, 409 |

### Games

| API | |
| --- | --- |
| URL | POST /leagues/:league/seasons/:season/games |
| Description | Create game |
| Auth | Yes |
| Request Parameters | home, away, scheduled |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /leagues/:league/seasons/:season/games |
| Description | Get all games |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /leagues/:league/seasons/:season/games/:game |
| Description | Get game |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | PUT /leagues/:league/seasons/:season/games/:game |
| Description | Create player |
| Auth | Yes |
| Request Parameters | home, away, scheduled |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | DELETE /leagues/:league/seasons/:season/games/:game |
| Description | Delete game |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

## Teams

| API | |
| --- | --- |
| URL | GET /leagues/:league/teams |
| Description | Get teams |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | POST /leagues/:league/teams |
| Description | Create team |
| Auth | Yes |
| Request Parameters | name, description |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /leagues/:league/teams/:team |
| Description | Get team |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | PUT /leagues/:league/teams/:team |
| Description | Update team |
| Auth | Yes |
| Request Parameters | name, description |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | DELETE /leagues/:league/teams/:team |
| Description | Delete team |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | POST /leagues/:league/teams/:team/players |
| Description | Create player |
| Auth | Yes |
| Request Parameters | name, height, weight, jersey, handed |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | PUT /leagues/:league/teams/:team/players/:player |
| Description | Create player |
| Auth | Yes |
| Request Parameters | name, height, weight, jersey, handed, icon |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | DELETE /leagues/:league/teams/:team/players/:player |
| Description | Delete player |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

## players

| API | |
| --- | --- |
| URL | GET /leagues/:league/players |
| Description | Get players |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |

| API | |
| --- | --- |
| URL | GET /leagues/:league/players/:player |
| Description | Get player |
| Auth | Yes |
| Request Parameters | |
| Errors | 401, 404, 409 |
