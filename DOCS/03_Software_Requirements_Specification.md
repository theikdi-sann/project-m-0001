# Software Requirement Specification (SRS)
## Hybrid Restaurant QR Menu & Management System

**Version:** 1.2
**Architect:** Systems Architect

---

## 1.0 Introduction
### 1.1 Purpose
The purpose of this document is to define the functional and non-functional requirements for the Hybrid Restaurant QR Menu & Management System. This system allows customers to order food via a web interface and enables staff to manage operations via a tablet application.

### 1.2 Scope
The system encompasses a Customer Web App, a Staff Tablet App, a Go-based Logic Server, and a Supabase Database/Realtime engine. It supports Buffet (time-limited), Take Away, and À la carte dining models.

## 2.0 User Personas & Roles

The system enforces strict Role-Based Access Control (RBAC) divided among six primary personas:
* **Admin:** Has operational oversight, system configuration, and manages staff accounts.
* **Manager:** Oversees operations, views reports, and handles complex scenarios (e.g., billing overrides, refunds).
* **Frontdesk:** Handles dining session creation (Buffet, À la carte, Take Away), table assignments, and final payment/cash drawer processing.
* **Waiter:** Receives instant notifications of new orders/requests, manages table status, and provides customer service.
* **Kitchen Staff:** Uses a Kitchen Display System (KDS) to view real-time incoming orders, update preparation status, and manage ingredient availability.
* **The Diner (Customer):** Views the menu, places orders (Dine-in or Take Away), and requests the bill via the QR web interface.

## 3.0 Functional Requirements

### 3.1 Authentication & Session Management
* **FR-01:** The system shall allow customers to **join an active dining session** (previously created by staff) by scanning a QR code, which authenticates them anonymously via a secure token.
* **FR-02:** The system shall allow staff to log in using an email and password, authenticated against the Supabase Auth provider.
* **FR-03:** The system shall allow staff to create **Buffet** (time-limited), **À la carte**, or **Take Away** sessions.
* **FR-04:** The Go backend shall automatically enforce Buffet timers and trigger a "Session Ending" alert to staff(waiter, frontdesk) 10 minutes before expiration.
* **FR-05:** The Go backend shall monitor active À la carte sessions for inactivity. It shall automatically notify the Waiter/Frontdesk to verify if the session should be extended (keep-alive) or ended (before 15 mins of hidden time limit of 3 hours). If no staff action is taken within 15 minutes of the prompt, the system shall automatically set the session status to "Expired".
* **FR-06:** Take Away sessions shall be strictly initiated by Frontdesk staff. The table assignment shall remain `null`. The Frontdesk can either input the order manually on behalf of the customer or generate a unique Session URL/QR Code to send to the customer via other apps or printed QR slip.
* **FR-07:** For Take Away orders placed via the Customer Web App, the system shall require the customer to select a payment intent (Cash on Pickup or Digital Payment) at checkout for payment processing/cancellation.
* **FR-08:** Upon order submission (by staff or customer), the Go backend shall immediately broadcast the order to the Kitchen Display System (KDS) with a clear "Take Away" tag, regardless of whether payment has been collected yet.
* **FR-09:** The Frontdesk system shall allow staff to view the payment status of active Take Away sessions. If the customer selected Cash, the Frontdesk can manually mark the session as Paid upon collection. If Digital was selected, the system shall automatically update the status upon successful gateway verification. Upon pickup and payment completion, the status transitions to "Completed".

### 3.2 Menu & Ordering
* **FR-10:** The system shall display the menu to customers, filtering items based on availability and category.
* **FR-11:** The system shall allow customers to add items to a cart and submit an order.
* **FR-12:** The Go backend shall validate every order submission (checking prices, item availability, and session status) before writing to the database.
* **FR-13:** The system shall broadcast new confirmed orders to the Staff/Kitchen interfaces via WebSockets (Supabase Realtime) within 1 second.

### 3.3 Staff Operations
* **FR-14:** The Staff Tablet App shall receive reliable push notifications (FCM) for new orders even when the app is in the background.
* **FR-15:** The system shall allow staff to update the status of an order (Pending -> Preparing -> Served).
* **FR-16:** The system shall allow staff to manually toggle menu item availability (In Stock / Out of Stock) in real-time.

### 3.4 Billing & Printing
* **FR-17:** The system shall generate a unique, public URL for a read-only digital e-receipt (PDF or Web View) that customers can download.
* **FR-18:** The system shall allow staff to trigger a print job to network-connected thermal printers (ESC/POS) via direct TCP connection from the tablet.
* **FR-19:** The system shall support "Split Bill" calculations strictly for À la carte sessions.

## 4.0 Non-Functional Requirements

### 4.1 Performance
* **NFR-01:** The Customer Web App shall load the menu (First Contentful Paint) in under **1.5 seconds** on a 4G network.
* **NFR-02:** API response time for order submission shall be less than **200ms**.

### 4.2 Security
* **NFR-03:** All sensitive logic (pricing, total calculation) must be executed on the Go backend, never trusted from the client.
* **NFR-04:** All data transmission shall be encrypted via TLS 1.2 or higher.
* **NFR-05:** Database access shall be restricted via Row Level Security (RLS) policies in PostgreSQL.

### 4.3 Scalability & Reliability
* **NFR-06:** The system shall support at least 50 concurrent active sessions on the specified VPS hardware (4 vCPU / 8GB RAM).
* **NFR-07:** The system shall automatically attempt to reconnect to WebSockets in case of network interruption.

## 5.0 External Interfaces
* **Supabase:** Used for PostgreSQL Database, Authentication, and Realtime subscriptions.
* **Stripe API:** Used for processing secure credit card payments.
* **Firebase Cloud Messaging (FCM):** Used for waking up staff tablets with high-priority alerts.
* **ESC/POS Printers:** Interfaced via raw TCP streams for ticket printing.
