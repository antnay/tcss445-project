CREATE INDEX idx_crime_incidents_address ON crime_incidents (address_id);
CREATE INDEX idx_crime_incidents_category ON crime_incidents (crime_category_id);
CREATE INDEX idx_crime_incidents_location ON crime_incidents (location_id);

CREATE INDEX idx_addresses_city ON addresses (city_id);

CREATE INDEX idx_cities_county ON cities (county_id);

CREATE INDEX idx_counties_state ON counties (state_id);

CREATE INDEX idx_states_country ON states (country_id);

---------

CREATE INDEX idx_crime_incidents_date_time ON crime_incidents (incident_date, incident_time);
CREATE INDEX idx_crime_incidents_created ON crime_incidents (created_at);

CREATE INDEX idx_crime_incidents_case_num ON crime_incidents (case_num);

CREATE INDEX idx_crime_categories_name ON crime_categories (category_name);

CREATE INDEX idx_locations_coords ON locations (latitude, longitude);

CREATE INDEX idx_user_sessions_user ON user_sessions (user_id);
CREATE INDEX idx_user_sessions_timestamp ON user_sessions (created_at);

CREATE INDEX idx_oauth_login_user ON oauth_login (user_id);
CREATE INDEX idx_password_login_user ON password_login (user_id);

---------------

-- fill factor
ALTER TABLE crime_incidents
    SET (FILLFACTOR = 90);
ALTER TABLE user_sessions
    SET (FILLFACTOR = 85);

-- parition
CREATE TABLE crime_incidents_partition
(
    LIKE crime_incidents
) PARTITION BY RANGE (incident_date);


CREATE INDEX idx_crime_incidents_partition_address ON crime_incidents_partition (address_id);
CREATE INDEX idx_crime_incidents_partition_category ON crime_incidents_partition (crime_category_id);
CREATE INDEX idx_crime_incidents_partition_location ON crime_incidents_partition (location_id);
CREATE INDEX idx_crime_incidents_partition_date_time ON crime_incidents_partition (incident_date, incident_time);
CREATE INDEX idx_crime_incidents_partition_created ON crime_incidents_partition (created_at);
CREATE INDEX idx_crime_incidents_partition_case_num ON crime_incidents_partition (case_num);

CREATE TABLE crime_incidents_2010
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2010-01-01') TO ('2011-01-01');

CREATE TABLE crime_incidents_2011
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2011-01-01') TO ('2012-01-01');

CREATE TABLE crime_incidents_2012
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2012-01-01') TO ('2013-01-01');

CREATE TABLE crime_incidents_2013
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2013-01-01') TO ('2014-01-01');

CREATE TABLE crime_incidents_2014
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2014-01-01') TO ('2015-01-01');

CREATE TABLE crime_incidents_2015
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2015-01-01') TO ('2016-01-01');

CREATE TABLE crime_incidents_2016
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2016-01-01') TO ('2017-01-01');

CREATE TABLE crime_incidents_2017
    PARTITION OF crime_incidents_partition
        FOR VALUES FROM ('2017-01-01') TO ('2018-01-01');



INSERT INTO crime_incidents_partition
(incident_id, address_id, crime_category_id, incident_date, incident_time,
 location_id, case_num, is_resolved, source_id, created_at)
VALUES (1, 1, 1, '2025-01-15', '14:30:00', 1, 'CASE001', FALSE, 1, CURRENT_TIMESTAMP);


ALTER TABLE crime_incidents
    ALTER COLUMN incident_date SET STATISTICS 1000;
ALTER TABLE crime_incidents
    ALTER COLUMN location_id SET STATISTICS 1000;

ALTER TABLE crime_incidents
    CLUSTER ON idx_crime_incidents_date_time;


ANALYZE crime_incidents;
ANALYZE crime_incidents_partition;
ANALYZE addresses;
ANALYZE locations;

ALTER TABLE crime_incidents
    SET (
        AUTOVACUUM_VACUUM_SCALE_FACTOR = 0.1,
        AUTOVACUUM_ANALYZE_SCALE_FACTOR = 0.05
        );


CREATE MATERIALIZED VIEW crime_stats_partition AS
SELECT cc.category_name,
       c.city_name,
       s.state_name,
       COUNT(*)                              AS incident_count,
       DATE_TRUNC('month', ci.incident_date) AS month
FROM crime_incidents_partition ci
         JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
         JOIN addresses a ON ci.address_id = a.address_id
         JOIN cities c ON a.city_id = c.city_id
         JOIN counties co ON c.county_id = co.county_id
         JOIN states s ON co.state_id = s.state_id
GROUP BY cc.category_name, c.city_name, s.state_name, DATE_TRUNC('month', ci.incident_date)
WITH DATA;

CREATE INDEX idx_crime_stats_category_partition ON crime_stats_partition (category_name);
CREATE INDEX idx_crime_stats_location_partition ON crime_stats_partition (city_name, state_name);


REINDEX TABLE crime_incidents;



