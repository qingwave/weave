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
    user_id BIGINT NOT NULL REFERENCES users(id),
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
    creator_id BIGINT,
    updater_id BIGINT,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

INSERT INTO groups (name, kind, describe, created_at) VALUES
    ('root', 'system', 'weave system group', LOCALTIMESTAMP),
    ('system:authenticated', 'system', 'system group contains all authenticated user', LOCALTIMESTAMP),
    ('system:unauthenticated', 'system', 'system group contains all unauthenticated user', LOCALTIMESTAMP)  ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS user_groups(
    group_id BIGINT NOT NULL REFERENCES groups(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    PRIMARY KEY(group_id, user_id)
);

INSERT INTO user_groups (group_id, user_id)
    SELECT  g.id, u.id FROM users AS u, groups AS g 
    WHERE (u.name='admin' AND g.name='root') ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS resources (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(256) NOT NULL,
    scope varchar(100),
    kind varchar(100)
);

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name varchar(100) NOT NULL UNIQUE,
    scope varchar(100),
    namespace varchar(100)
);

INSERT INTO roles (name, scope) VALUES
    ('cluster-admin', 'cluster'),
    ('authenticated', 'cluster'),
    ('unauthenticated', 'cluster');

CREATE TABLE IF NOT EXISTS rules (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    role_id bigint REFERENCES roles(id) on UPDATE CASCADE on DELETE CASCADE,
    resource varchar(100),
    operation varchar(100)
);

CREATE UNIQUE INDEX idx_role_rule
ON rules(role_id, resource, operation);

INSERT INTO rules (role_id, resource, operation) VALUES
  ((SELECT id FROM roles WHERE name = 'cluster-admin'), '*', '*'),
  ((SELECT id FROM roles WHERE name = 'authenticated'), 'users', '*'),
  ((SELECT id FROM roles WHERE name = 'authenticated'), 'auth', '*'),
  ((SELECT id FROM roles WHERE name = 'unauthenticated'), 'auth', 'create');

CREATE TABLE IF NOT EXISTS user_roles(
    user_id BIGINT NOT NULL REFERENCES users(id),
    role_id BIGINT NOT NULL REFERENCES roles(id),
    PRIMARY KEY(user_id, role_id)
);

CREATE TABLE IF NOT EXISTS group_roles(
    group_id BIGINT NOT NULL REFERENCES groups(id),
    role_id BIGINT NOT NULL REFERENCES roles(id),
    PRIMARY KEY(group_id, role_id)
);

INSERT INTO group_roles (group_id, role_id) VALUES
    ((SELECT id FROM groups WHERE name = 'root'), (SELECT id FROM roles WHERE name = 'cluster-admin')),
    ((SELECT id FROM groups WHERE name = 'system:authenticated'), (SELECT id FROM roles WHERE name = 'authenticated')),
    ((SELECT id FROM groups WHERE name = 'system:unauthenticated'), (SELECT id FROM roles WHERE name = 'unauthenticated'));
