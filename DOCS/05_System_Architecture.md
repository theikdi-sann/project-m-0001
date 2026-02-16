Unified System Architecture & Business Logic Specification
Hybrid Restaurant QR Menu & Management System
1. Technology Stack & Architecture
The system utilizes a high-performance, cost-effective hybrid stack managed by a solo developer, with a targeted infrastructure cost of under $45 USD/month.

Customer Frontend: Responsive Next.js web application for a "no-app-download" QR scanning and ordering experience.

Staff Frontend: Capacitor-wrapped application (reusing the web codebase) deployed on Android/iOS tablets.

Backend Logic Server: Go (Golang) server responsible for complex business logic, price validation, and enforcing timers.

Database & Realtime: Self-hosted Supabase instance (via Docker) providing PostgreSQL, Authentication, and WebSocket Realtime capabilities.

Hardware Integration: Direct TCP/IP connections to ESC/POS thermal printers for ticketing.

Push Notifications (FCM): Firebase Cloud Messaging is used to forcefully "wake up" staff tablets to ensure alerts are received immediately, even if the device screen is off or the app is in the background.

2. User Roles & Hierarchy
The system enforces strict Role-Based Access Control (RBAC) via Supabase Auth, following this hierarchy: Admin > Manager > Frontdesk > Waiter > Kitchen > Customer.

Customer (Diner): Authenticated anonymously via secure token upon scanning the QR code; can view the menu, place orders, and request the bill.

Kitchen: Views orders in real-time, updates preparation status, and can cancel orders with specific reasons (e.g., out of stock).

Waiter: Receives instant notifications (via FCM), manages order status, requests session extensions/ends, and handles customer service.

Frontdesk: Creates sessions (Buffet, À la carte, or Take Away), assigns tables, generates invoices, and processes payments.

Manager/Admin: Oversees active sessions, resolves overrides, manages staff, and accesses historical analytics and reporting.

3. Core Business Logic
3.1 Session Management
Initialization: Frontdesk staff create sessions and explicitly define the session type as Buffet, À la carte, or Take Away at creation.

Session Lifecycle & Statuses:
*   **Active:** Session created, dining in progress. (Table: Occupied).
*   **Expired:** Buffet time limit reached. Staff alerted. Guests still present. (Table: Occupied).
*   **Cancelled:** Session voided before payment. (Table: Available).
*   **Completed:** Payment successful, guests left. (Table: Available).

Buffet Rules (Dine-in): Features a fixed duration, with timers automatically enforced by the Go backend. The backend triggers a "Session Ending" alert 10 minutes before expiration. Extensions can be requested by waiters but must be executed by Frontdesk or higher. Requires a physical table assignment (Available -> Occupied).

À La Carte Rules (Dine-in): No time limits exist; the session ends either by staff decision or when the customer requests the bill. Requires a physical table assignment (Available -> Occupied).

Take Away Rules (To-Go): Logged as a distinct session type for accurate historical reporting. No physical table is assigned (the system tracks it via a generated "Order Number" or "Queue ID"). There are no time limits, and the session is marked as "Ended" immediately once the food is handed over to the customer.

3.2 Order Management & Notifications
Placement & Validation: Orders can be placed via the QR self-service interface or by waiters. The Go backend validates every submission before committing it to the database.

Real-time Sync & Alerts: Confirmed orders are broadcasted to Kitchen/Staff interfaces within 1 second via Supabase Realtime WebSockets. Simultaneously, FCM push notifications are dispatched to ensure idle tablets ring and alert the staff immediately.

Lifecycle & Modifications: Order status flows from Pending -> Preparing -> Ready -> Served (or "Handed Over" for Take Away). Customers can only modify orders before the kitchen starts preparation.

3.3 Menu Management
Structure: Hierarchical categories with dietary tags and manual real-time availability toggles for items (In Stock / Out of Stock).

Pricing Restrictions: Price changes made in the admin panel do not affect currently active orders.

3.4 Billing & Invoicing
Invoice Generation: * Buffet: Invoices are calculated by customer count × the fixed tier price (e.g., standard, premium, children).

À la carte: Invoices are generated based on actual ordered items when the bill is requested. Split bills are strictly supported for À la carte sessions.

Take Away: Invoices are generated based on actual ordered items. Unlike dine-in, invoices are typically generated immediately at the Frontdesk upon order placement or food collection.

Digital & Physical Receipts: The system generates a public URL for a digital e-receipt and allows staff to trigger physical prints to local network thermal printers.

Payment: Handled via Stripe API for credit cards or via manual cash verification, moving the session status to "Paid".

3.5 Reporting & Historical Analytics
Data Vault: Supabase (PostgreSQL) acts as the permanent system of record, storing all historical data securely.

Metrics Tracked: The system tracks financial transactions, payment history, session lifecycles (start, duration, end time), order cancellations, customer volume, and staff performance metrics.

Categorized Revenue Analysis: Historical data can be strictly filtered by the three distinct session types (Buffet, À la carte, and Take Away), allowing management to perfectly analyze the profitability and volume of walk-in diners versus to-go customers.

4. Project Scope & Delivery
Timeline: The Minimum Viable Product (MVP) is scheduled for delivery within 8 weeks using a phased Agile approach.

Performance Targets: Menu load times under 1.5 seconds, and Go API response times under 200ms.

Security: Sensitive logic is strictly isolated to the Go backend, data is TLS encrypted, and database access is secured using Supabase Row Level Security (RLS).
