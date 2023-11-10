CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "email" varchar,
  "password" varchar,
  "name" varchar,
  "phone_number" varchar,
  "identity" varchar,
  "address" text,
  "role" uint8,
  "created_at" timestamp_with_timezone,
  "updated_at" timestamp_with_timezone,
  "deleted_at" timestamp_with_timezone
);

CREATE TABLE "banks" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "logo" varchar,
  "created_at" timestamp_with_timezone,
  "updated_at" timestamp_with_timezone,
  "deleted_at" timestamp_with_timezone
);

CREATE TABLE "bank_accounts" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "bank_id" integer,
  "status" uint8,
  "account_owner" varchar,
  "sheba" varchar,
  "created_at" timestamp_with_timezone,
  "updated_at" timestamp_with_timezone,
  "deleted_at" timestamp_with_timezone
);

CREATE TABLE "payment_gateaways" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "commission_id" integer,
  "bank_account_id" integer,
  "name" varchar,
  "logo" varchar,
  "route" varchar,
  "status" uint8,
  "type" uint8,
  "created_at" timestamp_with_timezone,
  "updated_at" timestamp_with_timezone,
  "deleted_at" timestamp_with_timezone
);

CREATE TABLE "commissions" (
  "id" integer PRIMARY KEY,
  "amount_per_trans" float64,
  "percentage_per_trans" float64,
  "status" uint8,
  "created_at" timestamp_with_timezone,
  "updated_at" timestamp_with_timezone,
  "deleted_at" timestamp_with_timezone
);

CREATE TABLE "transactions" (
  "id" integer PRIMARY KEY,
  "payment_gateaways_id" integer,
  "payment_amount" float64,
  "status" uint8,
  "owner_bank_account" varchar,
  "purchaser_bank_account" varchar,
  "created_at" timestamp_with_timezone,
  "updated_at" timestamp_with_timezone,
  "deleted_at" timestamp_with_timezone
);

ALTER TABLE "payment_gateaways" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "payment_gateaways" ADD FOREIGN KEY ("commission_id") REFERENCES "commissions" ("id");

ALTER TABLE "payment_gateaways" ADD FOREIGN KEY ("bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("payment_gateaways_id") REFERENCES "payment_gateaways" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("bank_id") REFERENCES "banks" ("id");
