CREATE TABLE IF NOT EXISTS edToys (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    target_age VARCHAR(20) NOT NULL,
    genres text[] NOT NULL,
    skill_focus text[] NOT NULL,
    runtime INT NOT NULL
);
