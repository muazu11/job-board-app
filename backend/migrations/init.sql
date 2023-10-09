CREATE TABLE advertisements (
    id int PRIMARY KEY,
    title string,
    description string,
    company_id int REFERENCES companies(id),
    wage float,
    adress string,
    zip_code string,
    city string,
    work_time interval,
);

CREATE TABLE users (
    id int PRIMARY KEY,
    email string UNIQUE,
    name string,
    surname string,
    phone string,
    date_of_birth date,
);

CREATE TABLE account (
    user_id int PRIMARY KEY REFERENCES users(id),
    password_hash string,
    role role,
);   

CREATE TABLE companies (
    id int PRIMARY KEY,
    name string,
    logo_url string,
);

CREATE TABLE applications (
    advertisement_id int REFERENCES advertisements(id),
    applicant_id int REFERENCES users(id),
    message text,
    created_at timestampz,
);

CREATE TYPE role AS ENUM (
    'admin',
    'user',
);
