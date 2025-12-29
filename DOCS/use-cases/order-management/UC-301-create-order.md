# Use Case: Create Order

## Basic Information
- **Use Case ID**: UC-301
- **Primary Actor**: Waiter, Frontdesk, Customer (via QR code)
- **Secondary Actors**: System, Kitchen Staff, Customer
- **Priority**: High
- **Status**: Draft
- **Trigger**: Need to create a new food or beverage order for an active session

## Flow of Events

### Preconditions
1. User is authenticated and logged into the system (for staff) or has valid QR session access (for customers)
2. Active session exists for the table
3. Session has not ended and is not in "paid" status
4. User has appropriate permissions to create orders for the session
5. Restaurant is within operating hours

### Postconditions
1. New order is created and linked to the session
2. Order is sent to kitchen for preparation (if food items)
3. Order total is calculated and added to session bill
4. Inventory levels are updated if applicable
5. Customer can view order status in real-time
6. Audit log records order creation with timestamp and user

### Main Success Scenario

#### Order Creation by Staff (Waiter/Frontdesk)
1. User navigates to order creation interface and selects the target session/table
2. System displays available menu items with:
   - Real-time availability indicators
   - Pricing information
   - Dietary tags and icons
   - Buffet restrictions (if buffet session)
3. User selects items by:
   - Choosing menu categories
   - Selecting specific menu items
   - Specifying quantities (1-5 per item, maximum 15 items total per order)
   - Adding special instructions (up to 150 characters)
   - Marking dietary requirements (highlighted for kitchen)
4. User reviews order summary with:
   - Item list with quantities
   - Special instructions
   - Total amount
   - Estimated preparation time
5. User confirms order creation
6. System performs validations:
   - Session is active and not ended
   - All selected items are available
   - Quantity limits are respected (max 5 per item, max 15 items per order)
   - Special instructions within character limit
7. System creates order record with:
   - Unique Order ID
   - Session ID reference
   - Items with quantities and modifications
   - Special instructions and dietary tags
   - Creation timestamp
   - Status: "Pending"
   - Staff user ID (if created by staff)
8. System calculates order total
9. System updates session bill with new order amount
10. System notifies kitchen staff of new order
11. System displays order confirmation with Order ID
12. Customer receives notification (if via table device)

#### Order Creation by Customer (QR Code)
1. Customer scans QR code or accesses table device
2. System automatically links to their active session
3. System displays customer menu view with:
   - Only available menu items
   - Buffet restrictions applied (if buffet session)
   - Clear pricing for extra-charge buffet items
4. Customer builds order by:
   - Browsing categories
   - Adding items to cart
   - Setting quantities (1-5 per item, max 15 items total)
   - Adding special instructions (150 character limit)
5. Customer reviews cart and submits order
6. System performs backend validations:
   - Real-time stock availability check
   - Session status validation
   - Quantity limit enforcement
7. If validation fails for any item, system responds: "Cannot order [item name] - insufficient stock"
8. If validation passes, system creates order and notifies kitchen directly
9. System updates session bill
10. Customer sees order confirmation with estimated time

## Alternative Flows

### A1: Item Becomes Unavailable During Ordering
6a. System detects item is no longer available during backend validation
7a. System removes unavailable item from order
8a. System completes order with available items only
9a. System displays warning: "[Item] was unavailable and removed from your order"
10a. Use case continues with reduced order

### A2: Buffet Session with Extra Charge Items
3b. User selects items that are outside buffet inclusion
4b. System clearly marks these items as "Extra Charge"
5b. System shows additional costs in order summary
6b. Order proceeds with mixed pricing (buffet included + extra charges)

### A3: Quantity Adjustment Due to Low Stock
6c. System detects insufficient stock for requested quantity
7c. System automatically reduces quantity to available stock
8c. System completes order with adjusted quantities
9c. System displays: "Only [X] of [item] available. Quantity adjusted."
10c. Kitchen may further adjust if stock changes during preparation

### A4: Kitchen Stock Update During Preparation
11d. Kitchen discovers insufficient stock after order received
12d. Kitchen adjusts order quantities in system
13d. System updates order and notifies customer/wait staff
14d. Session bill is automatically adjusted

## Exception Flows

### E1: Session Ended During Order Creation
6e. System detects session has ended or is paid
7e. System prevents order creation
8e. System displays: "Cannot create order - session has ended"
9e. Use case ends unsuccessfully

### E2: Menu Item No Longer Available
6f. System detects multiple items are unavailable
7f. System cancels entire order creation
8f. System displays: "Selected items no longer available. Please create new order."
9f. Use case returns to step 2

### E3: Quantity Limits Exceeded
6g. System detects order exceeds 15 total items or 5 per item
7g. System rejects order creation
8g. System displays: "Order exceeds quantity limits. Maximum 15 items total and 5 per item."
9g. Use case returns to step 3

### E4: Kitchen System Unavailable
10h. System cannot notify kitchen
11h. System queues order for kitchen notification
12h. System displays: "Order created but kitchen notification delayed"
13h. System retries kitchen notification periodically

## Special Requirements
- **Response Time**: Order creation should complete within 3 seconds
- **Real-time Validation**: Stock and session status must be checked during order submission
- **Quantity Limits**: Strict enforcement of 15 items per order and 5 per item
- **Character Limits**: Special instructions limited to 150 characters
- **Availability Updates**: Menu items marked unavailable when stock reaches zero
- **Kitchen Workflow**: Orders must reach kitchen within 10 seconds of creation

## Business Rules

### Order Creation Rules
- **BR-ORDER-01**: Orders cannot be created for ended or paid sessions
- **BR-ORDER-02**: Maximum 15 items per order and maximum 5 of each item
- **BR-ORDER-03**: Special instructions limited to 150 characters
- **BR-ORDER-04**: Real-time stock validation during order submission
- **BR-ORDER-05**: Orders can be created until session ends, even after bill request

### Session Type Rules
- **BR-ORDER-06**: Buffet sessions may have restricted menu categories
- **BR-ORDER-07**: Extra-charge items in buffet sessions are clearly marked and priced
- **BR-ORDER-08**: À la carte sessions have full menu access
- **BR-ORDER-09**: Customer QR orders go directly to kitchen after stock validation

### Inventory Rules
- **BR-ORDER-10**: Low stock triggers automatic quantity adjustment
- **BR-ORDER-11**: Kitchen can modify orders due to stock issues
- **BR-ORDER-12**: Unavailable items are automatically removed from active orders
- **BR-ORDER-13**: Stock levels updated in real-time during order processing

### Notification Rules
- **BR-ORDER-14**: Kitchen receives immediate notification of new orders
- **BR-ORDER-15**: Dietary requirements and special instructions highlighted to kitchen
- **BR-ORDER-16**: Customers receive order confirmation and status updates

## User Interface Requirements

### Order Creation Interface
- Clear menu categorization with visual indicators
- Real-time stock availability display
- Quantity selectors with clear limits
- Special instructions text area with character counter
- Dietary requirement tagging system
- Order summary with running total
- Buffet vs extra charge item distinction

### Validation Feedback
- Immediate feedback on quantity limits
- Clear error messages for unavailable items
- Visual confirmation of successful order creation
- Warning messages for stock adjustments

### Customer QR Interface
- Simplified menu browsing
- Cart management with quantity controls
- Clear pricing and extra charge indications
- Order status tracking
- Session time remaining display

## Accessibility Requirements
- Screen reader compatible menu navigation
- Keyboard-friendly order creation
- High contrast mode support
- Clear visual hierarchy for menu categories
- Accessible error messages and notifications

## Integration Points
- **Session Management**: Validates session status and links orders
- **Menu Management**: Provides available items and pricing
- **Inventory System**: Checks stock levels and updates quantities
- **Kitchen Display System**: Sends orders for preparation
- **Billing System**: Updates session totals
- **Notification System**: Alerts customers and staff

## Data Fields

### Order Record
- Order ID (unique)
- Session ID (foreign key)
- Created timestamp
- Status (Pending, Preparing, Ready, Served, Cancelled)
- Total amount
- Staff user ID (if applicable)
- Special instructions
- Dietary tags

### Order Items
- Menu item ID
- Quantity (1-5)
- Unit price
- Modifications
- Special notes

## Notes
- Orders can be created after bill request since bill request doesn't automatically end session
- Final bill is generated only when session ends, incorporating all orders
- Kitchen has authority to adjust orders due to stock limitations
- Quantity limits help manage kitchen workflow and prevent bulk ordering issues
- Real-time validation ensures order accuracy despite changing stock levels