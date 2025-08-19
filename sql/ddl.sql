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

