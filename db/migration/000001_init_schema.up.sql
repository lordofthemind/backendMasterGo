-- Create account table
CREATE TABLE account (
  id bigserial PRIMARY KEY,
  owner varchar,
  balance bigint,
  currency varchar,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

-- Create entries table
CREATE TABLE entries (
  id bigserial PRIMARY KEY,
  account_id bigint NOT NULL,
  amount bigint NOT NULL,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES account (id)
);

-- Create transfers table
CREATE TABLE transfers (
  id bigserial PRIMARY KEY,
  from_account_id bigint NOT NULL,
  to_account_id bigint NOT NULL,
  amount bigint NOT NULL,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (from_account_id) REFERENCES account (id),
  FOREIGN KEY (to_account_id) REFERENCES account (id),
  CHECK (amount > 0) -- Ensure amount is positive
);

-- Add indexes
CREATE INDEX ON account (owner);
CREATE INDEX ON entries (account_id);
CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON transfers (from_account_id, to_account_id);

-- Add comments
COMMENT ON COLUMN entries.amount IS 'Can be positive or negative';
COMMENT ON COLUMN transfers.amount IS 'Must be positive';
