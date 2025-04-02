BEGIN;

CREATE TABLE security_questions (
    id SERIAL PRIMARY KEY,
    question TEXT NOT NULL
);

ALTER TABLE customertable 
ADD COLUMN security_question_id INTEGER NOT NULL REFERENCES security_questions(id),
ADD COLUMN security_answer_hash TEXT NOT NULL;

INSERT INTO security_questions (question) VALUES 
('What was the name of your first pet?'),
('What was your childhood nickname?'),
('In what city or town was your first job?'),
('What is the name of your favorite childhood friend?'),
('What is your mother''s maiden name?'),
('What high school did you attend?'),
('What was the make of your first car?'),
('What is your favorite movie?'),
('What is your favorite book?'),
('What was the street you lived on in third grade?');

CREATE INDEX idx_security_question_id ON customertable(security_question_id);

COMMIT;
