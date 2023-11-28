ALTER TABLE transactions
    ADD COLUMN temp_phone varchar;

UPDATE transactions
SET temp_phone = phone_number::varchar;

-- Step 3: Drop the old column
ALTER TABLE transactions
DROP COLUMN phone_number;

-- Step 4: Rename the new column to the original name
ALTER TABLE transactions
    RENAME COLUMN temp_phone TO phone_number;