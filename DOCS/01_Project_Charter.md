# Project Charter
## Hybrid Restaurant QR Menu & Management System

**Date:** February 9, 2026
**Project Manager:** Solo Developer
**Sponsor:** [Project Owner]

---

### 1. Project Purpose
The purpose of this project is to develop a cost-effective, hybrid software solution that modernizes the restaurant dining experience. By integrating a web-based customer ordering interface with a native staff tablet application, the system aims to eliminate inefficiencies inherent in paper-based workflows, reduce wait times, and minimize order errors for both buffet and à la carte establishments.

### 2. Measurable Objectives
*   **Delivery:** Launch a fully functional Minimum Viable Product (MVP) within **2 months**.
*   **Cost Efficiency:** Maintain operational infrastructure costs below **$40 USD/month** using a high-performance VPS and open-source self-hosted technologies.
*   **Performance:** Achieve sub-200ms API response times and real-time order synchronization latency under 500ms.
*   **Adoption:** Enable a "no-app-download" experience for customers, ensuring 100% accessibility via standard smartphone cameras.

### 3. Key Stakeholders
*   **Project Sponsor / Lead Developer:** [Your Name] (Responsible for funding, development, and execution).
*   **Restaurant Owners/Managers:** Primary users of the management and analytics features.
*   **Restaurant Staff (Waiters/Kitchen):** Daily users of the native tablet application for operations.
*   **Diners (Customers):** End-users of the QR menu and digital billing system.

### 4. High-Level Scope
**In-Scope:**
*   **Customer Interface:** A responsive Next.js web application for QR code scanning, menu browsing, cart management, and order placement (supporting both Buffet and À la carte modes).
*   **Staff Interface:** A Capacitor-wrapped Android/iOS tablet application for table management, real-time kitchen alerts, and billing.
*   **Backend System:** A centralized Go (Golang) backend for complex business logic (timers, price validation) and a self-hosted Supabase instance (PostgreSQL, Auth, Realtime) for data persistence.
*   **Billing:** Digital e-receipt generation (PDF/Web view) and integration with networked thermal printers (TCP/IP).
*   **Authentication:** Anonymous login for customers; secure Email/Password authentication for staff with Role-Based Access Control (RBAC).

**Out-of-Scope (for MVP):**
*   Native mobile application for customers (Web-only to reduce friction).
*   Integration with third-party delivery platforms (UberEats, GrabFood).
*   Advanced inventory stock tracking (beyond basic "Available/Unavailable" toggles).
*   Customer Loyalty/Rewards program.

### 5. High-Level Risks
*   **Resource Availability:** As a solo developer project, illness or burnout represents a single point of failure that could halt progress.
*   **Hardware Integration:** interfacing with various models of ESC/POS thermal printers via TCP/IP may present unforeseen compatibility challenges.
*   **Network Reliability:** The system's real-time features depend on the local network/internet stability; connectivity drops could disrupt kitchen operations.
*   **Payment Compliance:** Delays in verifying manual cash workflows or integrating local Myanmar payment providers could impact the billing cycle.
