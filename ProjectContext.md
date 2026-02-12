# PROJECT CONTEXT: Hybrid Restaurant QR Management System
I am building a comprehensive restaurant management system as a Solo Developer.

## 1. Core Concept
* Product: A hybrid QR ordering system for restaurants in Myanmar.
* Users:
    * Customers: Web-based (No install). Scan QR -> Order (Dine-in/Takeaway) -> Pay.
    * Staff (Kitchen/Waiters): Native Android Tablet App (via Capacitor). Receives real-time orders, manages tables, and handles billing.
* Unique Logic: Supports Buffet Sessions (Fixed Time Limit vs. Unlimited) and Digital Receipts.

## 2. Technical Stack (Strict Constraints)
* Infrastructure: Self-hosted on a $20-40 VPS (4 vCPU / 8GB RAM).
* Backend: Go (Golang) + Standard Library / pgx.
    * *Architecture:* Clean Architecture (Handler -> Usecase -> Repository).
    * *Testing:* TDD using testify (Mocking interfaces).
* Database & Auth: Supabase (Self-Hosted via Docker).
    * *DB:* PostgreSQL (managed by Supabase).
    * *Realtime:* Supabase Realtime (WebSockets) for Kitchen Display.
* Frontend: Next.js (Static Export).
    * *Styling:* Tailwind CSS.
    * *Mobile Wrapper:* Capacitor (for Staff App to handle Background Mode, Wake Lock, and TCP Thermal Printing).
* Deployment: Docker Compose + Nginx (Reverse Proxy) + Systemd (for Go binary).

## 3. Current Progress (Phase 1: Foundation)
* Completed:
    * Domain entities (`DiningSession`, `Order`) defined.
    * Database Schema (Tables: users, dining_tables, dining_sessions, `orders`) created via SQL migration.
    * Repository Interface defined.
* Current Task: Implementing the Create Dining Session Usecase.
    * *Logic:* Check if table is occupied -> If not, create session.
    * *Method:* TDD (Mocking the Repository first).

## 4. Development Rules for AI
1.  Solo Dev Mode: Solutions must be simple and maintainable. Avoid over-engineering (e.g., no microservices, no Kubernetes).
2.  TDD First: Always write the Test before the Implementation. Use testify/mock.
3.  Supabase CLI: We use supabase start for local dev, NOT manual Docker Compose.
4.  Clean Architecture: Keep business logic pure in usecase/. Database code lives strictly in repository/.
