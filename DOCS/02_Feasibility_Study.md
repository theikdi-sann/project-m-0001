# Feasibility Study Report (FSR)
## Hybrid Restaurant QR Menu & Management System

**Prepared By:** Technical Business Analyst
**Date:** February 9, 2026

---

### 1. Executive Summary
This report evaluates the feasibility of developing a Hybrid Restaurant QR Menu & Management System by a solo developer within a 2-month timeline. The proposed solution leverages a modern, high-performance tech stack (Next.js, Go, Supabase) hosted on a VPS. **Recommendation: GO.**

### 2. Technical Feasibility
*   **Tech Stack Suitability:**
    *   **Frontend (Next.js + Tailwind):** Excellent for creating a fast, responsive web interface for customers and can be easily wrapped with Capacitor for the staff tablet app, allowing for maximum code reuse (single codebase).
    *   **Backend (Go):** Go is statically typed and compiled, offering superior performance for concurrent tasks like managing multiple buffet timers and handling high-volume real-time WebSocket connections. It is an ideal choice for the "secure middleware" layer.
    *   **Database/Infra (Supabase on Docker):** Self-hosting Supabase provides an enterprise-grade feature set (Auth, Realtime, Postgres) without the high costs of managed services.
*   **Challenge:** Configuring and maintaining a self-hosted Supabase instance via Docker requires DevOps knowledge.
*   **Conclusion:** The stack is technically sound and highly optimized for performance and maintainability by a skilled solo developer.

### 3. Economic Feasibility (Cost-Benefit Analysis)
*   **Estimated Costs (Monthly):**
    *   **VPS (Hetzner/DigitalOcean, 8GB RAM):** $20 - $40 USD.
    *   **Domain Name:** ~$1 USD (amortized).
    *   **Software Licenses:** $0 (Open Source / Free Tier).
    *   **Total OpEx:** **<$45 USD/month**.
*   **Development Costs:**
    *   Primarily "sweat equity" (Solo Developer time).
*   **Benefits:**
    *   Reduction in paper printing costs for menus.
    *   Increased table turnover due to faster ordering.
    *   Reduction in order errors (saving food cost).
    *   Hardware savings by utilizing existing Android tablets and printers.
*   **Conclusion:** The project has extremely low overhead and high potential ROI.

### 4. Operational Feasibility
*   **Customer Adoption:** The "no-app-download" web-based approach removes the biggest barrier to entry. Customers are already habituated to scanning QR codes.
*   **Staff Adoption:** Moving from paper to tablets requires training. The UI/UX must be intuitive. The "Wake-on-LAN" style notifications (FCM) ensure staff don't miss orders even if the device is idle.
*   **Conclusion:** Operational success depends heavily on UX design and staff training, but no structural barriers exist.

### 5. Legal & Compliance Feasibility
*   **Data Privacy (GDPR/General):** The system uses **Anonymous Auth** for customers, meaning no Personally Identifiable Information (PII) like names or emails is strictly required to place an order, significantly reducing compliance burden.
*   **Payments (PCI-DSS):** Credit card processing is offloaded to **Stripe**, which handles PCI-DSS compliance. The app never touches raw card data.
*   **Local Compliance (Myanmar):** The "Manual Cash Verification" workflow aligns with local business practices and avoids complex regulatory hurdles for the MVP.
*   **Conclusion:** The project is legally viable with the current architecture.

### 6. Recommendation
**GO.** The project is technically feasible, economically attractive, and operationally viable. The risks are manageable through careful scope control and rigorous testing.
