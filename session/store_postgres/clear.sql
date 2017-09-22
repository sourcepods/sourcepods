-- name: clear-expired

DELETE FROM sessions
WHERE expires < now();
