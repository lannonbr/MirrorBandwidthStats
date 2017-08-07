-- schema for sqlite tables.
CREATE TABLE month (id integer PRIMARY KEY, time text NOT NULL, rx integer NOT NULL, tx integer NOT NULL, rate real NOT NULL);
CREATE TABLE day (id integer PRIMARY KEY, time text NOT NULL, rx integer NOT NULL, tx integer NOT NULL, rate real NOT NULL);
CREATE TABLE hour (id integer PRIMARY KEY, time text NOT NULL, rx integer NOT NULL, tx integer NOT NULL, rate real NOT NULL);
