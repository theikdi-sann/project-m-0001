# ERD and Models

This document outlines the database design for the Restaurant QR-based Ordering and Management System (RQRMOS). It includes the Entity-Relationship Diagram (ERD), the database schema in SQL, and the corresponding data models.

## Entity-Relationship Diagram (ERD)

The following diagram illustrates the relationships between the different entities in the system.

```mermaid
erDiagram
    USERS ||--o{ SESSIONS : "creates"
    USERS ||--o{ ORDERS : "places"
    ROLES ||--o{ USERS : "has"
    TABLES ||--o{ SESSIONS : "is on"
    SESSION_TYPES ||--o{ SESSIONS : "is of type"
    SESSIONS ||--o{ ORDERS : "has"
    SESSIONS ||--|| INVOICES : "generates"
    ORDERS ||--o{ ORDER_ITEMS : "contains"
    MENU_ITEMS ||--o{ ORDER_ITEMS : "is"
    MENU_CATEGORIES ||--o{ MENU_ITEMS : "belongs to"
    INVOICES ||--o{ INVOICE_ITEMS : "details"
    INVOICES ||--|| PAYMENTS : "receives"

    USERS {
        int user_id PK
        varchar(255) username
        varchar(255) password_hash
        int role_id FK
        datetime created_at
    }

    ROLES {
        int role_id PK
        varchar(50) role_name
    }

    TABLES {
        int table_id PK
        int table_number
        varchar(50) status
    }

    SESSION_TYPES {
        int session_type_id PK
        varchar(50) type_name
        decimal(10, 2) price
        int duration_minutes
    }

    SESSIONS {
        int session_id PK
        int table_id FK
        int user_id FK
        int session_type_id FK
        datetime start_time
        datetime end_time
        varchar(50) status
    }

    ORDERS {
        int order_id PK
        int session_id FK
        int user_id FK
        datetime order_time
        varchar(50) status
    }

    ORDER_ITEMS {
        int order_item_id PK
        int order_id FK
        int menu_item_id FK
        int quantity
        decimal(10, 2) price_at_order
    }

    MENU_CATEGORIES {
        int category_id PK
        varchar(100) category_name
    }

    MENU_ITEMS {
        int menu_item_id PK
        varchar(255) item_name
        text description
        decimal(10, 2) price
        int category_id FK
        boolean is_available
        boolean is_vegetarian
        boolean is_vegan
    }

    INVOICES {
        int invoice_id PK
        int session_id FK
        datetime invoice_date
        decimal(10, 2) total_amount
        varchar(50) status
    }

    INVOICE_ITEMS {
        int invoice_item_id PK
        int invoice_id FK
        int menu_item_id FK
        varchar(255) item_name
        int quantity
        decimal(10, 2) price
    }

    PAYMENTS {
        int payment_id PK
        int invoice_id FK
        datetime payment_date
        decimal(10, 2) amount_paid
        varchar(50) payment_method
    }

```

## MySQL Script

The following SQL script can be used to create the database schema in MySQL.

```sql
-- Create the database
CREATE DATABASE IF NOT EXISTS restaurant_qr_db;
USE restaurant_qr_db;

-- Roles Table
CREATE TABLE ROLES (
    role_id INT AUTO_INCREMENT PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE
);

-- Users Table
CREATE TABLE USERS (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES ROLES(role_id)
);

-- Tables Table
CREATE TABLE TABLES (
    table_id INT AUTO_INCREMENT PRIMARY KEY,
    table_number INT NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'Available' -- Available, Occupied
);

-- Session Types Table
CREATE TABLE SESSION_TYPES (
    session_type_id INT AUTO_INCREMENT PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL UNIQUE, -- Buffet, À la carte
    price DECIMAL(10, 2), -- Price for buffet
    duration_minutes INT -- Duration for buffet
);

-- Sessions Table
CREATE TABLE SESSIONS (
    session_id INT AUTO_INCREMENT PRIMARY KEY,
    table_id INT,
    user_id INT, -- User who created the session (Frontdesk)
    session_type_id INT,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    status VARCHAR(50) NOT NULL DEFAULT 'Active', -- Active, Ended, Paid
    FOREIGN KEY (table_id) REFERENCES TABLES(table_id),
    FOREIGN KEY (user_id) REFERENCES USERS(user_id),
    FOREIGN KEY (session_type_id) REFERENCES SESSION_TYPES(session_type_id)
);

-- Menu Categories Table
CREATE TABLE MENU_CATEGORIES (
    category_id INT AUTO_INCREMENT PRIMARY KEY,
    category_name VARCHAR(100) NOT NULL UNIQUE
);

-- Menu Items Table
CREATE TABLE MENU_ITEMS (
    menu_item_id INT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category_id INT,
    is_available BOOLEAN DEFAULT TRUE,
    is_vegetarian BOOLEAN DEFAULT FALSE,
    is_vegan BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (category_id) REFERENCES MENU_CATEGORIES(category_id)
);

-- Orders Table
CREATE TABLE ORDERS (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT,
    user_id INT, -- User who placed the order (Waiter or Customer)
    order_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'Pending', -- Pending, Preparing, Ready, Served, Cancelled
    FOREIGN KEY (session_id) REFERENCES SESSIONS(session_id),
    FOREIGN KEY (user_id) REFERENCES USERS(user_id)
);

-- Order Items Table
CREATE TABLE ORDER_ITEMS (
    order_item_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    menu_item_id INT,
    quantity INT NOT NULL,
    price_at_order DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES ORDERS(order_id),
    FOREIGN KEY (menu_item_id) REFERENCES MENU_ITEMS(menu_item_id)
);

-- Invoices Table
CREATE TABLE INVOICES (
    invoice_id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT,
    invoice_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Unpaid', -- Unpaid, Paid, Partially Paid
    FOREIGN KEY (session_id) REFERENCES SESSIONS(session_id)
);

-- Invoice Items Table
CREATE TABLE INVOICE_ITEMS (
    invoice_item_id INT AUTO_INCREMENT PRIMARY KEY,
    invoice_id INT,
    menu_item_id INT,
    item_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (invoice_id) REFERENCES INVOICES(invoice_id),
    FOREIGN KEY (menu_item_id) REFERENCES MENU_ITEMS(menu_item_id)
);

-- Payments Table
CREATE TABLE PAYMENTS (
    payment_id INT AUTO_INCREMENT PRIMARY KEY,
    invoice_id INT,
    payment_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    amount_paid DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL, -- Cash, Credit Card, etc.
    FOREIGN KEY (invoice_id) REFERENCES INVOICES(invoice_id)
);

-- Insert initial data for Roles
INSERT INTO ROLES (role_name) VALUES ('Admin'), ('Manager'), ('Frontdesk'), ('Waiter'), ('Kitchen'), ('Customer');
```

## Data Models

Here are the corresponding data models for the database schema.

### User
-   **user_id**: Integer (Primary Key)
-   **username**: String
-   **password_hash**: String
-   **role_id**: Integer (Foreign Key to Roles)
-   **created_at**: DateTime

### Role
-   **role_id**: Integer (Primary Key)
-   **role_name**: String

### Table
-   **table_id**: Integer (Primary Key)
-   **table_number**: Integer
-   **status**: String

### SessionType
-   **session_type_id**: Integer (Primary Key)
-   **type_name**: String
-   **price**: Decimal
-   **duration_minutes**: Integer

### Session
-   **session_id**: Integer (Primary Key)
-   **table_id**: Integer (Foreign Key to Tables)
-   **user_id**: Integer (Foreign Key to Users)
-   **session_type_id**: Integer (Foreign Key to SessionTypes)
-   **start_time**: DateTime
-   **end_time**: DateTime
-   **status**: String

### Order
-   **order_id**: Integer (Primary Key)
-   **session_id**: Integer (Foreign Key to Sessions)
-   **user_id**: Integer (Foreign Key to Users)
-   **order_time**: DateTime
-   **status**: String

### OrderItem
-   **order_item_id**: Integer (Primary Key)
-   **order_id**: Integer (Foreign Key to Orders)
-   **menu_item_id**: Integer (Foreign Key to MenuItems)
-   **quantity**: Integer
-   **price_at_order**: Decimal

### MenuCategory
-   **category_id**: Integer (Primary Key)
-   **category_name**: String

### MenuItem
-   **menu_item_id**: Integer (Primary Key)
-   **item_name**: String
-   **description**: String
-   **price**: Decimal
-   **category_id**: Integer (Foreign Key to MenuCategories)
-   **is_available**: Boolean
-   **is_vegetarian**: Boolean
-   **is_vegan**: Boolean

### Invoice
-   **invoice_id**: Integer (Primary Key)
-   **session_id**: Integer (Foreign Key to Sessions)
-   **invoice_date**: DateTime
-   **total_amount**: Decimal
-   **status**: String

### InvoiceItem
-   **invoice_item_id**: Integer (Primary Key)
-   **invoice_id**: Integer (Foreign Key to Invoices)
-   **menu_item_id**: Integer (Foreign Key to MenuItems)
-   **item_name**: String
-   **quantity**: Integer
-   **price**: Decimal

### Payment
-   **payment_id**: Integer (Primary Key)
-   **invoice_id**: Integer (Foreign Key to Invoices)
-   **payment_date**: DateTime
-   **amount_paid**: Decimal
-   **payment_method**: String