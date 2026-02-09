-- Initial Schema Migration

CREATE TABLE IF NOT EXISTS dining_tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_number TEXT NOT NULL UNIQUE,
    capacity INT NOT NULL DEFAULT 2,
    status TEXT NOT NULL DEFAULT 'available', -- available, occupied, reserved, maintenance
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dining_session_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE, -- e.g., 'Standard Buffet', 'Premium Buffet', 'À la carte'
    description TEXT,
    is_buffet BOOLEAN NOT NULL DEFAULT FALSE,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    duration_minutes INT, -- NULL for à la carte
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    full_name TEXT,
    role TEXT NOT NULL, -- admin, manager, frontdesk, waiter, kitchen
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dining_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id UUID NOT NULL REFERENCES dining_tables(id),
    session_type_id UUID NOT NULL REFERENCES dining_session_types(id),
    created_by UUID NOT NULL REFERENCES users(id),
    guest_count INT NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'active',
    start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    price_per_guest DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
