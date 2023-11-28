ALTER TABLE transactions
    ADD commission_amount float;
ALTER TABLE transactions
    ADD purchaser_card varchar(16);
ALTER TABLE transactions
    ADD card_month integer;
ALTER TABLE transactions
    ADD card_year integer;
ALTER TABLE transactions
    ADD phone_number varchar(10);
ALTER TABLE transactions
    ADD tracking_code varchar(255);

-- if use bellow migration:
-- ALTER TABLE transactions
--     ADD CommissionAmount float;
-- ALTER TABLE transactions
--     ADD PurchaserCard varchar(16);
-- ALTER TABLE transactions
--     ADD CardMonth integer;
-- ALTER TABLE transactions
--     ADD CardYear integer;
-- ALTER TABLE transactions
--     ADD PhoneNumber varchar(10);
-- ALTER TABLE transactions
--     ADD TrackingCode varchar(255);