CREATE TYPE role AS ENUM (
    'admin',
    'user'
);

CREATE TABLE companies (
    id serial PRIMARY KEY,
    name text CHECK (length(name) > 0),
    siren text UNIQUE CHECK (siren ~ '^\d{9}$'),
    logo_url text
);

CREATE TABLE users (
    id serial PRIMARY KEY,
    email text UNIQUE CHECK (email ~ '^[a-zA-Z0-9\-\.]+@[a-zA-Z0-9\-\.]+\.[a-zA-Z0-9\-\.]+$'),
    name text CHECK (length(name) > 0),
    surname text CHECK (length(surname) > 0),
    phone text CHECK (phone IS NULL OR phone ~ '^(\+\d{2}|0)[1-9]\d{8}$'),
    date_of_birth date CHECK (date_of_birth < now())
);

CREATE TABLE accounts (
    user_id int PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash text,
    auth_token uuid DEFAULT gen_random_uuid(),
    role role
);   

CREATE TABLE advertisements (
    id serial PRIMARY KEY,
    title text CHECK (length(title) > 1),
    description text,
    company_id int REFERENCES companies(id) ON DELETE CASCADE,
    wage float,
    address text,
    zip_code text,
    city text,
    work_time interval
);

CREATE TABLE applications (
    id serial PRIMARY KEY,
    advertisement_id int REFERENCES advertisements(id) ON DELETE CASCADE,
    applicant_id int REFERENCES users(id) ON DELETE CASCADE,
    message text,
    created_at timestamptz DEFAULT now(),
    UNIQUE (advertisement_id, applicant_id)
);
