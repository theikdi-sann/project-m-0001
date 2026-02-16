-- Migration: Make table_id nullable for Take Away
ALTER TABLE dining_sessions ALTER COLUMN table_id DROP NOT NULL;
