-- name: GetProfilesFoldersCatalogs :many
SELECT 
    profiles.user_id, 
    profiles.key AS profile_key, 
    folders.key AS folder_key, 
    catalogs.key AS catalog_key
FROM 
    profiles
INNER JOIN 
    folders ON profiles.key = folders.profile_key
INNER JOIN 
    catalogs ON folders.key = catalogs.folder_key
WHERE 
    profiles.user_id = $1;

-- name: GetProfilesCount :one
SELECT 
    COUNT(*) AS count
FROM
    profiles;

-- name: GetProfilesWithOffset :many
SELECT 
    user_id 
FROM
    profiles
LIMIT $1 OFFSET $2;
