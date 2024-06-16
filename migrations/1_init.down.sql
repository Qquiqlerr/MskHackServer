-- Disable foreign key constraints
ALTER TABLE route_features
    DROP CONSTRAINT IF EXISTS route_features_route_id_fkey;
ALTER TABLE route_features
    DROP CONSTRAINT IF EXISTS route_features_feature_id_fkey;
ALTER TABLE visit_permits
    DROP CONSTRAINT IF EXISTS visit_permits_route_id_fkey;
ALTER TABLE visit_permits
    DROP CONSTRAINT IF EXISTS visit_permits_group_id_fkey;
ALTER TABLE visit_permits
    DROP CONSTRAINT IF EXISTS visit_permits_visit_reason_fkey;
ALTER TABLE visit_permits
    DROP CONSTRAINT IF EXISTS visit_permits_visit_format_fkey;
ALTER TABLE visit_permits
    DROP CONSTRAINT IF EXISTS visit_permits_status_fkey;
ALTER TABLE visit_permits_photo_types
    DROP CONSTRAINT IF EXISTS visit_permits_photo_types_visit_permit_id_fkey;
ALTER TABLE visit_permits_photo_types
    DROP CONSTRAINT IF EXISTS visit_permits_photo_types_photo_type_id_fkey;
ALTER TABLE reports
    DROP CONSTRAINT IF EXISTS reports_type_fkey;
ALTER TABLE reports
    DROP CONSTRAINT IF EXISTS reports_statusID_fkey;
-- Drop tables in the correct order
DROP TABLE IF EXISTS route_features;
DROP TABLE IF EXISTS visit_permits;
DROP TABLE IF EXISTS group_permits;
DROP TABLE IF EXISTS permit_statuses;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS visit_reasons;
DROP TABLE IF EXISTS photo_types;
DROP TABLE IF EXISTS visit_format;
DROP TABLE IF EXISTS feature_types;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS zones;
DROP TABLE IF EXISTS type_of_reports;
DROP TABLE IF EXISTS reports_statuses;


