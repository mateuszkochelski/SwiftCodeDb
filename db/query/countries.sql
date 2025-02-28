-- name: CreateCountry :one
INSERT INTO countries (
    country_code,
    country_name
) VALUES (
    $1, $2
) ON CONFLICT (country_code) DO NOTHING 
RETURNING *;

-- name: GetCountry :one
SELECT * FROM countries
WHERE country_code = $1;

