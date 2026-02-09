# Use Case: View Dashboard

## Basic Information
- **Use Case ID**: UC-203
- **Primary Actor**: Admin, Manager, Frontdesk, Waiter, Kitchen Staff
- **Secondary Actors**: System, Analytics Service
- **Priority**: High
- **Status**: Draft
- **Trigger**: User logs into the system or navigates to dashboard
- **Frequency**: Multiple times daily

## Flow of Events

### Preconditions
1. User is authenticated and logged into the system
2. User has appropriate role-based permissions
3. System is operational and data services are available
4. User has navigated to the dashboard page

### Postconditions
1. User views relevant dashboard for their role
2. Real-time data is displayed and regularly updated (may be not real-time)
3. User can take action based on dashboard information
4. User session remains active

### Main Success Scenario

#### 1. Admin Dashboard View
1. Admin logs into system
2. System redirects to Admin Dashboard
3. System displays:
   - **System Overview Section:**
     - Total active sessions
     - System health status
     - Active users count
     - Server performance metrics
   - **Business Metrics Section:**
     - Daily revenue (today vs yesterday)
     - Monthly revenue trends
     - Popular menu items
     - Customer count statistics
   - **User Management Section:**
     - Recent user activities
     - System usage statistics
     - Security alerts (if any)
   - **Quick Actions:**
     - Manage Users
     - Generate Reports
     - System Settings

#### 2. Manager Dashboard View
1. Manager logs into system
2. System redirects to Manager Dashboard
3. System displays:
   - **Operations Overview:**
     - Current occupied tables
     - Active sessions count
     - Kitchen order queue status
     - Staff on duty
   - **Sales & Revenue:**
     - Today's revenue vs target
     - Payment methods breakdown
     - Table turnover rate
     - Average order value
   - **Inventory Alerts:**
     - Low stock items
     - Items needing restock
     - Waste tracking
   - **Staff Performance:**
     - Orders per waiter
     - Table service times
     - Customer feedback scores

#### 3. Frontdesk Dashboard View
1. Frontdesk staff logs into system
2. System redirects to Frontdesk Dashboard
3. System displays:
   - **Table Management Section:**
     - Table status grid (color-coded)
     - Available/occupied tables count
     - Reservation timeline
     - Session durations
   - **Session Overview:**
     - Active sessions list
     - New session creation button
     - Buffet vs À la carte breakdown
   - **Customer Flow:**
     - Customers waiting
     - Recent check-ins
     - Bill requests pending
   - **Quick Actions:**
     - Create New Session
     - View Reservations
     - Process Payment Invoice

#### 4. Waiter Dashboard View
1. Waiter logs into system
2. System redirects to Waiter Dashboard
3. System displays:
   - **Assigned Tables:**
     - Table numbers and status
     - Current orders
     - Customer requests
     - Bill readiness status
   - **Order Management:**
     - Pending orders to serve
     - Order preparation status
     - Special requests highlighted
   - **Notifications:**
     - New table assignments
     - Bill requests from customers
     - Kitchen order updates
   - **Performance:**
     - Tables served today
     - Average service time
     - Customer ratings

#### 5. Kitchen Staff Dashboard View
1. Kitchen staff logs into system
2. System redirects to Kitchen Dashboard
3. System displays:
   - **Order Queue:**
     - Pending orders with timestamps
     - Order priority indicators
     - Special preparation notes
     - Table numbers
   - **Preparation Status:**
     - Orders in progress
     - Completed orders
     - Delayed orders highlighted
   - **Inventory Status:**
     - Critical stock levels
     - Menu items availability
     - Ingredient usage rates
   - **Kitchen Metrics:**
     - Orders completed this shift
     - Average preparation time
     - Backlog status

### Alternative Flows

#### A1: Real-time Data Updates
3a. System automatically refreshes dashboard data every 30 seconds
4a. Updated metrics are displayed without page reload
5a. Visual indicators show data refresh status

#### A2: Role Change Detected
2b. System detects user role has changed since last login
3b. System displays: "Your role has been updated. Redirecting to appropriate dashboard."
4b. System redirects to new role's dashboard

#### A3: First Time Login
2c. System detects user's first login
3c. System displays welcome tour/highlights key features
4c. Use case continues normally

#### A4: Custom Dashboard Layout
3d. User has customized their dashboard layout
4d. System loads and applies user's saved preferences
5d. Widgets are arranged according to user customization

### Exception Flows

#### E1: Data Service Unavailable
3e. System cannot retrieve live data from analytics service
4e. System displays cached data with warning: "Displaying cached data. Live updates unavailable."
5e. Use case continues with limited functionality

#### E2: High System Load
3f. System detects high load and simplifies dashboard
4f. System displays essential metrics only with loading indicators
5f. Use case continues with performance optimization

#### E3: Browser Compatibility Issues
2g. System detects unsupported browser features
3g. System displays simplified dashboard version
4g. Use case continues with basic functionality

## Special Requirements

### Performance Requirements
- Dashboard should load within 2 seconds
- Real-time updates should occur every 30 seconds
- Support 50+ concurrent dashboard users
- Cache frequently accessed data for 1 minute

### Data Accuracy
- Financial data must be accurate to 2 decimal places
- Real-time session data must update within 10 seconds
- Inventory levels must reflect current stock
- Order status must be synchronized across all dashboards

### Security Requirements
- Role-based data access enforcement
- Sensitive financial data visible only to authorized roles
- Audit logging of dashboard access

## Business Rules

### Data Access Rules
- BR-DASH-01: Admin can view all system data and metrics
- BR-DASH-02: Manager can view operational and financial data
- BR-DASH-03: Frontdesk can view table and session data only
- BR-DASH-04: Waiter can view only tables (session) and orders
- BR-DASH-05: Kitchen staff can view order queue and inventory only

### Display Rules
- BR-DASH-06: Revenue data rounded to nearest whole number for non-admin roles
- BR-DASH-07: Real-time indicators for urgent items (red for delays, yellow for warnings)
- BR-DASH-08: Personal performance data visible only to the individual staff member
- BR-DASH-09: Historical data limited to 30 days for non-manager roles

### Notification Rules
- BR-DASH-10: Critical alerts must be prominently displayed
- BR-DASH-11: Notifications should not interrupt current workflow
- BR-DASH-12: Users can mute non-critical notifications

## User Interface Requirements

### Common Elements Across All Dashboards
- Role-specific color scheme and branding
- Consistent navigation menu
- User profile and logout access
- Responsive design for tablet and desktop
- Accessibility compliance (WCAG 2.1)

### Widget Requirements
- Draggable and resizable widgets (where applicable)
- Collapsible sections to reduce clutter
- Export functionality for reports
- Print-friendly versions available

### Visual Design
- Color-coded status indicators:
  - Green: Normal/Good
  - Yellow: Warning/Attention needed
  - Red: Critical/Immediate action required
- Clear typography hierarchy
- Intuitive iconography
- Consistent spacing and alignment

## Data Sources and Refresh Rates

| Data Type | Source | Refresh Rate | Role Access |
|-----------|--------|--------------|-------------|
| Session Data | Session Service | 10 seconds | All roles |
| Order Data | Order Service | 15 seconds | All roles |
| Financial Data | Analytics Service | 5 minutes | Admin, Manager |
| Inventory Data | Inventory Service | 30 minutes | Kitchen, Manager, Admin |
| Staff Performance | HR Service | 1 hour | Individual, Manager, Admin |

## Metrics and KPIs by Role

### Admin Dashboard KPIs
- System uptime percentage
- Active user sessions
- Revenue growth rate
- Customer satisfaction score
- Security incident count

### Manager Dashboard KPIs
- Table occupancy rate
- Average order value
- Staff efficiency metrics
- Inventory turnover rate
- Customer wait times

### Frontdesk Dashboard KPIs
- Tables occupied/available
- Session duration averages
- Bill processing time
- Customer check-in rate

### Waiter Dashboard KPIs
- Tables served per hour
- Order accuracy rate
- Customer feedback scores
- Average service time

### Kitchen Dashboard KPIs
- Orders completed per hour
- Preparation time averages
- Waste percentage
- Stock usage rates