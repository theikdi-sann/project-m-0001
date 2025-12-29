CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(30) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_roles_name (name)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(30) NOT NULL,
    description TEXT,
    can_create BOOLEAN DEFAULT FALSE,
    can_read BOOLEAN DEFAULT FALSE,
    can_update BOOLEAN DEFAULT FALSE,
    can_delete BOOLEAN DEFAULT FALSE,
	can_export BOOLEAN DEFAULT FALSE,
	can_import BOOLEAN DEFAULT FALSE,
    role_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    INDEX idx_permissions_name (name)
    );

CREATE TABLE users (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  	username VARCHAR(50) UNIQUE NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    phone VARCHAR(15),
    role_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT,

    INDEX idx_users_username (username),
    INDEX idx_users_phone (phone),
    INDEX idx_users_active (is_active)
);

CREATE TABLE tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_number INT UNIQUE NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0),
    location_description VARCHAR(100),
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved', 'maintenance')),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_tables_status (status)
);

CREATE TABLE session_statuses (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) NOT NULL,
	description TEXT
);

CREATE TABLE session_types (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) UNIQUE NOT NULL,
	display_name VARCHAR(255) NOT NULL,
	fixed_price DECIMAL(10,2) DEFAULT 0 NOT NULL
);


CREATE TABLE sessions (
   	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	token VARCHAR(255) UNIQUE NOT NULL,
	table_id UUID NOT NULL,
	session_type_id UUID NOT NULL,
	
	customer_count INT,
	
	start_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	end_time TIMESTAMPTZ,
	expires_at TIMESTAMPTZ,
	
	status_id UUID NOT NULL,
	
    created_by UUID NOT NULL,
    
    ended_by UUID,
    
    customer_notes TEXT,
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE RESTRICT,
    FOREIGN KEY (session_type_id) REFERENCES session_types(id) ON DELETE RESTRICT,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (ended_by) REFERENCES users(id) ON DELETE RESTRICT,
    
    INDEX idx_sessions_table (table_id),
    INDEX idx_sessions_status (status_id),
    INDEX idx_sessions_type (session_type_id),
    INDEX idx_sessions_timing (start_time, end_time)
);

CREATE TABLE menu_categories (
   	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(10) DEFAULT 'active' CHECK (status IN ('active', 'disabled')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_menu_category_name (name),
    INDEX idx_menu_category_status (status)
);

CREATE TABLE option_types (
   	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE options (
   	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) UNIQUE NOT NULL,
	display_name VARCHAR(255) NOT NULL,
	option_type_id UUID NOT NULL,
	
	FOREIGN KEY (option_type_id) REFERENCES option_types(id) ON DELETE RESTRICT
);

CREATE TABLE option_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	option_id UUID NOT NULL,
	value VARCHAR(255) NOT NULL,
	
	FOREIGN KEY (option_id) REFERENCES options(id) ON DELETE CASCADE
);

CREATE TABLE session_type_option_value_additional_prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	additional_price DECIMAL(10,2) NOT NULL DEFAULT 0,
	session_type_id UUID NOT NULL,
	option_value_id UUID NOT NULL,
	
	FOREIGN KEY (session_type_id) REFERENCES session_types(id) ON DELETE CASCADE,
	FOREIGN KEY (option_value_id) REFERENCES option_values(id) ON DELETE CASCADE
);

CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID NOT NULL,
    
    preparation_time INT,
    status VARCHAR(10) DEFAULT 'available' CHECK (status IN ('available', 'unavailable')),
    ingredients TEXT,
    chef_notes TEXT,
    image_url VARCHAR(255),
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (category_id) REFERENCES menu_categories(id) ON DELETE RESTRICT,

    INDEX idx_menu_items_category (category_id),
    INDEX idx_menu_items_status (status)
);

CREATE TABLE menu_item_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	menu_item_id UUID NOT NULL,
	option_id UUID NOT NULL,
	
	FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE,
	FOREIGN KEY (option_id) REFERENCES options(id) ON DELETE CASCADE
);

CREATE TABLE pricing_models (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	type VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE menu_item_session_details (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	session_type_id UUID NOT NULL,
	menu_item_id UUID NOT NULL,
	price DECIMAL(10,2) NOT NULL,
	pricing_model_id UUID NOT NULL,
	
	FOREIGN KEY (session_type_id) REFERENCES session_types(id) ON DELETE CASCADE,
	FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE,
	FOREIGN KEY (pricing_model_id) REFERENCES pricing_models(id) ON DELETE RESTRICT
);


CREATE TABLE orders (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL,
    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,
    served_by UUID NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'Pending', -- e.g., Pending, Preparing, Ready, Served, Cancelled

    special_instructions TEXT,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE RESTRICT,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (served_by) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE order_items (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    menu_item_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price_at_order DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL, -- (quantity * price_at_order)
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'QUEUED', -- e.g., QUEUED, COOKING, Ready, Served, Cancelled
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE RESTRICT
);

CREATE TABLE order_item_selected_options (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_item_id UUID NOT NULL,
    option_value_id UUID NOT NULL,

    additional_price DECIMAL(10, 2) DEFAULT 0, -- Store extra cost at time of order (derived from session_type_option_value_additional_prices)

    FOREIGN KEY (order_item_id) REFERENCES order_items(id) ON DELETE CASCADE,
    FOREIGN KEY (option_value_id) REFERENCES option_values(id) ON DELETE RESTRICT
);

-- Invoice & Billing

CREATE TABLE invoices (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_number VARCHAR(50) UNIQUE NOT NULL, -- Human readable sequence (e.g., INV-2023-0001)
    session_id UUID NOT NULL, -- One invoice per session usually
    
    -- Monetary Calculations
    items_total DECIMAL(10,2) NOT NULL DEFAULT 0, -- Sum of order_items
    service_charge_amount DECIMAL(10,2) DEFAULT 0,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    
    -- This handles the fixed_price from the session_types table (e.g., Buffet entrance fee)
    session_fee DECIMAL(10,2) DEFAULT 0, 
    
    grand_total DECIMAL(10,2) NOT NULL,
    
    status VARCHAR(20) DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'partially_paid', 'paid', 'void', 'refunded')),
    
    generated_by UUID NOT NULL, -- User who printed/created the bill

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE RESTRICT,
    FOREIGN KEY (generated_by) REFERENCES users(id) ON DELETE RESTRICT,
    
    INDEX idx_invoices_session (session_id),
    INDEX idx_invoices_status (status),
    INDEX idx_invoices_date (created_at)
);

CREATE TABLE invoice_adjustments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL, -- e.g., "VAT (7%)", "SERVICE CHARGE", "DISCOUNT"

    category VARCHAR (20) NOT NULL CHECK (category IN ('TAX', 'SERVICE_CHARGE', "DISCOUNT", "OTHER")),

    type VARCHAR (20) NOT NULL CHECK (type IN ('FIXED', 'PERCENTAGE')),

    value DECIMAL(15, 2) DEFAULT 0.00, -- Stores either % or $ amount
    
    is_enabled BOOLEAN DEFAULT TRUE,

    created_by UUID NOT NULL,
    updated_by UUID NOT NULL,

    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE payment_methods (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL, -- Cash, Credit Card, Debit Card, QR Code, Voucher
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE payments (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    payment_method_id UUID NOT NULL,
    
    amount DECIMAL(10,2) NOT NULL,
    
    transaction_reference VARCHAR(100), -- ID from the card terminal or bank transfer
    processed_by UUID NOT NULL, -- User who took the money
    
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE RESTRICT,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id) ON DELETE RESTRICT,
    FOREIGN KEY (processed_by) REFERENCES users(id) ON DELETE RESTRICT,
    
    INDEX idx_payments_invoice (invoice_id)
);


