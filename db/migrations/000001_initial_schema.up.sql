CREATE TABLE IF NOT EXISTS teachers
(
    id         SERIAL PRIMARY KEY,
    teacher_id VARCHAR(50)  NOT NULL UNIQUE,
    first_name VARCHAR(50)  NOT NULL,
    last_name  VARCHAR(50)  NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    phone      VARCHAR(20)  NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS classes
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(100) NOT NULL,
    homeroom_teacher VARCHAR(50)  NOT NULL REFERENCES teachers (teacher_id),
    description      TEXT         NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS students
(
    id         SERIAL PRIMARY KEY,
    student_id VARCHAR(50)  NOT NULL UNIQUE,
    classes_id INTEGER      NOT NULL REFERENCES classes (id),
    first_name VARCHAR(50)  NOT NULL,
    last_name  VARCHAR(50)  NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    phone      VARCHAR(20)  NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS attendances
(
    id          SERIAL PRIMARY KEY,
    student_id  VARCHAR(50) NOT NULL REFERENCES students (student_id),
    class_id    INTEGER     NOT NULL REFERENCES classes (id),
    date        DATE        NOT NULL,
    status      VARCHAR(20) NOT NULL CHECK (status IN ('present', 'absent', 'late')),
    description TEXT        NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS absent_requests
(
    id           SERIAL PRIMARY KEY,
    student_id   VARCHAR(50) NOT NULL REFERENCES students (student_id),
    class_id     INTEGER NOT NULL REFERENCES classes (id),
    request_date DATE        NOT NULL,
    reason       TEXT        NOT NULL,
    status       VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);