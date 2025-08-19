CREATE TYPE provider AS ENUM ('google', 'github');
CREATE TYPE role AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS public.locations
(
    location_id bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    latitude    numeric(10, 7)                      NOT NULL,
    longitude   numeric(10, 7)                      NOT NULL,
    CHECK (latitude >= -90.0 AND latitude <= 90.0),
    CHECK (longitude >= -180.0 AND longitude <= 180.0),
    UNIQUE (latitude, longitude)
);
CREATE TABLE IF NOT EXISTS public.states
(
    state_id   bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    state_name varchar(100)                        NOT NULL,
    UNIQUE (state_name)
);

CREATE TABLE IF NOT EXISTS public.counties
(
    county_id   bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    state_id    bigint                              NOT NULL,
    county_name varchar(100)                        NOT NULL,
    FOREIGN KEY (state_id) REFERENCES public.states (state_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    UNIQUE (state_id, county_name)
);

CREATE TABLE IF NOT EXISTS public.cities
(
    city_id      bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    city_name    varchar(100)                        NOT NULL,
    county_id    bigint                              NOT NULL,
    population   integer,
    area_sq_mi   numeric(10, 2),
    founded_year integer,
    is_active    boolean                             NOT NULL DEFAULT TRUE,
    created_at   timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (county_id) REFERENCES public.counties (county_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CHECK (population > 0 OR population IS NULL),
    CHECK (area_sq_mi > 0 OR area_sq_mi IS NULL),
    CHECK (founded_year > 0 OR founded_year IS NULL),
    UNIQUE (city_name, county_id)
);

CREATE TABLE IF NOT EXISTS public.neighborhoods
(
    neighborhood_id   bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    neighborhood_name varchar(100)                        NOT NULL,
--     city_id           bigint                              NOT NULL,
    created_at        timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP
--     FOREIGN KEY (city_id) REFERENCES public.cities (city_id) ON UPDATE CASCADE ON DELETE RESTRICT,
--     UNIQUE (neighborhood_name, city_id)
);

CREATE TABLE IF NOT EXISTS public.addresses
(
    address_id      bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    street_address  varchar(255),
    city_id         bigint                              NOT NULL,
    postal_code     varchar(20),
    neighborhood_id bigint,
    FOREIGN KEY (city_id) REFERENCES public.cities (city_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (neighborhood_id) REFERENCES public.neighborhoods (neighborhood_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS public.data_sources
(
    source_id   bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    source_name varchar(255)                        NOT NULL UNIQUE,
    source_url  varchar(500),
    is_active   boolean                             NOT NULL DEFAULT TRUE,
    created_at  timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.crime_categories
(
    crime_category_id bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    category_name     varchar(100)                        NOT NULL,
    created_at        timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (LENGTH(category_name) >= 3)
);


CREATE TABLE IF NOT EXISTS public.crime_incidents
(
    incident_id       bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    address_id        bigint,
    crime_category_id bigint                              NOT NULL,
    incident_date     date                                NOT NULL,
    incident_time     time without time zone,
    location_id       bigint,
    case_num          varchar(25),
    is_resolved       boolean                             NOT NULL DEFAULT FALSE,
    source_id         bigint                              NULL,
    created_at        timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (address_id) REFERENCES public.addresses (address_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (crime_category_id) REFERENCES public.crime_categories (crime_category_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (location_id) REFERENCES public.locations (location_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (source_id) REFERENCES data_sources (source_id) ON UPDATE CASCADE,
    UNIQUE (address_id, case_num)
);

CREATE TABLE IF NOT EXISTS crime_incidents_partition
(
    incident_id       bigint GENERATED BY DEFAULT AS IDENTITY,
    address_id        bigint,
    crime_category_id bigint                                 NOT NULL,
    incident_date     date                                   NOT NULL,
    incident_time     time,
    location_id       bigint,
    case_num          varchar(25),
    is_resolved       boolean                                NOT NULL,
    source_id         bigint,
    created_at        timestamp with time zone DEFAULT NOW() NOT NULL,
    PRIMARY KEY (incident_id, incident_date),
    UNIQUE (case_num, incident_date, address_id),
    FOREIGN KEY (address_id) REFERENCES addresses
        ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (crime_category_id) REFERENCES crime_categories
        ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (location_id) REFERENCES locations
        ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (source_id) REFERENCES data_sources
        ON UPDATE CASCADE ON DELETE RESTRICT
)
    PARTITION BY RANGE (incident_date);

CREATE TABLE crime_incidents_2018
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2018-01-01') TO ('2019-01-01');

CREATE TABLE crime_incidents_2019
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2019-01-01') TO ('2020-01-01');

CREATE TABLE crime_incidents_2020
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2020-01-01') TO ('2021-01-01');

CREATE TABLE crime_incidents_2021
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2021-01-01') TO ('2022-01-01');

CREATE TABLE crime_incidents_2022
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2022-01-01') TO ('2023-01-01');

CREATE TABLE crime_incidents_2023
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');

CREATE TABLE crime_incidents_2024
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

CREATE TABLE crime_incidents_2025
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');


CREATE TABLE IF NOT EXISTS public.user_profiles
(
    user_id    bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    username   varchar(50)                         NOT NULL UNIQUE,
    email      varchar(255)                        NOT NULL UNIQUE,
    role       role                                NOT NULL DEFAULT 'user',
    is_active  boolean                             NOT NULL DEFAULT TRUE,
    created_at timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE TABLE IF NOT EXISTS public.password_login
(
    user_id       bigint                   NOT NULL PRIMARY KEY,
    password_hash varchar(255)             NOT NULL,
    access_token  varchar(500)             NOT NULL,
    refresh_token varchar(500),
    created_at    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id),
    FOREIGN KEY (user_id) REFERENCES public.user_profiles (user_id) ON UPDATE CASCADE ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS public.oauth_login
(
    user_id          bigint                   NOT NULL,
    provider         provider                 NOT NULL,
    provider_user_id varchar(100)             NOT NULL,
    access_token     varchar(500)             NOT NULL,
    refresh_token    varchar(500),
    created_at       timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, provider),
    FOREIGN KEY (user_id) REFERENCES public.user_profiles (user_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    UNIQUE (user_id),
    UNIQUE (provider, provider_user_id),
    UNIQUE (access_token)
);


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




CREATE OR REPLACE FUNCTION add_crime_incident_partition(
    p_city_name varchar(100),
    p_state_name varchar(100),
    p_source_name varchar(255),
    p_latitude numeric(10, 7),
    p_longitude numeric(10, 7),
    p_street_address varchar(255),
    p_postal_code varchar(20),
    p_neighborhood varchar(100),
    p_crime_category_name varchar(100),
    p_incident_date date,
    p_incident_time time without time zone DEFAULT NULL,
    p_case_num varchar(25) DEFAULT NULL,
    p_is_resolved boolean DEFAULT FALSE
) RETURNS bigint AS
$$
DECLARE
    v_state_id          bigint;
    v_county_id         bigint;
    v_city_id           bigint;
    v_location_id       bigint;
    v_neighborhood_id   bigint;
    v_source_id         bigint;
    v_address_id        bigint;
    v_crime_category_id bigint;
    v_incident_id       bigint;
BEGIN
    SELECT state_id
    INTO v_state_id
    FROM public.states
    WHERE state_name = p_state_name;

    IF v_state_id IS NULL THEN
        INSERT INTO public.states (state_name)
        VALUES (p_state_name)
        RETURNING state_id INTO v_state_id;
    END IF;

    SELECT c.city_id, c.county_id
    INTO v_city_id, v_county_id
    FROM public.cities c
             JOIN public.counties co ON c.county_id = co.county_id
    WHERE co.state_id = v_state_id
      AND c.city_name = p_city_name
    LIMIT 1;

    IF v_city_id IS NULL THEN
        SELECT county_id
        INTO v_county_id
        FROM public.counties
        WHERE state_id = v_state_id
        LIMIT 1;

        IF v_county_id IS NULL THEN
            INSERT INTO public.counties (state_id, county_name)
            VALUES (v_state_id, 'Unknown County')
            RETURNING county_id INTO v_county_id;
        END IF;

        INSERT INTO public.cities (city_name, county_id)
        VALUES (p_city_name, v_county_id)
        RETURNING city_id INTO v_city_id;
    END IF;

    SELECT neighborhood_id
    INTO v_neighborhood_id
    FROM public.neighborhoods
    WHERE neighborhood_name = p_neighborhood;

    IF v_neighborhood_id IS NULL THEN
        INSERT INTO public.neighborhoods (neighborhood_name, city_id)
        VALUES (p_neighborhood, v_city_id)
        RETURNING neighborhood_id INTO v_location_id;
    END IF;

    SELECT location_id
    INTO v_location_id
    FROM public.locations
    WHERE latitude = p_latitude
      AND longitude = p_longitude;

    IF v_location_id IS NULL THEN
        INSERT INTO public.locations (latitude, longitude)
        VALUES (p_latitude, p_longitude)
        RETURNING location_id INTO v_location_id;
    END IF;

    SELECT source_id
    INTO v_source_id
    FROM public.data_sources
    WHERE source_name = p_source_name;

    IF v_source_id IS NULL THEN
        INSERT INTO public.data_sources (source_name)
        VALUES (p_source_name)
        RETURNING source_id INTO v_source_id;
    END IF;

    SELECT address_id
    INTO v_address_id
    FROM public.addresses
    WHERE city_id = v_city_id
      AND (street_address = p_street_address OR (street_address IS NULL AND p_street_address IS NULL))
      AND (postal_code = p_postal_code OR (postal_code IS NULL AND p_postal_code IS NULL));

    IF v_address_id IS NULL THEN
        INSERT INTO public.addresses (street_address, city_id, postal_code, neighborhood_id)
        VALUES (p_street_address, v_city_id, p_postal_code, v_neighborhood_id)
        RETURNING address_id INTO v_address_id;
    END IF;

    SELECT crime_category_id
    INTO v_crime_category_id
    FROM public.crime_categories
    WHERE category_name = p_crime_category_name;

    IF v_crime_category_id IS NULL THEN
        INSERT INTO public.crime_categories (category_name)
        VALUES (p_crime_category_name)
        RETURNING crime_category_id INTO v_crime_category_id;
    END IF;

    INSERT INTO public.crime_incidents_partition (address_id,
                                                  crime_category_id,
                                                  incident_date,
                                                  incident_time,
                                                  location_id,
                                                  case_num,
                                                  is_resolved,
                                                  source_id)
    VALUES (v_address_id,
            v_crime_category_id,
            p_incident_date,
            p_incident_time,
            v_location_id,
            p_case_num,
            p_is_resolved,
            v_source_id)
    RETURNING incident_id INTO v_incident_id;

    RETURN v_incident_id;

EXCEPTION
    WHEN OTHERS THEN
        RAISE EXCEPTION 'Error adding crime incident: %', sqlerrm;
END;
$$ LANGUAGE plpgsql;


SELECT *
FROM add_crime_incident_partition(p_city_name := 'Tacoma',
                                  p_state_name := 'Washington',
                                  p_country_name := 'United States',
                                  p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                                  p_latitude := 47.191788,
                                  p_longitude := -122.466841,
                                  p_street_address := '2300 S 72ND ST',
                                  p_postal_code := '98409',
                                  p_crime_category_name := 'Traffic Accident/Collision - Non Fatal - Non Injury',
                                  p_incident_date := '2019-03-22'::date,
                                  p_incident_time := '09:00:00'::time,
                                  p_case_num := '1908100621');

SELECT *
FROM add_crime_incident_partition(p_city_name := 'Tacoma',
                                  p_state_name := 'Washington',
                                  p_country_name := 'United States',
                                  p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                                  p_latitude := 47.260788,
                                  p_longitude := -122.4678411,
                                  p_street_address := '1500 N 10TH S',
                                  p_postal_code := '98403',
                                  p_crime_category_name := 'Fraud Offenses',
                                  p_incident_date := '2018-07-12'::date,
                                  p_incident_time := '09:00:00'::time,
                                  p_case_num := '1819300700');

SELECT *
FROM add_crime_incident_partition(p_city_name := 'Tacoma',
                                  p_state_name := 'Washington',
                                  p_country_name := 'United States',
                                  p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                                  p_latitude := 47.258788,
                                  p_longitude := -122.443841,
                                  p_street_address := '600 S BAKER ST',
                                  p_postal_code := '98402',
                                  p_crime_category_name := 'Larceny/Theft Offenses',
                                  p_incident_date := '2018-01-12'::date,
                                  p_incident_time := '16:00:00'::time,
                                  p_case_num := '1901809009');

SELECT *
FROM add_crime_incident_partition(p_city_name := 'Tacoma',
                                  p_state_name := 'Washington',
                                  p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                                  p_latitude := 47.179788,
                                  p_longitude := -122.462841,
                                  p_street_address := '8600 S HOSMER ST',
                                  p_postal_code := '98444',
                                  p_crime_category_name := 'Larceny/Theft Offenses',
                                  p_incident_date := '2018-11-14'::date,
                                  p_incident_time := '19:00:00'::time,
                                  p_case_num := '1832009015');

SELECT *
FROM add_crime_incident_partition(p_city_name := 'Tacoma',
                                  p_state_name := 'Washington',
                                  p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                                  p_latitude := 47.253788,
                                  p_longitude := -122.444841,
                                  p_street_address := '900 TACOMA AVE S',
                                  p_postal_code := '98402',
                                  p_crime_category_name := 'Destruction/Damage/Vandalism',
                                  p_incident_date := '2024-01-18'::date,
                                  p_incident_time := '06:00:00'::time,
                                  p_case_num := '2401909037');


SELECT *
FROM add_crime_incident_partition(p_city_name := 'Tacoma',
                                  p_state_name := 'Washington',
                                  p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                                  p_latitude := 47.174,
                                  p_longitude := -122.433,
                                  p_street_address := '9200 PACIFIC AVE',
                                  p_postal_code := '98444',
                                  p_crime_category_name := 'Motor Vehicle Theft',
                                  p_incident_date := '2025-04-04'::date,
                                  p_incident_time := '05:00:00'::time,
                                  p_case_num := '2509801347');
