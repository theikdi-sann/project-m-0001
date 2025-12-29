-- AUTHENTICATION & ADMINISTRATION

CREATE TABLE roles (
    id VARCHAR(25) PRIMARY KEY,
    name VARCHAR(30) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_roles_name (name)
    );

CREATE TABLE permissions (
    id VARCHAR(25) PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    description TEXT,
    can_create BOOLEAN DEFAULT FALSE,
    can_read BOOLEAN DEFAULT FALSE,
    can_update BOOLEAN DEFAULT FALSE,
    can_delete BOOLEAN DEFAULT FALSE,
	can_export BOOLEAN DEFAULT FALSE,
	can_import BOOLEAN DEFAULT FALSE,
    role_id VARCHAR(25) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (role_id) REFERENCES roles(id),
    INDEX idx_roles_name (name)
    );

CREATE TABLE users (
    id VARCHAR(25) PRIMARY KEY,
  	username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    fullname VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(15),
    role_id VARCHAR(25) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (role_id) REFERENCES roles(id),
    INDEX idx_users_role (role_id),
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_active (is_active)
);

-- SESSION MANAGEMENT

CREATE TABLE tables (
    id VARCHAR(25) PRIMARY KEY,
    table_number INT UNIQUE NOT NULL,
    capacity INT NOT NULL CHECK (capacity BETWEEN 1 AND 20),
    location_description VARCHAR(100),
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'occupied', 'reserved', 'maintenance')),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_tables_status (status)
);

CREATE TABLE session_statuses (
	id VARCHAR(25) PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	description TEXT
);

CREATE TABLE session_types (
	id VARCHAR(25) PRIMARY KEY,
	name VARCHAR(100) UNIQUE NOT NULL,
	display_name VARCHAR(255) NOT NULL,
	fixed_price DECIMAL(10,2) DEFAULT 0 NOT NULL
);

CREATE TABLE sessions (
	id VARCHAR(25) PRIMARY KEY,
	token VARCHAR(255) UNIQUE NOT NULL,
	table_id VARCHAR(25) NOT NULL,
	session_type_id VARCHAR(25) NOT NULL,
	
	customer_count INT,
	
	start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	end_time TIMESTAMP,
	expires_at TIMESTAMP,
	
	status_id VARCHAR(25) NOT NULL,
	
    created_by VARCHAR(25) NOT NULL,
    
    ended_by VARCHAR(25),
    
    customer_notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (table_id) REFERENCES tables(id),
    FOREIGN KEY (session_type_id) REFERENCES session_types(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (ended_by) REFERENCES users(id),
    
    INDEX idx_sessions_table (table_id),
    INDEX idx_sessions_status (status_id),
    INDEX idx_sessions_type (session_type_id),
    INDEX idx_sessions_timing (start_time, end_time)
);

-- MENU MANAGEMENT
    
CREATE TABLE menu_categories (
    id VARCHAR(25) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(50) NOT NULL,
    description TEXT,
    status ENUM('active', 'disabled') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_menu_category_name (name),
    INDEX idx_menu_category_status (status)
);

CREATE TABLE option_types (
	id VARCHAR(25) PRIMARY KEY,
	name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE options (
	id VARCHAR(25) PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	display_name VARCHAR(255) NOT NULL,
	option_type_id VARCHAR(25) NOT NULL,
	
	FOREIGN KEY (option_type_id) REFERENCES option_types(id)
);

CREATE TABLE option_values (
	id VARCHAR(25) PRIMARY KEY,
	option_id VARCHAR(25) NOT NULL,
	value VARCHAR(255) NOT NULL,
	
	FOREIGN KEY (option_id) REFERENCES options(id)
);

CREATE TABLE session_type_option_value_additional_prices (
	id VARCHAR(25) PRIMARY KEY,
	additional_price DECIMAL(10,2) NOT NULL DEFAULT 0,
	session_type_id VARCHAR(25) NOT NULL,
	option_value_id VARCHAR(25) NOT NULL,
	
	FOREIGN KEY (session_type_id) REFERENCES session_types(id),
	FOREIGN KEY (option_value_id) REFERENCES option_values(id)
);

CREATE TABLE menu_items (
    id VARCHAR(25) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id VARCHAR(25) NOT NULL,
    
    preparation_time INT,
    status VARCHAR(10) DEFAULT 'available' CHECK (status IN ('available', 'unavailable')),
    ingredients TEXT,
    chef_notes TEXT,
    image_url VARCHAR(255),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (category_id) REFERENCES menu_categories(id),
    INDEX idx_menu_items_category (category_id),
    INDEX idx_menu_items_status (status)
);

CREATE TABLE menu_item_options (
	id VARCHAR(25) PRIMARY KEY,
	menu_item_id VARCHAR(25) NOT NULL,
	option_id VARCHAR(25) NOT NULL,
	
	FOREIGN KEY (menu_item_id) REFERENCES menu_items(id),
	FOREIGN KEY (option_id) REFERENCES options(id)
);

CREATE TABLE pricing_models (
	id VARCHAR(25) PRIMARY KEY,
	type VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE menu_item_session_details (
	id VARCHAR(25) PRIMARY KEY,
	session_type_id VARCHAR(25) NOT NULL,
	menu_item_id VARCHAR(25) NOT NULL,
	price DECIMAL(10,2) NOT NULL,
	pricing_model_id VARCHAR(25) NOT NULL,
	
	FOREIGN KEY (session_type_id) REFERENCES session_types(id),
	FOREIGN KEY (menu_item_id) REFERENCES menu_items(id),
	FOREIGN KEY (pricing_model_id) REFERENCES pricing_models(id)
);


-- ORDER MANAGEMENT

CREATE TABLE orders (
    id VARCHAR(25) PRIMARY KEY,
    session_id VARCHAR(25) NOT NULL,
    created_by VARCHAR(25) NOT NULL,
    updated_by VARCHAR(25) NOT NULL,
    served_by VARCHAR(25) NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'Pending', -- e.g., Pending, Preparing, Ready, Served, Cancelled

    special_instructions TEXT,
    
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id),
    FOREIGN KEY (served_by) REFERENCES users(id)
);

CREATE TABLE order_items (
    id VARCHAR(25) PRIMARY KEY,
    order_id VARCHAR(25) NOT NULL,
    menu_item_id VARCHAR(25) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price_at_order DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL, -- (quantity * price_at_order)
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'QUEUED', -- e.g., QUEUED, COOKING, Ready, Served, Cancelled
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id)
);

CREATE TABLE order_item_selected_options (
    id VARCHAR(25) PRIMARY KEY,
    order_item_id VARCHAR(25) NOT NULL,
    option_value_id VARCHAR(25) NOT NULL,

    additional_price DECIMAL(10, 2) DEFAULT 0, -- Store extra cost at time of order (derived from session_type_option_value_additional_prices)

    FOREIGN KEY (order_item_id) REFERENCES order_items(id),
    FOREIGN KEY (option_value_id) REFERENCES option_values(id)
);

-- Invoice & Billing

CREATE TABLE invoices (
    id VARCHAR(25) PRIMARY KEY,
    invoice_number VARCHAR(50) UNIQUE NOT NULL, -- Human readable sequence (e.g., INV-2023-0001)
    session_id VARCHAR(25) UNIQUE NOT NULL, -- One invoice per session usually
    
    -- Monetary Calculations
    items_total DECIMAL(10,2) NOT NULL DEFAULT 0, -- Sum of order_items
    service_charge_amount DECIMAL(10,2) DEFAULT 0,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    
    -- This handles the fixed_price from the session_types table (e.g., Buffet entrance fee)
    session_fee DECIMAL(10,2) DEFAULT 0, 
    
    grand_total DECIMAL(10,2) NOT NULL,
    
    status VARCHAR(20) DEFAULT 'unpaid' CHECK (status IN ('unpaid', 'partially_paid', 'paid', 'void', 'refunded')),
    
    generated_by VARCHAR(25) NOT NULL, -- User who printed/created the bill
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (session_id) REFERENCES sessions(id),
    FOREIGN KEY (generated_by) REFERENCES users(id),
    
    INDEX idx_invoices_session (session_id),
    INDEX idx_invoices_status (status),
    INDEX idx_invoices_date (created_at)
);

CREATE TABLE invoice_adjustments (
    id VARCHAR(25) PRIMARY KEY,
    invoice_id VARCHAR(25) NOT NULL,
    name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL, -- e.g., "VAT (7%)", "SERVICE CHARGE", "DISCOUNT"

    category VARCHAR (20) NOT NULL CHECK (category IN ('TAX', 'SERVICE_CHARGE', "DISCOUNT", "OTHER")),

    type VARCHAR (20) NOT NULL CHECK (type IN ('FIXED', 'PERCENTAGE')),

    value DECIMAL(15, 2) DEFAULT 0.00, -- Stores either % or $ amount
    
    is_enabled BOOLEAN DEFAULT TRUE,

    created_by VARCHAR(25) NOT NULL,
    updated_by VARCHAR(25) NOT NULL,

    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE TABLE payment_methods (
    id VARCHAR(25) PRIMARY KEY,
    name VARCHAR(50) NOT NULL, -- Cash, Credit Card, Debit Card, QR Code, Voucher
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE payments (
    id VARCHAR(25) PRIMARY KEY,
    invoice_id VARCHAR(25) NOT NULL,
    payment_method_id VARCHAR(25) NOT NULL,
    
    amount DECIMAL(10,2) NOT NULL,
    
    transaction_reference VARCHAR(100), -- ID from the card terminal or bank transfer
    processed_by VARCHAR(25) NOT NULL, -- User who took the money
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (invoice_id) REFERENCES invoices(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id),
    FOREIGN KEY (processed_by) REFERENCES users(id),
    
    INDEX idx_payments_invoice (invoice_id)
);
