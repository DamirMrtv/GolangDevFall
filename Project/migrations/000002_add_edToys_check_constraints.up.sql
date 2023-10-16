ALTER TABLE edToys ADD CONSTRAINT edToys_runtime_check CHECK (runtime >= 0);
ALTER TABLE edToys ADD CONSTRAINT edToys_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE edToys ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);