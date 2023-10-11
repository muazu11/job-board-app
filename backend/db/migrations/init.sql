CREATE TYPE role AS ENUM (
    'admin',
    'user'
);

CREATE TABLE companies (
    id int PRIMARY KEY,
    name text,
    logo_url text
);

CREATE TABLE users (
    id int PRIMARY KEY,
    email text UNIQUE,
    name text,
    surname text,
    phone text,
    date_of_birth date
);

CREATE TABLE account (
    user_id int PRIMARY KEY REFERENCES users(id),
    password_hash text,
    role role
);   

CREATE TABLE advertisements (
    id int PRIMARY KEY,
    title text,
    description text,
    company_id int REFERENCES companies(id),
    wage float,
    address text,
    zip_code text,
    city text,
    work_time interval
);

CREATE TABLE applications (
    advertisement_id int REFERENCES advertisements(id),
    applicant_id int REFERENCES users(id),
    message text,
    created_at timestamptz
);
