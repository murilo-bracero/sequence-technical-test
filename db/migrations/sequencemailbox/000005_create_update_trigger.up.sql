CREATE OR REPLACE FUNCTION update_timestamp_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_timestamp_trigger
BEFORE UPDATE ON sequences
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp_column();