CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL CHECK (length(username) >= 1),
  pass TEXT NOT NULL,
  sudo INTEGER
);

CREATE TABLE events (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL CHECK (length(name) >= 1),
  description TEXT NOT NULL,
  venue TEXT,
  date TEXT,
  kind TEXT, 
  thumbnail BLOB
);

CREATE TABLE concerts (
  event_id INTEGER NOT NULL PRIMARY KEY REFERENCES events(id),
  artist TEXT
);

CREATE TABLE games (
  event_id INTEGER NOT NULL PRIMARY KEY REFERENCES events(id),
  team1 TEXT,
  team2 TEXT
);

CREATE TABLE tickets (
  id INTEGER NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  event_id INTEGER NOT NULL REFERENCES events(id),
  seat TEXT
);

