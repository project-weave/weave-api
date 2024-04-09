CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS events (
    id UUID DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    event_dates DATE[] NOT NULL,
    is_specific_dates BOOLEAN NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),      

    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS event_responses (
    user_id UUID,
    event_id UUID,
    alias TEXT DEFAULT '',
    availabilities TIMESTAMP[][] NOT NULL,

    PRIMARY KEY (user_id, event_id, alias),

    CONSTRAINT fk_user_id   
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_event_id  
        FOREIGN KEY (event_id) 
        REFERENCES events(id)
        ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION guest_uuid()
    RETURNS UUID 
    LANGUAGE plpgsql IMMUTABLE
    AS 
$$
BEGIN 
    RETURN '00000000-0000-0000-0000-000000000000'::UUID;
END;
$$;

INSERT INTO users(id, first_name, last_name, email) 
VALUES(guest_uuid(), '', '', '');