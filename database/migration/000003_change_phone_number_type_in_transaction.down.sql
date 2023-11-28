-- Step 1: Create a new column with the old type
ALTER TABLE transactions
    ADD COLUMN temp_phone varchar(10);

-- Step 2: Update the new column with the values from the new column
UPDATE transactions
SET temp_phone = phone_number::varchar(10);

-- Step 3: Drop the new column
ALTER TABLE transactions
DROP COLUMN phone_number;

-- Step 4: Rename the old column to the original name
ALTER TABLE transactions
    RENAME COLUMN temp_phone TO phone_number;