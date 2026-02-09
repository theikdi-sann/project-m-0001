# Business Logic

## 1. SESSION MANAGEMENT

### Session Types & Rules
- TWO session types: Buffet & À la carte
- Buffet: Auto-generate invoice when session created
- À la carte: Generate invoice only when customer requests bill
- Buffet: Fixed duration, auto-end possible + manual end by staff
- À la carte: No time limits, ends via customer bill request or staff decision
- Waiters can request end session for BOTH types
- Waiters can request extend session for BUFFET only
- Only Frontdesk+ roles can execute session end/extend

### Session End Triggers:
| Session Type | Auto End | Customer Request | Waiter Request | Staff Direct End |
|-------------|----------|------------------|----------------|------------------|
| **Buffet**  | ✅ Yes    | ✅ Via waiter    | ✅ Yes         | ✅ Yes           |
| **À La Carte** | ❌ No  | ✅ Bill request = end request | ✅ Yes | ✅ Yes |

### Session Creation
- Frontdesk creates sessions with table assignment
- Must specify session type (Buffet/À la carte) at creation
- Session timer starts immediately
- Table status: Available → Occupied
## 2. BILLING & INVOICING

### Invoice Generation Triggers:
| Session Type | Customer Bill Request | Auto Session End | Staff Manual End |
|-------------|----------------------|------------------|------------------|
| **Buffet**  | ✅ Generate Invoice  | ✅ Generate Invoice | ✅ Generate Invoice |
| **À La Carte** | ✅ Generate Invoice | ❌ No Auto End | ✅ Generate Invoice |

### Invoice Generation
- Buffet: ready to generate invoice upon session end
- À la carte && Buffet: Generate invoice when customer requests bill
- Buffet invoices: customer count × fixed price (different types: basic/medium/premium/children)
- À la carte invoices: actual ordered items
- Frontdesk generates invoices for à la carte upon request

### Payment Processing
- Buffet: Payment processed at session end (after fixed duration or early end)
- À la carte: Payment processed when bill requested
- Split bills supported for À la carte only
- Single payment method per session (for now)
- Session status flow: Active → Ended → Payment Processing → Paid
## 3. ORDER MANAGEMENT

### Order Creation
- Customers create orders via waiters or QR self-service
- Orders linked to specific sessions
- Kitchen receives orders in real-time
- Order status: Pending → Preparing → Ready → Served

### Order Modifications
- Orders can be modified ONLY BEFORE kitchen starts preparation
- Cannot modify after kitchen starts preparation
- Kitchen can update order status only (not items)
- Kitchen can cancel orders for specific reasons (out of stock, etc.)
- Special dietary requests must be highlighted to kitchen

## 4. USER ROLES & PERMISSIONS

### Role Hierarchy
Admin > Manager > Frontdesk > Waiter > Kitchen > Customer

### Key Permissions
- Admin: Full system access, user management, analytics, system config
- Manager: Operational oversight, reports, staff management, overrides
- Frontdesk: Session management, billing, table assignment, payment processing
- Waiter: Order management, customer service, request actions (end/extend sessions)
- Kitchen: View orders, update preparation status, cancel orders (with reason)
- Customer: View menu, place orders, request bill

## 5. MENU MANAGEMENT

### Menu Structure
- Hierarchical categories (Appetizers, Main Courses, Desserts, Beverages, etc.)
- Items have availability status (Available/Unavailable)
- Dietary tags (Vegetarian, Vegan, Gluten-free, etc.)
- Seasonal items with auto-activation dates
- Different buffet types with different pricing (lunch/dinner/children)

### Menu Access
- Customers see only available items
- Staff see all items with availability indicators
- Menu updates NOT real-time (manual refresh required)
- Price changes don't affect active orders

## 6. TABLE & RESERVATION SYSTEM

### Table States
- Available → Occupied (session created) → Available (session ended)
- Real-time table status updates
- Table assignment during session creation
- Reservation system: Yes (to be implemented)

### Session Duration
- Buffet: Fixed duration, auto-end possible, extension requests allowed
- À la carte: No time limits, customer-controlled via bill request
- Extension requests for Buffet only
- End requests for both session types

## 7. NOTIFICATION SYSTEM

### Real-time Notifications
- Kitchen: New orders and order modifications
- Frontdesk: Bill requests, end session requests, extend session requests
- Waiters: Orders ready for serving
- Customers: Order status updates (real-time)

### Request Workflow
Waiter requests → Frontdesk notification → Frontdesk action → Waiter confirmation

## 8. REPORTING & ANALYTICS

### Data Tracking
- All financial transactions and payment history
- Session lifecycle (start, duration, end, payment time)
- Order modifications and cancellation reasons
- User actions audit trail
- Customer count, timing metrics, and behavior patterns
- Revenue analysis by session type, time, menu items
- Staff performance (orders handled, efficiency, customer ratings)
- Menu analysis (popular items, profitability, seasonal trends)
- Table turnover rates and utilization

## 9. BUSINESS RULES & WORKFLOWS

### Customer Journey

#### Buffet Customer:
Arrival → Session Created → Dining → 
[Request Bill] → Generate Invoice → Payment → Session End

OR

[Auto-end] → Generate Invoice → Payment → Session End

#### À La Carte Customer:
Arrival → Session Created → Dining → 
[Request Bill] → Generate Invoice → Payment → Session End

### Staff Workflow
#### Management: 
- Oversight → Analytics → Reporting → Business decisions

#### Frontdesk:
- Generate invoices for BOTH types on customer request
- Generate invoices for Buffet at auto-end
- Can generate invoices for both types at staff discretion
- Process payments after invoice generation

#### Waiter:
- Can request bill generation for BOTH session types
- Customer service → Order management → Request actions (end/extend sessions)
- Customer bill request → Waiter → Frontdesk → Generate Invoice

#### Kitchen
- Order preparation → Status updates → Order cancellation (with reasons)

### Payment Flow
Buffet: Session created → Session ends OR Request Bill (invoice generated)→ Payment processed
À la carte: Session created → Orders placed → Bill requested (invoice generated) → Session ends → Payment processed

## 10. FUTURE ENHANCEMENTS

### Planned Features
- Reservation system
- Membership discounts and loyalty programs
- Promotional pricing and seasonal offers
- Inventory and stock management
- Multiple payment methods
- Real-time menu updates
- Advanced analytics and predictive reporting

## 11. TECHNICAL CONSTRAINTS

### System Limitations
- Menu updates require manual refresh (not real-time)
- Single payment method per session (initially)
- Order modifications restricted to pre-kitchen phase
- Buffet sessions have fixed duration limits
- Role-based permissions with clear hierarchy
