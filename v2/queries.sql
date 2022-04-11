SELECT 
    DISTINCT ON (s.name) name, 
    s.stat_id, s.profile_id, s.value, s.unit, s.date_created, s.last_modified, s.dataset_id, s.dataset_name 
FROM 
    key_stats_history s;
    
    
SELECT 
    DISTINCT s.name, s.date_created, ) name, date_created,
    s.stat_id, s.profile_id, s.value, s.unit, s.last_modified, s.dataset_id, s.dataset_name 
FROM 
    key_stats_history s
WHERE 
    s.date_created <= '2022-04-07 15:01:46.123856'
ORDER BY 
    s.date_created;


SELECT * FROM key_stats_history s WHERE s.date_created <= '2022-04-07 15:01:46.123856' ORDER BY s.date_created DESC;

SELECT DISTINCT s.date_created FROM key_stats_history s WHERE s.profile_id = $1 ORDER BY s.date_created DESC;


## Get a all key stats belonging to a specific version
SELECT DISTINCT ON (s.name)
    s.stat_id, 
    s.profile_id,
    s.name,
    s.value, 
    s.unit, 
    s.date_created, 
    s.last_modified, 
    s.dataset_id, 
    s.dataset_name 
FROM 
    key_stats_history s
WHERE
    s.date_created <= '2022-04-07 16:18:05.6306'
ORDER BY 
    s.name, s.date_created DESC;
