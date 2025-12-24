# Use Case: Update Order

## Basic Information
- **Use Case ID**: UC-303
- **Primary Actor**: Waiter, Frontdesk, Kitchen Staff, Customer (via QR Menu)
- **Secondary Actors**: System, Notification Service
- **Priority**: High
- **Status**: Draft
- **Trigger**: Need to modify an existing order before kitchen starts preparation

## Flow of Events

### Preconditions
1. User is authenticated and logged into the system (for staff) or has valid QR session access (for customers)
2. Active order exists that can be modified
3. Order is in "Pending" status (no items have started preparation)
4. User has appropriate permissions to update orders
5. Session is active and not ended

### Postconditions
1. Order is updated with requested changes
2. Order total is recalculated and session bill updated
3. Kitchen is notified if food items are affected
4. Customer is notified of changes (if applicable)
5. Audit log records the modification with timestamp and user
6. Original order version is preserved for audit purposes

### Main Success Scenario

#### Update Order Items (Before Kitchen Preparation)
1. User navigates to active orders and selects order to modify
2. System displays order details with clear status indicator: "Modifiable - Pending"
3. User makes changes:
   - Adds new items from available menu (respecting 15 total items, 5 per item limits)
   - Removes existing items
   - Changes quantities of existing items
   - Updates special instructions (150 character limit)
   - Modifies dietary requirements
4. User reviews updated order summary with new total
5. User confirms the changes
6. System validates:
   - Order is still in "Pending" status (no items in preparation)
   - All new items are available
   - Quantity limits are maintained
   - Session is still active
7. System updates order record with:
   - Modified items and quantities
   - Updated special instructions
   - New calculated total
   - Modification timestamp and user
8. System updates session bill with new amount
9. System notifies kitchen of order changes (if food items affected)
10. System displays success message: "Order updated successfully"
11. Customer receives update notification (if via table device)

#### Update Order Status (Kitchen Workflow)
1. Kitchen staff views pending orders
2. Kitchen selects order and clicks "Start Preparation"
3. System updates order status to "Preparing"
4. System locks entire order from further modifications
5. System notifies waiter: "Order #[ID] is now being prepared"
6. Customer view updates to show preparation status

#### Kitchen Item Cancellation
1. Kitchen discovers issue with specific item (out of stock, quality problem, etc.)
2. Kitchen selects the specific item and clicks "Cancel Item"
3. System prompts for cancellation reason:
   - Out of stock
   - Quality issue
   - Equipment failure
   - Customer request (via kitchen)
   - Other (with explanation)
4. Kitchen selects reason and confirms cancellation
5. System:
   - Removes item from order
   - Updates order total and session bill
   - Records cancellation with reason and staff ID
   - Notifies waiter and customer of cancelled item
6. Order continues with remaining items

## Alternative Flows

### A1: Order No Longer Modifiable (Kitchen Started Preparation)
- 6a. System detects order status is no longer "Pending"
- 7a. System prevents modifications and displays: "Cannot modify order - kitchen has started preparation"
- 8a. User can:
   - View current order status
   - Contact kitchen for emergency changes
   - Create a new order for additional items
- 9a. Use case ends for modification attempts

### A2: Item Becomes Unavailable During Update
- 6b. System detects selected new item is no longer available
- 7b. System removes unavailable item from update
- 8b. System completes update with available items only
- 9b. System displays warning: "[Item] was unavailable and removed from your update"

### A3: Customer QR Modification Request
- 1c. Customer requests modification via QR code
- 2c. System checks order modifiability (must be in "Pending" status)
- 3c. If modifiable, customer makes changes directly
- 4c. If not modifiable, system displays: "Order is being prepared and cannot be modified"
- 5c. Use case continues based on modifiability status

### A4: Quantity Limit Enforcement During Update
- 6d. System detects update would exceed 15 total items or 5 per item
- 7d. System rejects the update
- 8d. System displays: "Update would exceed quantity limits. Maximum 15 items total and 5 per item."
- 9d. Use case returns to modification step

## Exception Flows

### E1: Session Ended During Update
- 6e. System detects session has ended
- 7e. System prevents order update
- 8e. System displays: "Cannot update order - session has ended"
- 9e. Use case ends

### E2: Concurrent Modification Attempt
- 6f. System detects another user is modifying the same order
- 7f. System displays: "Another user is currently modifying - this order. Please try again."
- 8f. Use case returns to order list

### E3: Database Update Error
- 7g. System cannot save order changes due to database error
- 8g. System displays: "System error. Order not updated. Please try again."
- 9g. Use case ends unsuccessfully

## Special Requirements
- **Response Time**: Order updates should complete within 3 seconds
- **Modification Lock**: Entire order locks when ANY item moves to "Preparing" status
- **Real-time Sync**: All users see updated order information within 10 seconds
- **Audit Trail**: Complete history of all order modifications preserved
- **Notification Speed**: Kitchen notifications within 15 seconds of changes

## Business Rules

### Modification Rules
- **BR-ORDER-30**: Orders can only be modified when ALL items are in "Pending" status
- **BR-ORDER-31**: Once ANY item moves to "Preparing", entire order is locked from modifications
- **BR-ORDER-32**: Quantity limits (15 total items, 5 per item) apply during updates
- **BR-ORDER-33**: Special instructions limited to 150 characters during updates
- **BR-ORDER-34**: Customers can modify their own orders via QR when order is modifiable

### Kitchen Rules
- **BR-ORDER-35**: Kitchen can cancel specific items, not entire orders
- **BR-ORDER-36**: Item cancellations require mandatory reason selection
- **BR-ORDER-37**: Kitchen status updates trigger order modification locks
- **BR-ORDER-38**: Cancelled items are removed from billing automatically

### Notification Rules
- **BR-ORDER-39**: Kitchen must be notified of all order modifications affecting food items
- **BR-ORDER-40**: Customers must be notified when their orders are modified
- **BR-ORDER-41**: Waiters must be notified when kitchen cancels items
- **BR-ORDER-42**: All modification events must be audited

### Data Integrity Rules
- **BR-ORDER-43**: Original order version preserved for audit and reporting
- **BR-ORDER-44**: Modification timestamps and user IDs recorded for all changes
- **BR-ORDER-45**: Order totals recalculated automatically after modifications
- **BR-ORDER-46**: Session bills updated in real-time after order changes

## User Interface Requirements

### Modification Interface
- Clear visual indicator of order modifiability status
- Warning messages when approaching quantity limits
- Real-time total calculation during modification
- Easy add/remove item functionality
- Character counter for special instructions
- Confirmation dialog for significant changes

### Kitchen Interface
- Prominent "Start Preparation" button that locks orders
- Easy item cancellation with reason dropdown
- Clear display of dietary requirements and special instructions
- Urgent order highlighting

## Accessibility Requirements
- Screen reader compatible modification interface
- Keyboard navigation for all update functions
- High contrast mode support
- Clear error messages with suggested actions
- Accessible confirmation dialogs

## Integration Points
- **Order Management**: Provides order data and status tracking
- **Menu Management**: Validates item availability during updates
- **Kitchen Display**: Receives order updates and status changes
- **Billing System**: Updates session totals after modifications
- **Notification System**: Alerts relevant users of changes
- **Audit System**: Records all modification events

## Data Security
- Role-based permissions for order modifications
- Audit trail of all changes with user identification
- Data validation to prevent unauthorized modifications
- Session integrity checks during updates
- Concurrent modification protection

## Notes
- The entire order becomes unmodifiable once any item moves to "Preparing" status
- Kitchen can only cancel specific items with documented reasons
- Customers can modify their own orders via QR without waiter approval, but system enforces modifiability rules
- All modifications are tracked for reporting and customer service purposes
- Order locking prevents conflicts between kitchen preparation and customer modifications