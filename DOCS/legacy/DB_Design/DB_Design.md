# Database Design

This document outlines the database schema for the Restaurant Management System. The schema is optimized for MySQL/MariaDB and managed via Prisma.

## 1. ENUM Definitions

### `UserRole`
- `ADMIN`
- `MANAGER`
- `WAITER`
- `KITCHEN`
- `FRONT_DESK`
- `STAFF`
- `CUSTOMER`

### `TableStatus`
- `AVAILABLE`
- `OCCUPIED`
- `RESERVED`
- `MAINTENANCE`

### `DiningSessionStatus`
- `ACTIVE`
- `BILL_REQUESTED`
- `AWAITING_PAYMENT`
- `COMPLETED`
- `CANCELLED`

<!-- [X]
### `PricingModelType`
- `INCLUDED`
- `EXTRA_CHARGE` -->

<!-- [X] -->
### `OrderStatus`
- `PENDING`
- `PREPARING`
- `READY`
- `SERVED`
- `CANCELLED`

<!-- [X] -->
### `OrderItemStatus`
- `QUEUED`
- `COOKING`
- `READY`
- `SERVED`
- `CANCELLED`

<!-- [X] -->
### `InvoiceStatus`
- `UNPAID`
- `PARTIALLY_PAID`
- `PAID`
- `VOID`
- `REFUNDED`

<!-- [X] -->
### `AdjustmentCategory`
- `TAX`
- `SERVICE_CHARGE`
- `DISCOUNT`
- `OTHER`

<!-- [X] -->
### `AdjustmentType`
- `FIXED`
- `PERCENTAGE`

## 2. Authentication & Administration

### `users`
System users (Staff, Managers, etc.) - Integrated with NextAuth.
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(255)
- `username`: VARCHAR(50) UNIQUE
- `email`: VARCHAR(100) UNIQUE
- `emailVerified`: TIMESTAMP
- `image`: VARCHAR(255)
- `role`: UserRole (ENUM, Default: STAFF)
- `is_active` : BOOLEAN DEFAULT TRUE
- `createdAt`: TIMESTAMP
- `updatedAt`: TIMESTAMP

### `permissions`
Granular access control linked to roles.
- `id`: VARCHAR(25) PRIMARY KEY
- `role`: UserRole (ENUM) - The role this permission applies to.
- `resource`: VARCHAR(50) - e.g., 'ORDER', 'MENU', 'INVOICE', 'USER'.
- `action`: VARCHAR(20) - e.g., 'CREATE', 'READ', 'UPDATE', 'DELETE', 'MANAGE', 'EXPORT'.
- `description`: TEXT

## 3. Session Management

### `tables`
Physical dining tables in the restaurant.
- `id`: VARCHAR(25) PRIMARY KEY
- `table_number`: INT UNIQUE
- `capacity`: INT
- `location_description`: VARCHAR(100)
- `status`: TableStatus (ENUM)

### `dining_session_types`
Types of dining experiences (e.g., Buffet, À la carte).
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(100) UNIQUE
- `display_name`: VARCHAR(255)
- `fixed_price`: DECIMAL(10,2) (Used for buffet entry fees)

### `dining_sessions`
Active dining sessions linked to a table.
- `id`: VARCHAR(25) PRIMARY KEY
- `token`: VARCHAR(255) UNIQUE (For QR access)
- `table_id`: VARCHAR(25) (FK to tables)
- `dining_session_type_id`: VARCHAR(25) (FK to dining_session_types)
- `customer_count`: INT
- `start_time`: TIMESTAMP
- `end_time`: TIMESTAMP
- `expires_at`: TIMESTAMP
- `status`: DiningSessionStatus (ENUM)
- `created_by`: VARCHAR(25) (FK to users)
- `ended_by`: VARCHAR(25) (FK to users)
- `customer_notes`: TEXT
- `price_list_id`: VARCHAR(25) (FK to price_list)

### `price_list`
- `id`: PRIMARY KEY
- `name`: VARCHAR(100) UNIQUE
- `currency`: VARCHAR(3) DEFAULT 'MMK'
- `is_active`: BOOLEAN DEFAULT TRUE

## 4. Menu Management

### `menu_categories`
Groups for menu items.
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(50) UNIQUE
- `description`: TEXT
- `is_enabled`: BOOLEAN DEFAULT TRUE

### `menu_items`
Individual food and drink items.
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(255) UNIQUE
- `description`: TEXT
- `category_id`: VARCHAR(25) (FK to menu_categories)
- `preparation_time`: INT (minutes)
- `is_available`: BOOLEAN DEFAULT TRUE
- `ingredients`: TEXT
- `chef_notes`: TEXT
- `image_url`: VARCHAR(255)

### `option_types`
Defines the selection mechanism for the option group.
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(100) UNIQUE (e.g., 'CHECKBOX', 'RADIO_BOX', 'QUANTITY')

### `options`
Specific options (e.g., "Beef Doneness", "Pizza Toppings").
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(255) UNIQUE
- `display_name`: VARCHAR(255)
- `option_type_id`: VARCHAR(25) (FK to option_types)

### `option_values`
The specific choices available within an option (e.g., "Rare", "Medium", "Mushroom").
- `id`: VARCHAR(25) PRIMARY KEY
- `option_id`: VARCHAR(25) (FK to options)
- `name`: VARCHAR(255) (The display text, e.g., "Medium")
- `display_order`: INT (Sort order: 1 for Rare, 2 for Medium...)
- `is_default`: BOOLEAN DEFAULT FALSE (Pre-select this value?)
- `is_available`: BOOLEAN DEFAULT TRUE (Mark specific choice as out of stock)

### `menu_item_options`
Many-to-many relationship between items and options.
- `id`: VARCHAR(25) PRIMARY KEY
- `menu_item_id`: VARCHAR(25) (FK to menu_items)
- `option_id`: VARCHAR(25) (FK to options)

<!-- ### `menu_item_dining_session_details`
Session-specific pricing for menu items.
- `id`: VARCHAR(25) PRIMARY KEY
- `dining_session_type_id`: VARCHAR(25) (FK to session_types)
- `menu_item_id`: VARCHAR(25) (FK to menu_items)
- `price`: DECIMAL(10,2)
- `pricing_model`: PricingModelType (ENUM) -->

### `price_list_items`
- `id`: VARCHAR(25) PRIMARY KEY,
- `menu_item_id`: VARCHAR(25) (FK to menu_items)
- `price_list_id`: VARCHAR(25) (FK to price_list)
- `is_included`: BOOLEAN DEFAULT TRUE
- `price`: DECIMAL(10,2)
UNIQUE(price_list_id, menu_item_id)

<!-- 
### `dining_session_type_option_value_additional_prices`
Session-specific surcharges for options.
- `id`: VARCHAR(25) PRIMARY KEY
- `additional_price`: DECIMAL(10,2)
- `session_type_id`: VARCHAR(25) (FK to session_types)
- `option_value_id`: VARCHAR(25) (FK to option_values) -->

### `price_list_option_values`
- `id`: VARCHAR(25)
- `price_list_id`: VARCHAR(25) (FK to price_lists)
- `option_value_id`: VARCHAR(25) (FK to option_values)
- `price`: DECIMAL(10,2)
- `is_included`: BOOLEAN DEFAULT TRUE
UNIQUE(price_list_id, option_value_id)

## 5. Order Management

### `orders`
Groups of items ordered together.
- `id`: VARCHAR(25) PRIMARY KEY
- `dining_session_id`: VARCHAR(25) (FK to sessions)
- `created_by`: VARCHAR(25) (FK to users) ('CUSTOMER', 'WAITER', 'FRONT_DESK')
- `updated_by`: VARCHAR(25) (FK to users) ('WAITER', 'CUSTOMER', 'FRONT_DESK', 'KITCHEN')
- `served_by`: VARCHAR(25) (FK to users) ('WAITER')
- `status`: OrderStatus (ENUM)
- `special_instructions`: TEXT

### `order_items`
Individual items within an order.
- `id`: VARCHAR(25) PRIMARY KEY
- `order_id`: VARCHAR(25) (FK to orders)
- `menu_item_id`: VARCHAR(25) (FK to menu_items)
- `menu_name_at_order`: VARCHAR(255)
- `quantity`: INT
- `price_at_order`: DECIMAL(10,2)
- `subtotal`: DECIMAL(10,2) (Calculated: (price_at_order + (options_price * option_qty)) * quantity) 
- `notes`: TEXT
- `status`: OrderItemStatus (ENUM)
- `cancellation_reason`: TEXT (Required if status is CANCELLED)

### `order_item_selected_options`
Options selected for a specific order item.
- `id`: VARCHAR(25) PRIMARY KEY
- `order_item_id`: VARCHAR(25) (FK to order_items)
- `option_value_id`: VARCHAR(25) (FK to option_values)
- `quantity`: INT (For 'double' 'triple' a specific topping/option)
- `additional_price_at_order`: DECIMAL(10,2) (Captured at time of order)

<!-- [X] -->
## 6. Invoicing & Billing

### `invoices`
Final bill for a dining session.
- `id`: VARCHAR(25) PRIMARY KEY
- `invoice_number`: VARCHAR(50) UNIQUE
- `session_id`: VARCHAR(25)
- `items_total`: DECIMAL(10,2)
- `service_charge_amount`: DECIMAL(10,2)
- `tax_amount`: DECIMAL(10,2)
- `discount_amount`: DECIMAL(10,2)
- `dining_session_fee`: DECIMAL(10,2) (Fixed fee from session type)
- `grand_total`: DECIMAL(10,2)
- `status`: InvoiceStatus (ENUM)
- `generated_by`: VARCHAR(25) (FK to users) ('FRONT_DESK')

### `invoice_adjustments`
Taxes, service charges, and discounts applied to invoices.
- `id`: VARCHAR(25) PRIMARY KEY
- `invoice_id`: VARCHAR(25) (FK to invoices)
- `name`: VARCHAR(255)
- `display_name`: VARCHAR(255)
- `category`: AdjustmentCategory (ENUM)
- `type`: AdjustmentType (ENUM)
- `value`: DECIMAL(15,2)
- `is_enabled`: BOOLEAN DEFAULT TRUE

### `payment_methods`
Available payment options (Cash, Card, etc.).
- `id`: VARCHAR(25) PRIMARY KEY
- `name`: VARCHAR(50) UNIQUE
- `is_active`: BOOLEAN DEFAULT TRUE

### `payments`
Transaction records for invoices.
- `id`: VARCHAR(25) PRIMARY KEY
- `invoice_id`: VARCHAR(25) (FK to invoices)
- `payment_method_id`: VARCHAR(25) (FK to payment_methods)
- `amount`: DECIMAL(15,2)
- `change`: DECIMAL(15,2)
- `transaction_reference`: VARCHAR(100)
- `processed_by`: VARCHAR(25) (FK to users)