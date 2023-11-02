CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE
);

CREATE TABLE segments(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    percent INT
);

CREATE TABLE operations(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    segment_name VARCHAR(50) NOT NULL,
    operation_type VARCHAR(20) NOT NULL,
    operation_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE users_segments(
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    segment_id INT REFERENCES segments(id) ON DELETE CASCADE,
    timeout INT,
    CONSTRAINT users_segments_pk PRIMARY KEY(user_id,segment_id)
);