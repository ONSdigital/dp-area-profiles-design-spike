-- Get profile ID for area with code x.
SELECT 
    profile_id, name, area_code 
FROM 
   area_profiles 
WHERE 
   area_code = $1;


-- Get key stats for a given area profile.
SELECT 
   s.profile_id, s.stat_id, t.type_id, t.name, s.value, s.unit, s.date_created, s.dataset_id, s.dataset_name 
FROM 
   key_stats s
INNER JOIN
    key_stat_types t
ON
    t.type_id = s.stat_type
WHERE 
   s.profile_id = $1;


-- Get a lst of versions.
SELECT DISTINCT 
    s.date_created 
FROM 
    key_stats_history s 
WHERE 
    s.profile_id = $1 
ORDER BY 
    s.date_created DESC;


-- Get all key stats belonging to the specified version.
SELECT DISTINCT ON 
    (t.name) s.profile_id, s.stat_id, t.name, t.type_id, s.value, s.unit, s.date_created, s.dataset_id, s.dataset_name 
FROM 
    key_stats_history s 
INNER JOIN
    key_stat_types t
ON
    t.type_id = s.stat_type
WHERE 
    s.profile_id = $1 AND s.date_created <= $2
ORDER BY 
    t.name, s.date_created DESC;


-- Get dataset Recipes
SELECT 
    r.dataset_id,
    r.dataset_edition,
    r.cantabular_query,
    r.stat_type,
    g.name
FROM
    key_stats_recipes r
INNER JOIN
    recipe_geographies rg ON rg.recipe_id = r.recipe_id
INNER JOIN 
    geography_types g ON g.id = rg.geography_type_id
WHERE
    r.dataset_id = 'Test dataset 1' AND r.dataset_edition = 'abc123';