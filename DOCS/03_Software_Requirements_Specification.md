# Software Requirement Specification (SRS)
## Hybrid Restaurant QR Menu & Management System

**Version:** 1.0
**Architect:** Systems Architect

---

## 1.0 Introduction
### 1.1 Purpose
The purpose of this document is to define the functional and non-functional requirements for the Hybrid Restaurant QR Menu & Management System. This system allows customers to order food via a web interface and enables staff to manage operations via a tablet application.

### 1.2 Scope
The system encompasses a Customer Web App, a Staff Tablet App, a Go-based Logic Server, and a Supabase Database/Realtime engine. It supports Buffet (time-limited) and À la carte dining models.

## 2.0 User Personas
*   **The Diner (Customer):** Tech-savvy or casual diner who wants to order quickly without waiting for a server.
*   **The Waiter (Staff):** Needs to receive instant notifications of new orders and requests (e.g., "Bill please") and manage table status.
*   **The Kitchen Staff:** Needs a clear, real-time display of incoming orders sorted by time and priority.
*   **The Manager:** Needs to oversee active sessions, configure the menu, and handle complex billing scenarios.

## 3.0 Functional Requirements

### 3.1 Authentication & Session Management
*   **FR-01:** The system shall allow customers to **join an active dining session** (previously created by staff) by scanning a QR code, which authenticates them anonymously via a secure token.
*   **FR-02:** The system shall allow staff to log in using an email and password, authenticated against the Supabase Auth provider.
*   **FR-03:** The system shall allow staff to create **Buffet** (time-limited), **À la carte**, or **Take Away** sessions.
*   **FR-04:** The Go backend shall automatically enforce buffet timers and trigger a "Session Ending" alert 10 minutes before expiration.

### 3.2 Menu & Ordering
*   **FR-05:** The system shall display the menu to customers, filtering items based on availability and category.
*   **FR-06:** The system shall allow customers to add items to a cart and submit an order.
*   **FR-07:** The Go backend shall validate every order submission (checking prices, item availability, and session status) before writing to the database.
*   **FR-08:** The system shall broadcast new confirmed orders to the Staff/Kitchen interfaces via WebSockets (Supabase Realtime) within 1 second.

### 3.3 Staff Operations
*   **FR-09:** The Staff Tablet App shall receive reliable push notifications (FCM) for new orders even when the app is in the background.
*   **FR-10:** The system shall allow staff to update the status of an order (Pending -> Preparing -> Served).
*   **FR-11:** The system shall allow staff to manually toggle menu item availability (In Stock / Out of Stock) in real-time.

### 3.4 Billing & Printing
*   **FR-12:** The system shall generate a unique, public URL for a read-only digital e-receipt (PDF or Web View) that customers can download.
*   **FR-13:** The system shall allow staff to trigger a print job to network-connected thermal printers (ESC/POS) via direct TCP connection from the tablet.
*   **FR-14:** The system shall support "Split Bill" calculations for À la carte sessions.

## 4.0 Non-Functional Requirements

### 4.1 Performance
*   **NFR-01:** The Customer Web App shall load the menu (First Contentful Paint) in under **1.5 seconds** on a 4G network.
*   **NFR-02:** API response time for order submission shall be less than **200ms**.

### 4.2 Security
*   **NFR-03:** All sensitive logic (pricing, total calculation) must be executed on the Go backend, never trusted from the client.
*   **NFR-04:** All data transmission shall be encrypted via TLS 1.2 or higher.
*   **NFR-05:** Database access shall be restricted via Row Level Security (RLS) policies in PostgreSQL.

### 4.3 Scalability & Reliability
*   **NFR-06:** The system shall support at least 50 concurrent active sessions on the specified VPS hardware (4 vCPU / 8GB RAM).
*   **NFR-07:** The system shall automatically attempt to reconnect to WebSockets in case of network interruption.

## 5.0 External Interfaces
*   **Supabase:** Used for PostgreSQL Database, Authentication, and Realtime subscriptions.
*   **Stripe API:** Used for processing secure credit card payments.
*   **Firebase Cloud Messaging (FCM):** Used for waking up staff tablets with high-priority alerts.
*   **ESC/POS Printers:** Interfaced via raw TCP streams for ticket printing.
