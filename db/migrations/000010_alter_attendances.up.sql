ALTER TABLE attendances
   ADD COLUMN created_by_level VARCHAR(255) NOT NULL DEFAULT 'student';