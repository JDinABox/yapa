CREATE TABLE Users (
  id          BLOB    NOT NULL PRIMARY KEY, -- UUID
  created     INTEGER NOT NULL,
  email       TEXT    NOT NULL UNIQUE,
  name        TEXT    NOT NULL
);

-- Permissions should be a static table
CREATE TABLE Permissions (
  id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name        TEXT    NOT NULL UNIQUE
);

INSERT INTO Permissions VALUES ('admin');

CREATE TABLE UserPermissions (
  id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  userID      BLOB    NOT NULL REFERENCES Users(id), -- UUID
  created     INTEGER NOT NULL,
  permission  TEXT    NOT NULL REFERENCES Permissions(name)
)

CREATE TABLE Providers (
  id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  type        TEXT    NOT NULL UNIQUE
);

INSERT INTO Providers VALUES ('apple'), ('google');

CREATE TABLE AuthMethods (
  id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  userID      BLOB    NOT NULL REFERENCES Users(id), -- UUID
  created     INTEGER NOT NULL,
  provider    TEXT    NOT NULL REFERENCES Providers(type),
  identifier  TEXT    NOT NULL UNIQUE,
  secret      TEXT    NOT NULL
);

CREATE UNIQUE INDEX AuthMethods_provider_identifier_index ON AuthMethods(provider, identifier);

CREATE TABLE Sessions (
  id          BLOB    NOT NULL PRIMARY KEY, -- UUID
  created     INTEGER NOT NULL,
  updated     INTEGER NOT NULL,
  expires     INTEGER NOT NULL,
  data        BLOB    NOT NULL
)

CREATE INDEX Sessions_expires_index ON Sessions(expires);
