CREATE TABLE IF NOT EXISTS notifications (
     id SERIAL PRIMARY KEY,
     sender VARCHAR(255),
     recipient VARCHAR(255),
     category VARCHAR(100),
     title VARCHAR(255),
     body TEXT,
     template TEXT,
     status VARCHAR(50),
     created_at TIMESTAMP,
     updated_at TIMESTAMP,
     deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS fcm_tokens (
    id SERIAL PRIMARY KEY,
    token text NOT NULL,
    user_id int(11) NOT NULL,
    os varchar(20) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);