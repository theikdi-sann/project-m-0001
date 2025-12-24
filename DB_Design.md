# AUTHENTICATION

user {
    id      : VARCHAR(25)
    username: VARCHAR(50) UNIQUE
    password: VARCHAR(50)
    fullname: VARCHAR(100)
    email   : VARCHAR(100) UNIQUE
    phone   : VARCHAR(15)
    role    : Role (DEFAULT Waiter)
    isActive: BOOLEAN (DEFAULT TRUE)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
}

user_session {
    id            : VARCHAR(25)
    user_id       : VARCHAR(50)
    session_token : VARCHAR(255)
    ip_address    : INET
    user_agent    : VARCHAR(255)
    expires_at    : TIMESTAMP
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
}

role {
    id:        : VARCHAR(25)
    name       : VARCHAR(30)       
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
}

permission {
    id: 
    role_id: 
    name: VARCHAR(25)
    can_create : bool
    can_read : bool
    can_update : bool
    can_delete : bool
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
}

# ADMINISTRATION

# MENU MANAGEMENT

menu_category {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    description TEXT,
    status : ENUM('active', 'inactive'),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
}

menu_item {
    id : VARCHAR(25),
    name : VARCHAR(255) NOT NULL,
    description : VARCHAR(255),
    category_id : menu_categories(id),
    
    -- Operational
    preparation_time INT, -- in minutes
    status VARCHAR(10) DEFAULT 'available' CHECK (status IN ('available', 'unavailable')),
    ingredients TEXT,
    chef_notes TEXT,
    image_url VARCHAR(255),

    
    -- Options
    option
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_menu_items_category (category_id),
    INDEX idx_menu_items_status (status),
    INDEX idx_menu_items_buffet (available_for_buffet)
}

options {
    id: 
    name: 
    
}

menu_item_session_details {
    id : 
    session_type_id : 
    menu_item_id : 
    price : 
    pricing_model : ENUM ("INCLUDED", "EXTRA_CHARGE")
}

# ORDER MANAGEMENT

# SESSION & BILLING

session (
    id : VARCHAR(25)
    token: VARCHAR(255) UNIQUE
    table_id : Table(id),
    session_type : SessionTypes(id)
    
    customer_count : INTEGER,
    
    -- Timing
    start_time:  TIMESTAMP NOT NULL,
    end_time : TIMESTAMP,
    expires_at : TIMESTAMP,
    
    -- Status
    status : ENUM('active', 'bill_requested', "awaiting_payment", 'ended', 'paid', 'cancelled'),
    
    -- Staff assignment
    created_by UUID NOT NULL REFERENCES users(id),
    
    -- Customer information
    customer_notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_sessions_table (table_id),
    INDEX idx_sessions_status (status),
    INDEX idx_sessions_type (session_type),
    INDEX idx_sessions_timing (start_time, end_time),
);

session_type {
    id           : VARCHAR(25)
    name         : VARCHAR(255)
    display_name : VARCAHR(255)
}

payment {
    id                : VARCHAR(25),
    invoice_id        : Invoices(id),
    payment_method_id : payment_method(id),
    amount            : DECIMAL,
    transaction_id    : VARCHAR(100),
    status            : ENUM('pending', 'completed', 'failed', 'refunded'),
}

payment_method {
    id   : VARCHAR(25)
    name : VARCHAR(50)
}


# MYSQL

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

    SHOW TABLES;
    SELECT * FROM permissions;