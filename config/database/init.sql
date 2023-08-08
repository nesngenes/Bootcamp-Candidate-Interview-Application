CREATE TABLE candidate (
    candidate_id VARCHAR(100) PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    address TEXT,
    date_of_birth DATE,
    application_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interviewer (
    interviewer_id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) REFERENCES user(user_id),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100),
    phone VARCHAR(20),
    specialization VARCHAR(100)
);

CREATE TABLE hr_recruitment (
    hr_id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) REFERENCES user(user_id),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100),
    phone VARCHAR(20)
);

CREATE TABLE bootcamp (
    bootcamp_id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100),
    start_date DATE,
    end_date DATE,
    location VARCHAR(200)
);

CREATE TABLE interview_process (
    process_id VARCHAR(100) PRIMARY KEY,
    candidate_id VARCHAR(100) REFERENCES candidate(candidate_id),
    interviewer_id VARCHAR(100) REFERENCES interviewer(interviewer_id),
    hr_id VARCHAR(100) REFERENCES hr_recruitment(hr_id),
    resume_id VARCHAR(100) REFERENCES resume(resume_id),
    bootcamp_id VARCHAR(100) REFERENCES bootcamp(bootcamp_id),
    interview_date TIMESTAMP,
    result VARCHAR(50)
);

CREATE TABLE resume (
	resume_id VARCHAR(100) PRIMARY KEY,
    candidate_id VARCHAR(100) REFERENCES candidate(candidate_id),
	cv_file VARCHAR(100)
);

CREATE TABLE user_role (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    password_hash TEXT,
    role_id INT REFERENCES user_role(role_id)
);