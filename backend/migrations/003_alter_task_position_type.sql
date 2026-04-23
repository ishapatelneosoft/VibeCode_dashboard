-- Migration: 003_alter_task_position_type.sql
-- Changes task.position column from INTEGER to DOUBLE PRECISION to support fractional indexing

-- First, drop the index that depends on the column (optional, but safe)
DROP INDEX IF EXISTS idx_tasks_column_position;

-- Alter column type
ALTER TABLE tasks 
ALTER COLUMN position TYPE DOUBLE PRECISION;

-- Recreate the index
CREATE INDEX idx_tasks_column_position ON tasks(column_id, position);

-- Update any existing integer positions to float (they will be automatically converted)
-- No explicit conversion needed as PostgreSQL handles INTEGER -> DOUBLE PRECISION