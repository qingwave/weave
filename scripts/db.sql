CREATE DATABASE weave;

\c weave;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(50) UNIQUE NOT NULL,
    email varchar(256) UNIQUE NOT NULL,
    password varchar(256) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

INSERT INTO users (name, email, password) VALUES ('admin', 'admin@weave.com', '$2a$10$5whQjJqSdL18PrEP.z/gZOubMKhFB38K0CvHWdnaQodb/H3yeG4J2')
