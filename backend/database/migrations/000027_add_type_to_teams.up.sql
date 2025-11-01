ALTER TABLE teams ADD COLUMN type VARCHAR(50) CHECK (type IN ('squad', 'chapter', 'tribe', 'guild'));
UPDATE teams SET type = 'squad' WHERE type IS NULL;

