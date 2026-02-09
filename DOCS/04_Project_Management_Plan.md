# Project Management Plan (PMP)
## Hybrid Restaurant QR Menu & Management System

**Project Manager:** Agile Project Manager
**Duration:** 2 Months (8 Weeks)
**Methodology:** Agile (Kanban/Scrum Hybrid)

---

### 1. Work Breakdown Structure (WBS) - Phased Approach

| Phase | Duration | Key Deliverables |
| :--- | :--- | :--- |
| **Phase 1: Foundation** | Weeks 1-2 | - VPS Setup & Docker Configuration (Supabase)<br>- Git Repo & CI/CD Setup<br>- Database Schema Design (Postgres)<br>- Go Backend "Hello World" |
| **Phase 2: Core Development** | Weeks 3-5 | - **Customer Web:** Menu UI, Cart, QR Scanning<br>- **Staff App:** Login, Dashboard Shell<br>- **Backend:** Order Validation Logic, Realtime Websockets |
| **Phase 3: Logic & Integration** | Weeks 6-7 | - **Buffet Logic:** Timer implementation<br>- **Billing:** Digital Receipt generation, Stripe Integration<br>- **Hardware:** Thermal Printer TCP Integration |
| **Phase 4: Polish & Launch** | Week 8 | - **Testing:** End-to-End User Flow testing<br>- **Bug Fixes:** UI/UX refinements<br>- **Deployment:** Production VPS Live |

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
| **M-01** | **Infrastructure Ready:** Supabase running on VPS, accessible via URL. | Week 2 End |
| **M-02** | **"First Order":** A mock customer can place an order and it appears in the DB. | Week 4 End |
| **M-03** | **Realtime Sync:** Kitchen tablet updates instantly when order is placed. | Week 5 End |
| **M-04** | **Hardware Integ:** Receipt prints successfully on thermal printer. | Week 7 End |
| **M-05** | **MVP Launch:** System deployed and ready for live customers. | Week 8 End |
