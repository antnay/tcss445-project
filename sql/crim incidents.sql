INSERT INTO public.crime_categories (category_name)
VALUES ('Fraud Offenses'),
       ('Traffic Accident/Collision - Non Fatal - Non Injury'),
       ('Larceny/Theft Offenses'),
       ('Destruction/Damage/Vandalism');


-- WITH loc_insert AS (INSERT INTO locations (latitude, longitude)
--     VALUES (47.260788, -122.467841),
--            (47.191788, -122.466841),
--            (47.258788, -122.443841)
--     RETURNING location_id),
--      city_query AS (SELECT city_id
--                     FROM cities
--                     WHERE city_name = 'Tacoma'),
--      address_insertion AS (INSERT
--          INTO addresses (street_address, city_id, postal_code)
--              SELECT data.street_address,
--                     c.city_id,
--                     data.postal_code
--              FROM city_query c
--                       CROSS JOIN (VALUES ('1500 N 10TH ST', 98403),
--                                          ('2300 S 72ND ST', 98402),
--                                          ('600 S BAKER ST', 98402))
--                  AS data(street_address, postal_code)
--              RETURNING address_id)
-- INSERT
-- INTO crime_incidents (address_id, crime_category_id, incident_date, incident_time, location_id, case_num, is_resolved,
--                       source_id, created_at)
-- SELECT a.address_id,
--     FROM city_query as c,
--          address_insertion as a,
--          loc_insert as l,
--
-- ;


CREATE OR REPLACE FUNCTION add_crime_incident(
    p_city_name varchar(100),
    p_state_name varchar(100),
    p_country_name varchar(50),
    p_source_name varchar(255),
    p_latitude numeric(10, 7),
    p_longitude numeric(10, 7),
    p_street_address varchar(255),
    p_postal_code varchar(20),
    p_crime_category_name varchar(100),
    p_incident_date date,
    p_incident_time time without time zone DEFAULT NULL,
    p_case_num varchar(25) DEFAULT NULL,
    p_is_resolved boolean DEFAULT FALSE
) RETURNS bigint AS
$$
DECLARE
    v_country_id        bigint;
    v_state_id          bigint;
    v_county_id         bigint;
    v_city_id           bigint;
    v_location_id       bigint;
    v_source_id         bigint;
    v_address_id        bigint;
    v_crime_category_id bigint;
    v_incident_id       bigint;
BEGIN
    -- 1. Handle country (lookup first, insert if not found)
    SELECT country_id
    INTO v_country_id
    FROM public.countries
    WHERE country_name = p_country_name;

    IF v_country_id IS NULL THEN
        INSERT INTO public.countries (country_name)
        VALUES (p_country_name)
        RETURNING country_id INTO v_country_id;
    END IF;

    -- 2. Handle state (lookup first, insert if not found)
    SELECT state_id
    INTO v_state_id
    FROM public.states
    WHERE country_id = v_country_id
      AND state_name = p_state_name;

    IF v_state_id IS NULL THEN
        INSERT INTO public.states (country_id, state_name)
        VALUES (v_country_id, p_state_name)
        RETURNING state_id INTO v_state_id;
    END IF;

    -- 3. Handle city (lookup first across ALL counties in the state)
    SELECT c.city_id, c.county_id
    INTO v_city_id, v_county_id
    FROM public.cities c
             JOIN public.counties co ON c.county_id = co.county_id
    WHERE co.state_id = v_state_id
      AND c.city_name = p_city_name
    LIMIT 1;

    -- If city doesn't exist, we need to create it, but first need a county
    IF v_city_id IS NULL THEN
        -- Look for an existing county in this state (prefer first one found)
        SELECT county_id
        INTO v_county_id
        FROM public.counties
        WHERE state_id = v_state_id
        LIMIT 1;

        -- If no county exists in this state, create a default one
        IF v_county_id IS NULL THEN
            INSERT INTO public.counties (state_id, county_name)
            VALUES (v_state_id, 'Unknown County')
            RETURNING county_id INTO v_county_id;
        END IF;

        -- Now create the city
        INSERT INTO public.cities (city_name, county_id)
        VALUES (p_city_name, v_county_id)
        RETURNING city_id INTO v_city_id;
    END IF;

    -- 4. Handle location (lookup first, insert if not found)
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

    -- 5. Handle data source (lookup first, insert if not found)
    SELECT source_id
    INTO v_source_id
    FROM public.data_sources
    WHERE source_name = p_source_name;

    IF v_source_id IS NULL THEN
        INSERT INTO public.data_sources (source_name)
        VALUES (p_source_name)
        RETURNING source_id INTO v_source_id;
    END IF;

    -- 6. Handle address (lookup first, insert if not found)
    SELECT address_id
    INTO v_address_id
    FROM public.addresses
    WHERE city_id = v_city_id
      AND (street_address = p_street_address OR (street_address IS NULL AND p_street_address IS NULL))
      AND (postal_code = p_postal_code OR (postal_code IS NULL AND p_postal_code IS NULL));

    IF v_address_id IS NULL THEN
        INSERT INTO public.addresses (street_address, city_id, postal_code)
        VALUES (p_street_address, v_city_id, p_postal_code)
        RETURNING address_id INTO v_address_id;
    END IF;

    -- 7. Handle crime category (lookup first, insert if not found)
    SELECT crime_category_id
    INTO v_crime_category_id
    FROM public.crime_categories
    WHERE category_name = p_crime_category_name;

    IF v_crime_category_id IS NULL THEN
        INSERT INTO public.crime_categories (category_name)
        VALUES (p_crime_category_name)
        RETURNING crime_category_id INTO v_crime_category_id;
    END IF;

    -- 8. Finally, insert the crime incident
    INSERT INTO public.crime_incidents (address_id,
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
FROM add_crime_incident(p_city_name := 'Tacoma',
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
FROM add_crime_incident(p_city_name := 'Tacoma',
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
FROM add_crime_incident(p_city_name := 'Tacoma',
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
FROM add_crime_incident(p_city_name := 'Tacoma',
                        p_state_name := 'Washington',
                        p_country_name := 'United States',
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
FROM add_crime_incident(p_city_name := 'Tacoma',
                        p_state_name := 'Washington',
                        p_country_name := 'United States',
                        p_source_name := 'City of Tacoma Reported Crime (Tacoma)',
                        p_latitude := 47.253788,
                        p_longitude := -122.444841,
                        p_street_address := '900 TACOMA AVE S',
                        p_postal_code := '98402',
                        p_crime_category_name := 'Destruction/Damage/Vandalism',
                        p_incident_date := '2024-01-18'::date,
                        p_incident_time := '06:00:00'::time,
                        p_case_num := '2401909037');
