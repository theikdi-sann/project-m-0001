# Project Management Plan (PMP)
## Hybrid Restaurant QR Menu & Management System

**Project Manager:** Agile Project Manager
**Duration:** 2 Months (8 Weeks)
**Methodology:** Agile (Kanban/Scrum Hybrid)

---

### 1. Work Breakdown Structure (WBS) - Phased Approach

| Phase | Status | Key Deliverables |
| :--- | :--- | :--- |
| **Phase 1: Foundation (Backend)** | ✅ **Complete** | - VPS & Docker (Supabase/Postgres)<br>- Go Backend Clean Arch<br>- Session Management (Buffet Timers)<br>- Menu Management API |
| **Phase 2: Core Logic & Security** | ✅ **Complete** | - Order Management (Validation, State Machine)<br>- **RBAC & Auth** (JWT Middleware, Secure API)<br>- Realtime Order Broadcast (Websockets)<br>- **Take Away Lifecycle** (Nullable Tables, Unlimited Duration) |
| **Phase 3: Frontend Ecosystem** | 🔄 **In Progress** | - **Customer Web App:** Menu, Cart, QR Entry, I18n, Tests (Done)<br>- **Staff App (Tablet):**<br>  - KDS (Kitchen Display) (FR-13, FR-15)<br>  - Table Management & Notifications (FR-14)<br>- **Logic:** À La Carte Inactivity Monitor (FR-05) |
| **Phase 4: Billing & Hardware** | 📅 **Next** | - **Billing:** Split Bill Logic (FR-19), Digital Receipts (FR-17)<br>- **Payment:** Take Away Payment Intent (Cash/Digital) (FR-07, FR-09)<br>- **Hardware:** TCP Thermal Printer Integration (FR-18)<br>- **Reporting:** Revenue Analytics |

### 2. Resource Allocation

| Role | Count | Responsibilities |
| :--- | :--- | :--- |
| **Full Stack Developer** | 1 (Solo) | Frontend (Next.js), Mobile (Capacitor), Backend (Go), DevOps (Docker), Database (SQL) |
| **QA Tester** | 1 (Self) | Manual testing of user flows on actual devices (Phone + Tablet) |

### 3. Communication Plan (Solo Dev Context)

*   **Daily Stand-up (Self):** 10-minute morning review of the Trello/Jira board. "What did I do yesterday? What will I do today? Blockers?"
*   **Weekly Sprint Review:** Friday afternoon. Review code quality, merge feature branches, and update the "ProjectBook.md" with progress.
*   **Stakeholder Update:** Bi-weekly demo to the hypothetical (or real) restaurant owner to verify requirements.

### 4. Milestones

| Milestone ID | Description | Target Date |
| :--- | :--- | :--- |
| **M-01** | **Backend Core:** API, DB, Auth, Realtime ready. | ✅ Done |
| **M-02** | **Customer MVP:** QR Scan -> Order -> Kitchen Receive. | ✅ Done |
| **M-03** | **Staff Operations:** KDS, Table Management, Take Away handling. | Week 5 End |
| **M-04** | **Hardware/Billing:** Printer connected, Payments integrated. | Week 7 End |
| **M-05** | **Enterprise Launch:** Stress tests, Analytics, Production Deploy. | Week 8 End |
