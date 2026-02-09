# Use Case: Request Invoice

## Basic Information
- **Use Case ID**: UC-505
- **Primary Actor**: Waiter, Customer (via QR/tablet)
- **Secondary Actors**: System, Frontdesk, Notification Service
- **Priority**: High
- **Status**: Reviewed
- **Trigger**: Customer wants to receive the bill/invoice for their session

## Flow of Events

### Preconditions
1. User is authenticated and logged into the system (for waiter) or has valid session access (for customer)
2. Active session exists that can be billed
3. Session is in billable state (Active or Ended but not Paid)
4. System is operational and accessible
5. At least one order exists in the session

### Postconditions
1. Invoice request is sent and notified to Frontdesk
2. Real-time invoice request notification is created and displayed
3. Session status may be updated to "Awaiting Payment"
4. Audit log records the invoice request with timestamp and requestor
5. Customer receives confirmation that request was sent

### Main Success Scenario

#### Customer Invoice Request (via QR/Tablet)
1. Customer accesses their session via QR code or table device
2. System displays session interface with "Request Bill" button in navigation
3. Customer clicks "Request Bill" button
4. System displays confirmation: "Request bill for your session? Waiter will bring your bill shortly."
5. Customer confirms the request
6. System validates:
   - Session exists and is active or ended but unpaid
   - Table is currently occupied
   - Session has at least one order
   - No existing pending invoice request for this session
7. System creates invoice request with:
   - Session ID and table number
   - Customer count
   - Current bill total
   - Request timestamp
   - Requestor: "Customer"
8. System sends real-time notification to Frontdesk
9. System updates session status to "Bill Requested"
10. System displays confirmation to customer: "Bill request sent. Your waiter will assist you shortly."
11. System logs the request in audit trail

#### Waiter Invoice Request
1. Waiter navigates to Sessions view
2. Waiter selects the target session from their assigned tables
3. Waiter clicks "Request Invoice" button
4. System displays confirmation: "Request invoice for Table [X]? This will notify frontdesk to generate the bill."
5. Waiter confirms the request
6. System performs same validation as customer request (step 6 above)
7. System creates invoice request with requestor: "Waiter [Name]"
8. System sends notification to Frontdesk
9. System updates session status to "Bill Requested"
10. System displays confirmation to waiter: "Invoice request sent to frontdesk"
11. System logs the request in audit trail

#### Frontdesk Receiving Invoice Request
1. Frontdesk receives real-time notification via:
   - Dashboard alert
   - Sound notification (if enabled)
   - Visual highlight in sessions list
2. System displays bill request notification with:
   - "BILL REQUEST: Table {tableNumber} | Session {sessionID}"
   - Customer count and current total
   - Requestor (Customer or Waiter name)
   - Request timestamp
   - Time since request
3. Frontdesk can:
   - Click to view session details
   - Generate invoice immediately
   - Acknowledge the request
4. System marks notification as "Seen" when frontdesk views it

## Alternative Flows

### A1: Customer Requests Bill Then Continues Ordering
6a. System detects new orders added after bill request
7a. System updates the bill request with new total
8a. System notifies frontdesk: "Bill updated - new orders added for Table [X]"
9a. Frontdesk generates updated invoice with all items

### A2: Multiple Rapid Requests
6b. System detects duplicate request within 2 minutes
7b. System ignores duplicate and displays: "Bill request already sent. Frontdesk notified."
8b. Use case ends

### A3: Buffet Session Invoice Request
6c. System detects buffet session type
7c. System includes buffet-specific information:
   - Buffet tier and base price
   - Customer count
   - Included items summary
   - Any extra charge items
8c. Use case continues normally

## Exception Flows

### E1: Session Already Paid
6d. System detects session is already in "Paid" status
7d. System prevents request and displays: "Cannot request bill - session already paid."
8d. Use case ends

### E2: No Orders in Session
6e. System detects session has no orders
7e. System prevents request and displays: "Cannot request bill - no orders in session."
8e. Use case ends

### E3: Table Not Occupied
6f. System detects table is not occupied (session ended or cancelled)
7f. System prevents request and displays: "Cannot request bill - session not active."
8f. Use case ends

### E4: Frontdesk Notification Failed
8g. System cannot deliver notification to frontdesk
9g. System queues notification for retry
10g. System displays warning: "Bill request sent but frontdesk notification delayed."
11g. System retries notification every 30 seconds

## Special Requirements
- **Response Time**: Bill requests should process within 2 seconds
- **Notification Delivery**: Frontdesk should receive requests within 10 seconds
- **Real-time Updates**: Session status should update immediately across all devices
- **Audit Trail**: All bill requests must be logged with complete details
- **Mobile Optimization**: Customer interface must work smoothly on mobile devices

## Business Rules

### Request Validation Rules
- **BR-BILL-01**: Invoice requests only allowed for active or ended but unpaid sessions
- **BR-BILL-02**: Session must have at least one order to request invoice
- **BR-BILL-03**: Table must be occupied at time of request
- **BR-BILL-04**: Duplicate requests within 2 minutes are ignored

### Notification Rules
- **BR-BILL-05**: Frontdesk must receive real-time notifications for bill requests
- **BR-BILL-06**: Notifications must include session details and requestor information
- **BR-BILL-07**: Urgent requests (long wait times) should be highlighted
- **BR-BILL-08**: All bill requests must be audited

### Session State Rules
- **BR-BILL-09**: Session status updates to "Bill Requested" upon successful request
- **BR-BILL-10**: New orders can still be added after bill request (up to session end)
- **BR-BILL-11**: Bill requests don't automatically end the session

## User Interface Requirements

### Customer Bill Request Interface
- Clear "Request Bill" button in session navigation
- Simple confirmation dialog
- Loading indicator during request processing
- Success confirmation message
- Option to cancel request if made in error

### Waiter Bill Request Interface
- "Request Invoice" button in session details
- Quick access from session list view
- Visual indicators for sessions with pending bill requests
- Confirmation of request sent

### Frontdesk Notification Interface
- Prominent alert system for new bill requests
- Session summary in notifications
- Quick actions (View Session, Generate Invoice)
- Time tracking for request response times
- Bulk acknowledgment capabilities

## Integration Points
- **Session Management**: Validates session state and updates status
- **Order Management**: Checks for existing orders and calculates totals
- **Notification System**: Delivers real-time alerts to frontdesk
- **Invoice Generation**: Prepares for the actual invoice creation (separate use case)
- **Audit System**: Logs all bill request activities

## Data Fields

### Bill Request Record
- Request ID (unique)
- Session ID
- Table number
- Request timestamp
- Requestor type (Customer/Waiter)
- Requestor details (Waiter ID or Customer session)
- Current bill total at time of request
- Status (Pending, Processed, Cancelled)

### Notification Payload
- Session details
- Customer count
- Bill total
- Requestor information
- Time elapsed since request
- Urgency indicator