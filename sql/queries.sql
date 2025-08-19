SELECT n.neighborhood_name,
       c.city_name,
       cc.category_name
FROM crime_incidents_2022 ci
         JOIN addresses a ON ci.address_id = a.address_id
         JOIN cities c ON a.city_id = c.city_id
         JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
         JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id;


SELECT crime_category_id
FROM crime_incidents_2022
GROUP BY crime_category_id
HAVING COUNT(*) > 50;


SELECT ci.*
FROM crime_incidents_2025 ci
WHERE NOT EXISTS (SELECT 1
                  FROM crime_incidents_2025 ci2
                  WHERE ci2.address_id = ci.address_id
                    AND ci2.incident_id <> ci.incident_id);


SELECT *
FROM crime_incidents_2021 ci1
         FULL OUTER JOIN addresses a ON ci1.address_id = a.address_id;


SELECT incident_id
FROM crime_incidents_partition
EXCEPT
SELECT incident_id
FROM crime_incidents_2024;


SELECT n.neighborhood_name, COUNT(*) AS total_incidents
FROM crime_incidents_partition ci
         JOIN addresses a ON ci.address_id = a.address_id
         JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
GROUP BY n.neighborhood_name
HAVING COUNT(*) > 100;


SELECT COUNT(*) AS total_incidents_2025
FROM crime_incidents_2025;


SELECT c.city_name, COUNT(DISTINCT n.neighborhood_id) AS neighborhood_count
FROM neighborhoods n
         JOIN addresses a ON n.neighborhood_id = a.neighborhood_id
         JOIN cities c ON a.city_id = c.city_id
GROUP BY c.city_name
HAVING COUNT(DISTINCT n.neighborhood_id) > 5;


SELECT c.city_name, cc.category_name, COUNT(*) AS total
FROM crime_incidents_partition ci
         JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
         JOIN addresses a ON ci.address_id = a.address_id
         JOIN cities c ON a.city_id = c.city_id
GROUP BY c.city_name, cc.category_name;


SELECT neighborhood_name, COUNT(*) AS incident_count
FROM crime_incidents_2022 ci
         JOIN addresses a ON ci.address_id = a.address_id
         JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
GROUP BY n.neighborhood_name
ORDER BY incident_count DESC
LIMIT 3;



