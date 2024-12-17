CREATE TABLE devices (
   id SERIAL PRIMARY KEY,
   user_id TEXT,
   name TEXT NOT NULL,
   serial_number TEXT,
   status TEXT,
   created_at timestamptz
);

-- Вставим тестовые устройства.
INSERT INTO devices (user_id, name, serial_number, status, created_at) VALUES (1,'термостат', '6R342', 'off', now());


CREATE TABLE IF NOT EXISTS outbox_events(
     id SERIAL PRIMARY KEY,
     payload JSONB NOT NULL,
     status TEXT NOT NULL DEFAULT 'new' CHECK(status IN ('new', 'done')),
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);