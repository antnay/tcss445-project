SELECT cc.category_name,
       latitude,
       longitude,
       street_address,
       city_name,
       state_name,
       country_name,
       postal_code,
       incident_date,
       incident_time
FROM crime_incidents_partition
         INNER JOIN public.addresses a ON crime_incidents_partition.address_id = a.address_id
         INNER JOIN public.cities c ON c.city_id = a.city_id
         INNER JOIN public.counties c2 ON c2.county_id = c.county_id
         INNER JOIN public.states s ON s.state_id = c2.state_id
         INNER JOIN public.countries c3 ON c3.country_id = s.country_id
         INNER JOIN public.crime_categories cc ON crime_incidents_partition.crime_category_id = cc.crime_category_id
         INNER JOIN public.locations l ON l.location_id = crime_incidents_partition.location_id;


SELECT *
FROM crime_stats_partition;


INSERT INTO crime_incidents_partition
SELECT *
FROM crime_incidents;

SELECT *
FROM crime_incidents_partition
WHERE incident_date BETWEEN '2025-01-01' AND '2025-12-31';

SELECT *
FROM addresses
where public.addresses.neighborhood_id is NULL;


DELETE FROM crime_incidents_partition a
WHERE EXISTS (
    SELECT 1
    FROM crime_incidents_partition b
    WHERE a.case_num = b.case_num
      AND a.incident_date = b.incident_date
      AND a.address_id = b.address_id
      AND a.ctid < b.ctid
);


WITH DuplicateCases AS (
    SELECT case_num, COUNT(*) as occurrence_count
    FROM crime_incidents_partition
    WHERE case_num IS NOT NULL
    GROUP BY case_num
    HAVING COUNT(*) > 1
)
SELECT
    cp.case_num,
    cp.incident_date,
    a.street_address,
    cc.category_name as crime_type,
    cp.incident_time
FROM crime_incidents_partition cp
         JOIN DuplicateCases dc ON cp.case_num = dc.case_num
         JOIN addresses a ON cp.address_id = a.address_id
         JOIN crime_categories cc ON cp.crime_category_id = cc.crime_category_id
ORDER BY cp.case_num, cp.incident_date;


WITH AddressCounts AS (
    SELECT cp.address_id, COUNT(*) as incident_count
    FROM crime_incidents_partition cp
    GROUP BY cp.address_id
    HAVING COUNT(*) > 1
)
SELECT
    a.street_address,
    ci.city_name,
    a.postal_code,
    cc.category_name as crime_type,
    cp.incident_date,
    cp.incident_time,
    cp.case_num,
    ac.incident_count as total_incidents_at_address
FROM crime_incidents_partition cp
         JOIN AddressCounts ac ON cp.address_id = ac.address_id
         JOIN addresses a ON cp.address_id = a.address_id
         JOIN cities ci ON a.city_id = ci.city_id
         JOIN crime_categories cc ON cp.crime_category_id = cc.crime_category_id
ORDER BY
    a.street_address,
    cp.incident_date DESC,
    cp.incident_time;


-- First, insert 20 user profiles with more realistic usernames
INSERT INTO user_profiles (username, email, role, is_active, created_at)
VALUES
    ('techExplorer42', 'alex.tech@example.com', 'user', true, NOW()),
    ('dataWizard89', 'sarah.code@example.com', 'user', true, NOW()),
    ('pythonMaster_23', 'mike.dev@example.com', 'admin', true, NOW()),
    ('webCrafter365', 'jessica.web@example.com', 'user', true, NOW()),
    ('codingNinja77', 'david.ninja@example.com', 'user', true, NOW()),
    ('debugQueen', 'emma.debug@example.com', 'user', FALSE, NOW()),
    ('systemGuru_42', 'james.sys@example.com', 'admin', true, NOW()),
    ('algorithmPro', 'lisa.algo@example.com', 'user', true, NOW()),
    ('devOpsHero', 'chris.ops@example.com', 'user', true, NOW()),
    ( 'securitySage', 'maria.sec@example.com', 'user', true, NOW()),
    ( 'cloudArchitect', 'peter.cloud@example.com', 'user', false, NOW()),
    ( 'aiResearcher', 'helen.ai@example.com', 'user', true, NOW()),
    ( 'networkWizard', 'robert.net@example.com', 'admin', true, NOW()),
    ( 'mobileDev_88', 'anna.mobile@example.com', 'user', true, NOW()),
    ( 'dataScientist42', 'thomas.data@example.com', 'user', true, NOW()),
    ( 'frontEndArtist', 'sophie.web@example.com', 'user', true, NOW()),
    ( 'backEndPro', 'daniel.backend@example.com', 'user', true, NOW()),
    ( 'systemAdmin_24', 'rachel.admin@example.com', 'admin', true, NOW()),
    ( 'mlEngineer', 'kevin.ml@example.com', 'user', true, NOW()),
    ( 'quantumCoder', 'laura.quantum@example.com', 'user', true, NOW());

-- Insert 10 password logins with bcrypt-like hashes
INSERT INTO password_login (user_id, password_hash, access_token, refresh_token, created_at)
VALUES
    (1, '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewKxpSeYyTo1I.Ae', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.8KJ2', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxLCJyZWZyZXNoIjp0cnVlfQ.9k3M', NOW()),
    (2, '$2a$12$k8Y6Zb7rkptZwKJvx.6RuOXPZpGGv89v.WO9LXD/1EspYL2tpfGi2', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.9mK3', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoyLCJyZWZyZXNoIjp0cnVlfQ.7j2N', NOW()),
    (3, '$2a$12$3.FmXyJ5K/n4m.8kl2Kk3O9M6H5M0ONzKF5g2.Y9GM4Z3ZyxsFp.', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozfQ.5pL4', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjozLCJyZWZyZXNoIjp0cnVlfQ.2m8P', NOW()),
    (4, '$2a$12$9.8kJz2/eM9K7U6TzJ.bS.4Q.ZGm4QUx6vxJ3/BiLYSA2dGaRnK7C', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0fQ.3nR5', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo0LCJyZWZyZXNoIjp0cnVlfQ.4k7L', NOW()),
    (5, '$2a$12$mK4.7uLXh9F8fz/qJsKTBOGd.HoEW.O7m3xZw0Ip4FwJ8cQF9Tx.S', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1fQ.1kM6', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo1LCJyZWZyZXNoIjp0cnVlfQ.6n5K', NOW()),
    (6, '$2a$12$Xk7.GxIGp.YPuE4qCj/SKeQZ9yJ5WBXm5U5WRZQxk5cZJyGSD8.a', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2fQ.7mP8', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo2LCJyZWZyZXNoIjp0cnVlfQ.5h4J', NOW()),
    (7, '$2a$12$QW.8kJHf2/N8K5U6TzJ.bS.4Q.ZGm4QUx6vxJ3/BiLYSA2dGaRnK7C', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3fQ.2pQ9', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo3LCJyZWZyZXNoIjp0cnVlfQ.8t6M', NOW()),
    (8, '$2a$12$Zt5.UjV4K9N.p/qRsTBxw.oPZg5gmq5xK4Q8gE9qJ3zP5F5s8Kj2', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4fQ.4nS7', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo4LCJyZWZyZXNoIjp0cnVlfQ.1p9N', NOW()),
    (9, '$2a$12$Rk9.HjI2p.n4K8U6TzJ.bS.4Q.ZGm4QUx6vxJ3/BiLYSA2dGaRnK7C', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo5fQ.6kL5', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjo5LCJyZWZyZXNoIjp0cnVlfQ.3m2P', NOW()),
    (10, '$2a$12$Yk2.MnV8K0N3p5U6TzJ.bS.4Q.ZGm4QUx6vxJ3/BiLYSA2dGaRnK7C', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMH0.5nT4', 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxMCwicmVmcmVzaCI6dHJ1ZX0.7k8M', NOW());

-- Insert 10 OAuth logins alternating between Google and GitHub
INSERT INTO oauth_login (user_id, provider, provider_user_id, access_token, refresh_token, created_at)
VALUES
    (11, 'google', 'g_115974234761957623851', 'ya29.a0AfH6SMBx7_GF8QN_KhsXj4J6HHs1dk2NF8QK9_GvjOM', '1//04dKzP9_HGf6JCgYIARAAGAQSNwF-L9IrqD7vYE5Uj2xs8JH_8aFGkw', NOW()),
    (12, 'github', 'gh_854123697', 'gho_16C7e42F3zWFw085jPEKN6Y2P9G3cD3pb7ij', 'ghr_J86xH8w085jPEKN6Y2P9G3cD3pb7kl9', NOW()),
    (13, 'google', 'g_234859760123497612384', 'ya29.a0AfH6SMCpL4K8_PQ9_MhK7Njk2L9IrHGf6JCgYI', '1//04tNz8_JKf6LMPgYIARAAGAQSNwF-L9IrWE2NF8QK9_GvjPM', NOW()),
    (14, 'github', 'gh_697854123', 'gho_26C7e42F3zWFw085jPEKN6Y2P9G3cD3xy8kl', 'ghr_K97yI9x085jPEKN6Y2P9G3cD3pb7mn0', NOW()),
    (15, 'google', 'g_345970812634897162534', 'ya29.a0AfH6SMDqM5N9_LiJ8_QR8Njk2L9IrHGf6JCgYI', '1//04uOA9_IHg7KNQgYIARAAGAQSNwF-L9IrXF3NF8QK9_GvkQN', NOW()),
    (16, 'github', 'gh_321985674', 'gho_36C7e42F3zWFw085jPEKN6Y2P9G3cD3vw4rs', 'ghr_L08zJ0y085jPEKN6Y2P9G3cD3pb7op1', NOW()),
    (17, 'google', 'g_456081923745806213645', 'ya29.a0AfH6SMEvR6O0_NjK9_SiK9Njk2L9IrHGf6JCgYI', '1//04vPB0_JIh8LORhYIARAAGAQSNwF-L9IrYG4NF8QK9_GvlRM', NOW()),
    (18, 'github', 'gh_147852369', 'gho_46C7e42F3zWFw085jPEKN6Y2P9G3cD3ut2qp', 'ghr_M19aK1z085jPEKN6Y2P9G3cD3pb7qr2', NOW()),
    (19, 'google', 'g_567192034856917324756', 'ya29.a0AfH6SMFwS7P1_OkL0_TjL0Njk2L9IrHGf6JCgYI', '1//04wQC1_KJi9MPSiYIARAAGAQSNwF-L9IrZH5NF8QK9_GvmSN', NOW()),
    (20, 'github', 'gh_963852741', 'gho_56C7e42F3zWFw085jPEKN6Y2P9G3cD3ts1po', 'ghr_N20bL2a085jPEKN6Y2P9G3cD3pb7st3', NOW());