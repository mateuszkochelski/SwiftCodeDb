-- name: CreateBank :one
INSERT INTO banks (
    swift_code,
    bank_name,
    bank_address,
    country_iso2_code,
    country,
    bank_type
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetBankBySwiftCode :one
SELECT * FROM banks
WHERE swift_code = $1;

-- name: GetBanksBySwiftCodePrefix :many
SELECT * FROM banks
WHERE swift_code like $1;

-- name: GetBanksByCountryISO2Code :many
SELECT * FROM banks
WHERE country_iso2_code = $1;

-- name: DeleteBankBySwiftCode :exec
DELETE FROM banks
WHERE $1 = swift_code;




