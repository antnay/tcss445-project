INSERT INTO countries (country_name)
VALUES ('United States');

WITH country_query AS (SELECT country_id
                       FROM countries
                       WHERE country_name = 'United States')
INSERT
INTO states (country_id, state_name)
SELECT c.country_id,
       state_data.state_name
FROM country_query c
         CROSS JOIN
     (VALUES ('Alabama'),
             ('Alaska'),
             ('Arizona'),
             ('Arkansas'),
             ('California'),
             ('Colorado'),
             ('Connecticut'),
             ('Delaware'),
             ('Florida'),
             ('Georgia'),
             ('Hawaii'),
             ('Idaho'),
             ('Illinois'),
             ('Indiana'),
             ('Iowa'),
             ('Kansas'),
             ('Kentucky'),
             ('Louisiana'),
             ('Maine'),
             ('Maryland'),
             ('Massachusetts'),
             ('Michigan'),
             ('Minnesota'),
             ('Mississippi'),
             ('Missouri'),
             ('Montana'),
             ('Nebraska'),
             ('Nevada'),
             ('New Hampshire'),
             ('New Jersey'),
             ('New Mexico'),
             ('New York'),
             ('North Carolina'),
             ('North Dakota'),
             ('Ohio'),
             ('Oklahoma'),
             ('Oregon'),
             ('Pennsylvania'),
             ('Rhode Island'),
             ('South Carolina'),
             ('South Dakota'),
             ('Tennessee'),
             ('Texas'),
             ('Utah'),
             ('Vermont'),
             ('Virginia'),
             ('Washington'),
             ('West Virginia'),
             ('Wisconsin'),
             ('Wyoming'))
         AS state_data (state_name);

WITH state_query AS (SELECT state_id
                     FROM states
                     WHERE state_name = 'Washington')
INSERT
INTO counties (state_id, county_name)
SELECT s.state_id,
       county_data.county_name
FROM state_query s
         CROSS JOIN (VALUES ('Adams'),
                            ('Asotin'),
                            ('Benton'),
                            ('Chelan'),
                            ('Clallam'),
                            ('Clark'),
                            ('Columbia'),
                            ('Cowlitz'),
                            ('Douglas'),
                            ('Ferry'),
                            ('Franklin'),
                            ('Garfield'),
                            ('Grant'),
                            ('Grays Harbor'),
                            ('Island'),
                            ('Jefferson'),
                            ('King'),
                            ('Kitsap'),
                            ('Kittitas'),
                            ('Klicktitat'),
                            ('Lewis'),
                            ('Lincoln'),
                            ('Mason'),
                            ('Okanogan'),
                            ('Pacific'),
                            ('Pend Oreille'),
                            ('Pierce'),
                            ('San Juan'),
                            ('Skagit'),
                            ('Skamania'),
                            ('Snohomish'),
                            ('Spokane'),
                            ('Stevens'),
                            ('Thurston'),
                            ('Wahkiakum'),
                            ('Walla Walla'),
                            ('Whatcom'),
                            ('Whitman'),
                            ('Yakima'))
    AS county_data(county_name);


WITH state_query AS (SELECT state_id
                     FROM states
                     WHERE state_name = 'Washington'),
     county_insert AS (SELECT county_id
                       FROM counties
                       WHERE county_name = 'Pierce')
INSERT
INTO cities (city_name, county_id, population, area_sq_mi, founded_year)
SELECT city_data.city_name,
       c.county_id,
       city_data.population,
       city_data.area_sq_mi,
       city_data.founded_year
FROM state_query s
         CROSS JOIN county_insert c
         CROSS JOIN (VALUES ('Tacoma', 222906, 62.63, 1875),
                            ('Lakewood', 62303, 18.9, 1996),
                            ('Puyallup', 42179, 14.1, 1890),
                            ('Federal Way', 97701, 21.16, 1990),
                            ('Fife', 10723, 5.676, 1957))
    AS city_data(city_name, population, area_sq_mi, founded_year);


WITH state_query AS (SELECT state_id
                     FROM states
                     WHERE state_name = 'Washington'),
     county_insert AS (SELECT county_id
                       FROM counties
                       WHERE county_name = 'King')
INSERT
INTO cities (city_name, county_id, population, area_sq_mi, founded_year)
SELECT city_data.city_name,
       c.county_id,
       city_data.population,
       city_data.area_sq_mi,
       city_data.founded_year
FROM state_query s
         CROSS JOIN county_insert c
         CROSS JOIN (VALUES ('Seattle', 755078, 83.78, 1851),
                            ('Bothell', 50213, 12.08, 1909),
                            ('Burien', 50730, 15.68, 1993),
                            ('SeaTac', 31799, 10.12, 1990),
                            ('Bellevue', 151574, 33.9, 1953),
                            ('Tukwila', 21135, 9.6, 1908))
    AS city_data(city_name, population, area_sq_mi, founded_year);

WITH state_query AS (SELECT state_id
                     FROM states
                     WHERE state_name = 'Washington'),
     county_insert AS (SELECT county_id
                       FROM counties
                       WHERE county_name = 'Snohomish')
INSERT
INTO cities (city_name, county_id, population, area_sq_mi, founded_year)
SELECT city_data.city_name,
       c.county_id,
       city_data.population,
       city_data.area_sq_mi,
       city_data.founded_year
FROM state_query s
         CROSS JOIN county_insert c
         CROSS JOIN (VALUES ('Mill Creek', 20742, 4.691, 1983))
    AS city_data(city_name, population, area_sq_mi, founded_year);


WITH city_query AS (SELECT city_id
                    FROM cities
                    WHERE city_name = 'Tacoma')
INSERT
INTO neighborhoods (neighborhood_name)
SELECT data.neighborhood_name,
       s.city_id
FROM city_query s
         CROSS JOIN(VALUES ('South End'))
    AS data(neighborhood_name)
ON CONFLICT (neighborhood_name, city_id) DO NOTHING
RETURNING neighborhood_id;

