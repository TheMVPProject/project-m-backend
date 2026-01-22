-- Create the user table
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    position INT DEFAULT Null
);

-- Create an index on the email column for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- create a trigger to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    New.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Create the forgot_password_sessions table
CREATE TABLE IF NOT EXISTS forgot_password_sessions(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    time_initiated TIMESTAMPTZ NOT NULL,
    time_last_request TIMESTAMPTZ NOT NULL,
    otp VARCHAR(10) NOT NULL,
    attempt_count INT DEFAULT 0,
    is_verified BOOLEAN DEFAULT FALSE,
    varification_token VARCHAR(255),
    is_password_updated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_fpt_user_id ON forgot_password_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_fpt_session_token ON forgot_password_sessions(session_token);

-- create a trigger to automatically update the updated_at timestamp
CREATE TRIGGER set_timestamp_forgot_password_sessions
BEFORE UPDATE ON forgot_password_sessions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();