CREATE DATABASE weave;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(50) UNIQUE NOT NULL,
    email varchar(256) UNIQUE NOT NULL,
    password varchar(256) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);
