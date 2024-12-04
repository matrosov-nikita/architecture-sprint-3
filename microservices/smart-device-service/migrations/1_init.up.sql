CREATE TABLE devices (
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   serial_number TEXT,
   status TEXT
);