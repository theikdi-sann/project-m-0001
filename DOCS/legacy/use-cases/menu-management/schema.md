-- Menu Items Table
CREATE TABLE menu_items (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    a_la_carte_price DECIMAL(10,2) NOT NULL,
    available_for_buffet BOOLEAN DEFAULT FALSE,
    buffet_pricing_model ENUM('always_included', 'always_extra_charge', 'tier_based'),
    always_extra_charge_price DECIMAL(10,2) DEFAULT 0,
    status ENUM('available', 'unavailable') DEFAULT 'available',
    preparation_time INT,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Option Groups Table
CREATE TABLE option_groups (
    id INT PRIMARY KEY AUTO_INCREMENT,
    menu_item_id INT,
    group_name VARCHAR(100) NOT NULL,
    required BOOLEAN DEFAULT FALSE,
    selection_type ENUM('single', 'multiple') DEFAULT 'single',
    display_order INT DEFAULT 0,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE
);

-- Option Choices Table
CREATE TABLE option_choices (
    id INT PRIMARY KEY AUTO_INCREMENT,
    option_group_id INT,
    choice_name VARCHAR(100) NOT NULL,
    additional_price DECIMAL(10,2) DEFAULT 0,
    available BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0,
    FOREIGN KEY (option_group_id) REFERENCES option_groups(id) ON DELETE CASCADE
);