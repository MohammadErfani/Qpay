ALTER TABLE transactions
    DROP COLUMN IF EXISTS tracking_code;
ALTER TABLE transactions
    DROP COLUMN IF EXISTS phone_number;
ALTER TABLE transactions
    DROP COLUMN IF EXISTS card_year;
ALTER TABLE transactions
    DROP COLUMN IF EXISTS card_month;
ALTER TABLE transactions
    DROP COLUMN IF EXISTS purchaser_card;
ALTER TABLE transactions
    DROP COLUMN IF EXISTS commission_amount;

-- ALTER TABLE transactions
-- DROP COLUMN IF EXISTS TrackingCode;
--
-- ALTER TABLE transactions
-- DROP COLUMN IF EXISTS PhoneNumber;
--
-- ALTER TABLE transactions
-- DROP COLUMN IF EXISTS CardYear;
--
-- ALTER TABLE transactions
-- DROP COLUMN IF EXISTS CardMonth;
--
-- ALTER TABLE transactions
-- DROP COLUMN IF EXISTS PurchaserCard;
--
-- ALTER TABLE transactions
-- DROP COLUMN IF EXISTS CommissionAmount;