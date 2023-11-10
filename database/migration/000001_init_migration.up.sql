CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "email" varchar,
  "password" varchar,
  "name" varchar,
  "phone_number" varchar,
  "identity" varchar,
  "address" text,
  "role" integer,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone
);

CREATE TABLE "banks" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "logo" varchar,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone
);

CREATE TABLE "bank_accounts" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "bank_id" integer,
  "status" integer,
  "account_owner" varchar,
  "sheba" varchar,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone
);

CREATE TABLE "gateways" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "commission_id" integer,
  "bank_account_id" integer,
  "name" varchar,
  "logo" varchar,
  "route" varchar,
  "status" integer,
  "type" integer,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone
);

CREATE TABLE "commissions" (
  "id" integer PRIMARY KEY,
  "amount_per_trans" float,
  "percentage_per_trans" float,
  "status" integer,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone
);

CREATE TABLE "transactions" (
  "id" integer PRIMARY KEY,
  "gateways_id" integer,
  "payment_amount" float,
  "status" integer,
  "owner_bank_account" varchar,
  "purchaser_bank_account" varchar,
  "created_at" timestamp with time zone,
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone
);

ALTER TABLE "gateways" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "gateways" ADD FOREIGN KEY ("commission_id") REFERENCES "commissions" ("id");

ALTER TABLE "gateways" ADD FOREIGN KEY ("bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("gateways_id") REFERENCES "gateways" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("bank_id") REFERENCES "banks" ("id");
