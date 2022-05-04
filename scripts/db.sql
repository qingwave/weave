CREATE DATABASE weave;

\c weave;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(50) NOT NULL,
    email varchar(256) NOT NULL,
    password varchar(256),
    avatar varchar(256),
    auth_type varchar(256) NOT NULL DEFAULT 'nil',
    auth_id varchar(256),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    constraint name_auth_type unique(name, auth_type)
);

INSERT INTO users (name, email, password) VALUES ('admin', 'admin@weave.com', '$2a$10$5whQjJqSdL18PrEP.z/gZOubMKhFB38K0CvHWdnaQodb/H3yeG4J2')
