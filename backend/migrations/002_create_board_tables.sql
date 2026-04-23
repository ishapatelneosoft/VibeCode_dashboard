-- Migration: 002_create_board_tables.sql
-- Creates columns and tasks tables for shared board functionality

-- Columns table (fixed 8 columns for shared board)
CREATE TABLE IF NOT EXISTS columns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    "order" INTEGER NOT NULL UNIQUE CHECK ("order" >= 0 AND "order" <= 7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    column_id UUID NOT NULL REFERENCES columns(id) ON DELETE CASCADE,
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    due_date TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_tasks_column_id ON tasks(column_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks(created_by);
CREATE INDEX IF NOT EXISTS idx_tasks_column_position ON tasks(column_id, position);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);

-- Seed the 8 fixed columns
INSERT INTO columns (id, name, "order") VALUES
    (gen_random_uuid(), 'Backlog', 0),
    (gen_random_uuid(), 'To Do', 1),
    (gen_random_uuid(), 'In Progress', 2),
    (gen_random_uuid(), 'Review', 3),
    (gen_random_uuid(), 'Testing', 4),
    (gen_random_uuid(), 'Done', 5),
    (gen_random_uuid(), 'Blocked', 6),
    (gen_random_uuid(), 'Archived', 7)
ON CONFLICT DO NOTHING;

-- Function to update updated_at timestamp for tasks
CREATE OR REPLACE FUNCTION update_tasks_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at for tasks
CREATE TRIGGER update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_tasks_updated_at_column();