CREATE DATABASE weave;

\c weave;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(100) NOT NULL UNIQUE,
    email varchar(256) NOT NULL,
    password varchar(256),
    avatar varchar(256),
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

INSERT INTO users (name, email, password, created_at) VALUES
    ('admin', 'admin@weave.com', '$2a$10$5whQjJqSdL18PrEP.z/gZOubMKhFB38K0CvHWdnaQodb/H3yeG4J2', LOCALTIMESTAMP),
    ('demo', 'admin@weave.com', '$2a$10$5whQjJqSdL18PrEP.z/gZOubMKhFB38K0CvHWdnaQodb/H3yeG4J2', LOCALTIMESTAMP);

CREATE TABLE IF NOT EXISTS auth_infos (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id BIGSERIAL NOT NULL REFERENCES users(id),
    url varchar(256),
    auth_type varchar(256),
    auth_id varchar(256),
    access_token varchar(256),
    refresh_token varchar(256),
    expiry timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(100) NOT NULL UNIQUE,
    describe varchar(1024),
    creator_id BIGSERIAL,
    updater_id BIGSERIAL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

INSERT INTO groups (name, describe, created_at) VALUES
    ('root', 'weave system group', LOCALTIMESTAMP),
    ('tenant', 'weave tenant group', LOCALTIMESTAMP);

CREATE TABLE IF NOT EXISTS user_groups(
    group_id BIGSERIAL NOT NULL REFERENCES groups(id),
    user_id BIGSERIAL NOT NULL REFERENCES users(id),
    PRIMARY KEY(group_id, user_id)
);

insert into user_groups (group_id, user_id)
    select  g.id, u.id from users as u, groups as g 
    where (u.name='admin' and g.name='root') or (u.name='demo' and g.name='tenant');
