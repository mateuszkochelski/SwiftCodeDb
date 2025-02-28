-- name: CreateBank :one
INSERT INTO banks (
    swift_code,
    bank_name,
    bank_address,
    country_code,
    bank_type
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBankBySwiftCodeWithCountry :one
SELECT b.swift_code, b.bank_name, b.bank_address, b.country_code, c.country_name, b.bank_type FROM banks as b 
INNER JOIN countries as c ON b.country_code = c.country_code
WHERE swift_code = $1 LIMIT 1;

-- name: GetBanksBranchesBySwiftCodePrefix :many
SELECT b.swift_code, b.bank_name, b.bank_address, b.country_code, b.bank_type FROM banks as b 
WHERE swift_code like $1 AND swift_code != $2;

-- name: GetBanksByCountryCode :many
SELECT b.swift_code, b.bank_name, b.bank_address, b.country_code, b.bank_type FROM banks as b
WHERE b.country_code = $1;

-- name: DeleteBankBySwiftCode :exec
DELETE FROM banks
WHERE $1 = swift_code;




