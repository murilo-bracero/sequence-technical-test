-- In real production apps the users password should be created by the DBAs or DevOps teams and not stored in the source code.
CREATE USER sequenceapi WITH PASSWORD 'sequenceapi-password';

GRANT CONNECT ON DATABASE sequencemailbox TO sequenceapi;