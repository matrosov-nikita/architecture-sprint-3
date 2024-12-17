CREATE TABLE devices_events (
   id SERIAL PRIMARY KEY,
   device_id INT,
   event_type TEXT,
   data JSONB,
   occured_on TIMESTAMP WITH TIME ZONE
);