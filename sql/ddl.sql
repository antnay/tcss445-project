-- CREATE TYPE provider AS ENUM ('google', 'github');
-- CREATE TYPE role AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS public.locations
(
    location_id bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    latitude    numeric(10, 7)                      NOT NULL,
    longitude   numeric(10, 7)                      NOT NULL,
    CHECK (latitude >= -90.0 AND latitude <= 90.0),
    CHECK (longitude >= -180.0 AND longitude <= 180.0),
    UNIQUE (latitude, longitude)
);

CREATE TABLE IF NOT EXISTS public.countries
(
    country_id   bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    country_name varchar(50)                         NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.states
(
    state_id   bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    country_id bigint                              NOT NULL,
    state_name varchar(100)                        NOT NULL,
    FOREIGN KEY (country_id) REFERENCES public.countries (country_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    UNIQUE (country_id, state_name)
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

CREATE TABLE IF NOT EXISTS public.addresses
(
    address_id     bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    street_address varchar(255),
    city_id        bigint                              NOT NULL,
    postal_code    varchar(20),
    FOREIGN KEY (city_id) REFERENCES public.cities (city_id) ON UPDATE CASCADE ON DELETE RESTRICT
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
    UNIQUE (provider, provider_user_id),
    UNIQUE (access_token)
);

CREATE TABLE IF NOT EXISTS public.user_sessions
(
    session_id    bigint GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    user_id       bigint                              NOT NULL,
    session_token varchar(255)                        NOT NULL UNIQUE,
    expires_at    timestamp with time zone            NOT NULL,
    created_at    timestamp with time zone            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.user_profiles (user_id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CHECK (expires_at > created_at)
);
