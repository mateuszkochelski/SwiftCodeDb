CREATE TYPE "bank_type" AS ENUM (
  'headquarter',
  'branch'
);

CREATE TABLE "banks" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "swift_code" varchar(11) NOT NULL,
  "bank_name" TEXT NOT NULL,
  "bank_address" TEXT,
  "country_iso2_code" varchar(2) NOT NULL,
  "country" TEXT NOT NULL,
  "bank_type" bank_type NOT NULL
);

CREATE INDEX ON "banks" ("swift_code");

CREATE INDEX ON "banks" ("country_iso2_code");

ALTER TABLE banks ADD CONSTRAINT bank_name CHECK (LENGTH(bank_name) > 0);

ALTER TABLE banks ADD CONSTRAINT country_iso2_code CHECK (LENGTH(country_iso2_code) = 2);

ALTER TABLE banks ADD CONSTRAINT country CHECK (LENGTH(country) > 0);

ALTER TABLE banks ADD CONSTRAINT swift_code_11_letters CHECK (char_length(swift_code) = 11);

ALTER TABLE banks ADD CONSTRAINT country_iso2_code_uppercase CHECK (UPPER(country_iso2_code) = country_iso2_code);

ALTER TABLE banks ADD CONSTRAINT country_name_uppercase CHECK (UPPER(country) = country);

ALTER TABLE banks ADD CONSTRAINT swift_code_end_with_XXX_implies_headquarter_bank_type
CHECK (
    ((swift_code LIKE '%XXX') and bank_type IN ('headquarter')) or (swift_code NOT LIKE '%XXX' and bank_type IN ('branch'))
);


