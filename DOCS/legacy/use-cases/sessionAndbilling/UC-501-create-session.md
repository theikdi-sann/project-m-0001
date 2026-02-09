# Use Case: Create Session

## Basic Information
- **Use Case ID**: UC-501
- **Primary Actor**: Frontdesk
- **Secondary Actors**: System, Customer, Waiter (notification)
- **Priority**: High
- **Status**: Draft
- **Trigger**: Customer arrives at restaurant and needs a table

## Flow of Events

### Preconditions
1. Frontdesk staff is authenticated and logged into the system
2. System is operational and accessible
3. At least one table is available in the system
4. Frontdesk has appropriate permissions to create sessions

### Postconditions
1. New session is created and active
2. Table status is updated to "occupied"
3. Session timer starts automatically
4. For buffet sessions: invoice is ready to generate upon session end
5. For à la carte sessions: orders can be placed immediately
6. QR code is generated for customer menu access
7. Waiter may be assigned (if applicable)

### Main Success Scenario

#### Session Creation Process
1. Frontdesk navigates to Session Management dashboard
2. System displays:
   - Available tables with status indicators
   - Current active sessions
   - Table capacity and location information
3. Frontdesk clicks "Create New Session" button
4. System displays session creation form with:

   **Session Configuration:**
   - Table Selection: [Dropdown of available tables]
   - Session Type: ○ Buffet  ○ À La Carte
   - Number of Customers: [Input field, 1-20]
   - Customer Notes: [Optional, 200 characters max]

   **If Buffet Session Selected:**
   - Buffet Tier: [Dropdown - Standard, Premium, VIP]
   - Buffet Tier Base Price: [Auto-displayed based on selection]
   - Total Buffet Charge: [Auto-calculated: base price × customer count]

   **If À La Carte Session Selected:**
   - No upfront pricing (billed per order)

5. Frontdesk selects:
   - Table number
   - Session type (Buffet/À La Carte)
   - Number of customers
   - Additional configuration based on session type
6. Frontdesk clicks "Create Session" button
7. System validates:
   - Table is still available
   - Customer count is valid (1-20)
   - All required fields are populated
   - Buffet tier is selected if buffet session
8. System creates session record with:
   - Unique Session ID
   - Table number
   - Session type
   - Start timestamp
   - Customer count
   - Buffet tier (if applicable)
   - Status: "Active"
   - Session duration timer: 00:00:00
9. System updates table status to "occupied"
10. System generates QR code linking to:
    - Digital menu (filtered by session type)
    - Order placement interface
    - Session information
11. **For Buffet Sessions**: System prepares invoice framework (not generated yet, but ready for session end)
12. System displays session confirmation with:
    - Session ID and table number
    - Session type and start time
    - QR code for customer access
    - Customer count
    - Session duration timer
13. System optionally prints session ticket with QR code
14. Frontdesk directs customers to their table

## Alternative Flows

### A1: Table No Longer Available
7a. System detects selected table is no longer available
8a. System displays error: "Table [X] is no longer available. Please select another table."
9a. Use case resumes at step 4 with updated table availability

### A2: Walk-in During Peak Hours - No Tables Available
2b. System shows no available tables
3b. Frontdesk can:
    - Add customers to waiting list
    - Estimate wait time based on active session durations
    - Suggest alternative time or takeaway
4b. Use case continues with waiting list management (separate use case)

### A3: Customer with Special Requirements
5c. Frontdesk adds special requirements in customer notes:
    - Wheelchair accessibility
    - Celebration (birthday/anniversary)
    - Dietary restrictions
    - Special seating preferences
6c. System records notes and may notify relevant staff

### A4: Large Group - Multiple Tables
4d. Frontdesk selects multiple adjacent tables
5d. System creates linked sessions
6d. System marks tables as "combined seating"
7d. Use case continues with single session covering multiple tables

## Exception Flows

### E1: System Database Error
8e. System cannot save session data
9e. System displays: "System error. Session not created. Please try again."
10e. Use case ends unsuccessfully

### E2: Invalid Customer Count
7f. System detects invalid customer count (0, negative, or >20)
8f. System displays: "Customer count must be between 1 and 20."
9f. Use case resumes at step 4

### E3: Printer Offline for QR Code
13g. System cannot print session ticket
14g. System displays: "Ticket printing failed. QR code is available on screen."
15g. Frontdesk can show digital QR code or manually note session details

## Special Requirements
- **Response Time**: Session creation should complete within 3 seconds
- **QR Code Generation**: Must be scannable by standard smartphones
- **Real-time Updates**: Table status must update immediately across all devices
- **Capacity Management**: Enforce maximum customer limits per table
- **Session Timing**: Accurate session duration tracking from creation

## Business Rules

### Session Creation Rules
- **BR-SESSION-01**: Maximum 20 customers per table for safety regulations
- **BR-SESSION-02**: Table must be available at moment of session creation
- **BR-SESSION-03**: Session start time is system timestamp, not manual entry
- **BR-SESSION-04**: Buffet sessions require tier selection upfront
- **BR-SESSION-05**: À la carte sessions have no upfront pricing

### Buffet Session Rules
- **BR-SESSION-06**: Buffet tier pricing is per person
- **BR-SESSION-07**: Buffet invoice is ready to generate upon session end (not at creation)
- **BR-SESSION-08**: Buffet sessions may have time limits (configurable by tier)
- **BR-SESSION-09**: Buffet customers can order included items at no extra charge

### Table Management Rules
- **BR-SESSION-10**: Table status must update in real-time across all interfaces
- **BR-SESSION-11**: Tables cannot be double-booked
- **BR-SESSION-12**: Table capacity limits must be enforced
- **BR-SESSION-13**: Combined tables require manager approval if exceeding normal capacity

### Customer Experience Rules
- **BR-SESSION-14**: QR code must be generated for all sessions
- **BR-SESSION-15**: Session information must be accessible to customers via QR
- **BR-SESSION-16**: Special requirements must be communicated to relevant staff

## User Interface Requirements

### Session Creation Form
- Visual table layout with availability indicators
- Clear session type toggle (Buffet/À La Carte)
- Real-time customer count validation
- Buffet tier selection with price display
- Customer notes field with character counter
- Session summary preview before confirmation

### Table Status Display
```
Table Layout View:
┌───┬──────────┬──────────┬─────────┐
│ # │ Status   │ Capacity │ Session │
├───┼──────────┼──────────┼─────────┤
│ 1 │ 🟢 Available │ 4       │ -       │
│ 2 │ 🔴 Occupied │ 6       │ S-102   │
│ 3 │ 🟢 Available │ 2       │ -       │
│ 4 │ 🟡 Reserved │ 4       │ R-015   │
└───┴──────────┴──────────┴─────────┘
```

### Session Confirmation Display
```
Session Created Successfully!
──────────────────────────────
Session ID: S-2024-0012
Table: 8 | Type: 🅱️ Buffet (Premium)
Customers: 4 | Start: 14:30
Total: $160.00 (4 × $40.00)
──────────────────────────────
[QR Code Display]
[Print Ticket] [View Session] [Close]
```

## Integration Points

### With Menu System
- **Buffet Sessions**: Filter menu to show included and extra-charge items
- **À La Carte Sessions**: Show full menu with standard pricing
- **QR Code**: Links to session-specific menu view

### With Billing System
- **Buffet Sessions**: Prepares invoice framework for session end billing
- **À La Carte Sessions**: Ready for order-based billing
- **Pricing**: Applies correct buffet tier pricing

### With Table Management
- **Real-time Updates**: Immediate table status changes
- **Capacity Tracking**: Customer count per table
- **Layout Management**: Table combinations and configurations

### With Notification System
- **Kitchen**: New session created (for preparation readiness)
- **Wait Staff**: New table assignment notifications
- **Management**: Session creation analytics

## Data Fields

### Session Record
- Session ID (unique)
- Table number
- Session type (Buffet/À La Carte)
- Start timestamp
- Customer count
- Buffet tier (if applicable)
- Status (Active, Ended, Paid)
- Customer notes
- Created by (staff ID)
- QR code data

### Buffet Session Specific
- Buffet tier
- Base price per person
- Total buffet charge (calculated)
- Time limit (if applicable)

## Notes
- Buffet sessions do NOT generate invoices at creation - invoice is ready to generate upon session end
- À la carte sessions accumulate charges through individual orders
- QR code provides customers with menu access and order placement
- Session timing is critical for buffet time limits and operational efficiency
- Table status management prevents double-booking and ensures accurate capacity planning