CREATE INDEX IF NOT EXISTS edToys_title_idx ON edToys USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS edToys_genres_idx ON edToys USING GIN (genres);