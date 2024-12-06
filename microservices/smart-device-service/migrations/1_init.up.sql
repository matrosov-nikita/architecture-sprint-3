CREATE TABLE devices (
   id SERIAL PRIMARY KEY,
   user_id TEXT,
   name TEXT NOT NULL,
   serial_number TEXT,
   status TEXT
);

-- Вставим тестовые устройства.
INSERT INTO devices (user_id, name, serial_number, status) VALUES (1,'термостат', '6R342', 'off');