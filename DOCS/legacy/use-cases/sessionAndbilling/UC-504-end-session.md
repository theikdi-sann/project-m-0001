# Use Case: End Session

## Basic Information
- **Use Case ID**: UC-505
- **Primary Actor**: Frontdesk, Waiter (via request)
- **Secondary Actors**: System, Payment System, Notification Service
- **Priority**: High
- **Status**: Draft
- **Trigger**: Customers are leaving or session needs to be closed

## Flow of Events

### Preconditions
0. User is authenticated and logged into the system
1. An active session exists that can be closed
2. User has appropriate permissions to end sessions
3. All orders are either served or cancelled

### Postconditions
0. Session status is changed to "ended"
1. Final bill is generated and processed if not already done
2. Table status is updated to "available"
3. Session data is archived for reporting
4. Audit log records the session closure

### Main Success Scenario

#### Frontdesk-Initiated Session End
0. Frontdesk navigates to session management
1. System displays active sessions
2. Frontdesk selects session to end
3. System displays session summary:
   - Final bill amount
   - Payment status
   - Outstanding orders
   - Session duration
4. Frontdesk clicks "End Session" button
5. System validates session can be closed:
   - All served orders are billed
   - No pending food orders
   - Payment is processed or waived
6. System displays confirmation with final details
7. Frontdesk confirms session closure
8. System:
   - Updates session status to "ended"
   - Records end time and final duration
   - Updates table status to "available"
   - Archives session data
   - Generates session closure report
9. System displays success message: "Session ended successfully"
10. System notifies cleaning staff if table needs resetting

#### Waiter-Initiated Session End (Request)
0. Waiter views active sessions (all tables)
1. Waiter determines customers are leaving at a specific table
2. Waiter navigates to that session
3. Waiter clicks "Request Session End"
4. System shows confirmation: "Request to end Table [X] session?"
5. Waiter confirms and may add notes (reason, customer count, condition)
6. System sends real-time notification to frontdesk with:
   - Table number and session details
   - Waiter name and request time
   - Final customer count and notes
7. System shows waiter: "Session end request sent to frontdesk"
8. Frontdesk receives notification and processes final billing
9. Use case continues from step 5 of Frontdesk-Initiated Session End

### Alternative Flows

#### A0: Outstanding Orders
5a. System detects unserved orders
6a. System displays: "[X] orders not served. Cancel orders or mark as no-show?"
7a. Frontdesk chooses to cancel orders or mark as no-show
8a. System updates order status accordingly
9a. Use case resumes at step 6

#### A1: Unpaid Balance
5b. System detects unpaid balance
6b. System displays: "Outstanding balance: $[amount]. Process payment first."
7b. Frontdesk processes payment through billing system
8b. Use case resumes at step 6

#### A2: Early Session End
0c. Session is ending significantly before expected time
1c. System prompts for early closure reason
2c. Frontdesk selects reason (customer complaint, emergency, etc.)
3c. System logs reason and proceeds with closure

#### A3: Frontdesk Modifies Waiter Request
8d. Frontdesk reviews waiter end request and modifies details
9d. Frontdesk adjusts final customer count or adds notes
10d. System records both original request and final modifications

### Exception Flows

#### E0: Session Already Ended
4e. System detects session is already closed
5e. System displays: "Session already ended."
6e. Use case ends

#### E1: System Archive Error
8f. System cannot archive session data
9f. System displays warning but continues: "Session ended but archive failed. Contact administrator."
10f. Use case ends with partial success

#### E2: Payment Processing Failure
7g. Payment system is unavailable
8g. System displays: "Payment system offline. Session ended with balance due."
9g. System flags session for follow-up
10g. Use case ends with warning

## Special Requirements
- Response time: Session ending should complete within 2 seconds
- Data integrity: Complete session data must be preserved
- Payment integration: Seamless connection with billing system
- Notifications: Appropriate staff alerted for table turnover

## Business Rules

### Closure Rules
- BR-SESSION-26: Sessions can be ended by authorized staff most of the time System automatically end the sessions
- BR-SESSION-27: All served orders must be paid or accounted for
- BR-SESSION-28: Sessions exceeding maximum duration auto-flag for review
- BR-SESSION-29: Early endings require reason documentation

### Financial Rules
- BR-SESSION-30: Final bill must include all served items
- BR-SESSION-31: Cancelled orders must be documented with reason
- BR-SESSION-32: No-show items require manager approval to remove
- BR-SESSION-33: Session data must be preserved for 7 years for tax purposes

### Waiter Permission Rules
- BR-SESSION-34: Waiters can view all active sessions and orders
- BR-SESSION-35: Waiters can serve orders for any table when they accept them
- BR-SESSION-36: Waiters can request session endings and extending for any table
- BR-SESSION-37: Waiters cannot process payments or generate final bills but they can request

### Operational Rules
- BR-SESSION-38: Table must be marked available for cleaning after session end
- BR-SESSION-39: Waiters receive performance metrics after session closure
- BR-SESSION-40: Session duration affects table turnover calculations
- BR-SESSION-41: Special events or large groups require additional documentation

## User Interface Requirements
- Clear warning messages for outstanding items
- Final confirmation with comprehensive summary
- Easy access to billing and payment functions
- Quick reason selection for exceptions
- Real-time notification system for waiter requests

## Integration Points
- Payment processing system
- Table management system
- Staff notification system
- Reporting and analytics platform
- Inventory management (for consumed items)

---

## Important Notes on Waiter Permissions

### Waiter Session Access:
- Waiters can view ALL active sessions and orders across ALL tables
- Waiters can serve orders for any table when they accept responsibility for them
- Waiters can request session extensions for any table
- Waiters can request session endings for any table
- Waiters cannot directly extend sessions or process payments - these require frontdesk authorization

### Order Serving Logic:
- When a waiter accepts an order, they take responsibility for serving it
- Multiple waiters can serve different orders for the same table
- Waiters can view order status across all tables to coordinate service
- The system tracks which waiter served each order for performance metrics

### Request-Based Workflow:
- All critical session modifications (extension, ending) go through frontdesk
- Waiters initiate requests, frontdesk authorizes and executes
- Real-time notifications ensure quick response times
- Audit trails track both requests and executions

## Relationship Between Use Cases

- **Extend Session** (UC-503) and **End Session** (UC-504) are complementary
- Both support the main **View Session** (UC-503) functionality
- Waiters can initiate both through requests, but only frontdesk can execute
- Business rules ensure consistency between extension limits and ending procedures
- Audit trails maintain complete session lifecycle tracking

