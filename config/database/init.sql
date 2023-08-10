
CREATE USER bootcampinterview WITH CREATEDB NOSUPERUSER INHERIT PASSWORD 'password';
CREATE DATABASE bootcampinterview OWNER bootcampinterview;

\c bootcampinterview bootcampinterview
CREATE TABLE "user_roles" (
  "id" VARCHAR(100) PRIMARY KEY,
  "name" VARCHAR(30)
);

CREATE TABLE "users" (
  "id" VARCHAR(100) PRIMARY KEY,
  "email" VARCHAR(50) UNIQUE NOT NULL,
  "username" VARCHAR(50) UNIQUE NOT NULL,
  "password" VARCHAR(100) NOT NULL,
  "role_id" VARCHAR(100),
  FOREIGN KEY ("role_id") REFERENCES "user_roles" ("id")
);

CREATE TABLE "bootcamp" (
  "id" VARCHAR(100) PRIMARY KEY,
  "name"VARCHAR(50) NOT NULL,
  "start_date" date,
  "end_date" date,
  "location" VARCHAR(100)
);

CREATE TABLE "candidate" (
  "id" VARCHAR(100) PRIMARY KEY,
  "full_name" VARCHAR(30) NOT NULL,
  "phone" VARCHAR(30),
  "email" VARCHAR(50) UNIQUE NOT NULL,
  "date_of_birth" DATE NOT NULL,
  "address" VARCHAR(100) NOT NULL,
  "cv_link" VARCHAR(255),
  "bootcamp_id" VARCHAR(100),
  "instansi_pendidikan"VARCHAR(50),
  "hackerrank_score" INT,
  FOREIGN KEY ("bootcamp_id") REFERENCES "bootcamp" ("id")
);

CREATE TABLE "status" (
  "id" VARCHAR(100) PRIMARY KEY,
  "name" VARCHAR(30) NOT NULL
);

CREATE TABLE "interviewer" (
  "id" VARCHAR(100) PRIMARY KEY,
  "full_name" VARCHAR(50) NOT NULL,
  "user_id" VARCHAR(100) NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE "hr_recruitment" (
  "id" VARCHAR(100) PRIMARY KEY,
  "full_name" VARCHAR(50) NOT NULL,
  "user_id" VARCHAR(100) NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

CREATE TABLE "interviews_process" (
  "id" VARCHAR(100) PRIMARY KEY,
  "candidate_id" VARCHAR(100) NOT NULL,
  "interviewer_id" VARCHAR(100) NOT NULL,
  "interview_datetime" timestamp NOT NULL,
  "meeting_link" VARCHAR(100),
  "form_interview" VARCHAR(100),
  "status_id" VARCHAR(100),
  FOREIGN KEY ("candidate_id") REFERENCES "candidate" ("id"),
  FOREIGN KEY ("interviewer_id") REFERENCES "interviewer" ("id"),
  FOREIGN KEY ("status_id") REFERENCES "status" ("id")
);

CREATE TABLE "result" (
  "id" VARCHAR(100) PRIMARY KEY,
  "name" VARCHAR(30)
);

CREATE TABLE "interview_result" (
  "id" VARCHAR(100) PRIMARY KEY,
  "interview_id" VARCHAR(100) NOT NULL,
  "result_id" VARCHAR(100),
  "note" VARCHAR(100),
  FOREIGN KEY ("interview_id") REFERENCES "interviews_process" ("id"),
  FOREIGN KEY ("result_id") REFERENCES "result" ("id")
);
