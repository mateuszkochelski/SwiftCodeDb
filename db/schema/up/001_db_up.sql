CREATE TYPE "bank_type" AS ENUM (
  'headquarter',
  'branch'
);

CREATE TABLE "countries"(
  "country_code" VARCHAR(2) PRIMARY KEY UNIQUE,
  "country_name" TEXT NOT NULL
);

CREATE TABLE "banks" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "swift_code" varchar(11) NOT NULL,
  "bank_name" TEXT NOT NULL,
  "bank_address" TEXT,
  "country_code" varchar(2) NOT NULL,
  "bank_type" bank_type NOT NULL
);

CREATE INDEX ON "banks" ("swift_code");

CREATE INDEX ON "banks" ("country_code");

ALTER TABLE banks ADD CONSTRAINT fk_banks_country FOREIGN KEY (country_code) REFERENCES countries(country_code);

ALTER TABLE banks ADD CONSTRAINT bank_name CHECK (LENGTH(bank_name) > 0);

ALTER TABLE banks ALTER COLUMN country_code SET DATA TYPE VARCHAR(2) COLLATE "C";

ALTER TABLE banks ADD CONSTRAINT country_code CHECK (LENGTH(country_code) = 2);

ALTER TABLE banks ADD CONSTRAINT swift_code_11_letters CHECK (char_length(swift_code) = 11);

ALTER TABLE banks ADD CONSTRAINT country_code_uppercase CHECK (country_code = UPPER(country_code::TEXT));

ALTER TABLE countries ADD CONSTRAINT country_name_uppercase CHECK (UPPER(country_name) = country_name);

ALTER TABLE banks ADD CONSTRAINT swift_code_end_with_XXX_implies_headquarter_bank_type
CHECK (
    ((swift_code LIKE '%XXX') and bank_type IN ('headquarter')) 
    or (swift_code NOT LIKE '%XXX' and bank_type IN ('branch'))
);


