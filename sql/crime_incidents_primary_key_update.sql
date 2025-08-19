-- First, drop any existing primary key if there is one
ALTER TABLE crime_incidents_partition DROP CONSTRAINT IF EXISTS crime_incidents_partition_pkey;

-- Add the primary key that includes both incident_id and incident_date
ALTER TABLE crime_incidents_partition 
    ADD PRIMARY KEY (incident_id, incident_date);
