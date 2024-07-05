CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL CHECK (length(username) >= 1),
  pass TEXT NOT NULL, 
  sudo INTEGER
);

CREATE TABLE events (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL CHECK (length(name) >= 1),
  description TEXT NOT NULL,
  venue TEXT,
  date TEXT,
  kind TEXT, 
  thumbnail BLOB,
  UNIQUE (id, name)
);

CREATE TABLE concerts (
  event_id INTEGER PRIMARY KEY REFERENCES events(id),
  artist TEXT
);

CREATE TABLE games (
  event_id INTEGER NULL PRIMARY KEY REFERENCES events(id),
  team1 TEXT,
  team2 TEXT
);

CREATE TABLE tickets (
  user_id INTEGER REFERENCES users(id),
  event_id INTEGER REFERENCES events(id),
  seat TEXT,
  PRIMARY KEY (event_id, seat)
);

CREATE TABLE user_events ( 
  user_id INTEGER REFERENCES users(id),
  event_id INTEGER REFERENCES events(id),
  PRIMARY KEY (user_id, event_id)
)
