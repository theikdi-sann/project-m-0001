# Use Case: Extend Session

## Basic Information
- **Use Case ID**: UC-503
- **Primary Actor**: Frontdesk, Waiter (via request)
- **Secondary Actors**: System, Notification Service
- **Priority**: Medium
- **Status**: Draft
- **Trigger**: Customer requests more time or session is approaching time limit

## Flow of Events

### Preconditions
1. User is authenticated and logged into the system
2. An active session exists that can be extended
3. User has appropriate permissions to extend sessions
4. Session has not exceeded maximum allowable duration

### Postconditions
1. Session duration is extended by the requested time
2. Session end time is updated in the system
3. Any automatic billing calculations are adjusted
4. Audit log records the extension
5. Relevant staff are notified of the extension

### Main Success Scenario

#### Frontdesk-Initiated Extension
1. Frontdesk navigates to the session management section
2. System displays active sessions with current durations
3. Frontdesk selects the session that needs extension
4. System displays session details including current end time
5. Frontdesk clicks "Extend Session" button
6. System shows extension options (30 minutes, 60 minutes, 90 minutes, custom)
7. Frontdesk selects the desired extension duration
8. Frontdesk may add a reason for extension (optional)
9. Frontdesk confirms the extension
10. System validates the extension:
    - New end time doesn't conflict with reservations
    - Session hasn't reached maximum duration limit
    - Extension is within business hours
11. System updates session end time
12. System recalculates any time-based charges if applicable
13. System displays success message: "Session extended by [duration]"
14. System logs the extension in audit trail

#### Waiter-Initiated Extension (Request)
1. Waiter views active sessions (all tables)
2. Waiter selects session that needs extension
3. System shows session details with "Request Extension" button
4. Waiter clicks "Request Extension"
5. System shows extension duration options
6. Waiter selects desired extension and may add customer reason
7. Waiter confirms the request
8. System sends real-time notification to frontdesk with:
    - Table number and session details
    - Requested extension duration
    - Waiter name and request time
    - Customer reason if provided
9. System shows waiter: "Extension request sent to frontdesk"
10. Frontdesk receives notification and processes the request
11. Use case continues from step 5 of Frontdesk-Initiated Extension

### Alternative Flows

#### A1: Maximum Duration Reached
10a. System detects session would exceed maximum duration
11a. System displays: "Cannot extend beyond maximum [X] hours. Please close session."
12a. Use case ends

#### A2: Table Reservation Conflict
10b. System detects table reservation at new end time
11b. System displays: "Table reserved at [time]. Maximum extension: [Y] minutes."
12b. Frontdesk adjusts extension duration or contacts reservation holder
13b. Use case resumes at step 6

#### A3: Business Hours Limitation
10c. System detects extension would go beyond business hours
11c. System displays: "Extension beyond business hours not allowed."
12c. Use case ends

#### A4: Frontdesk Denies Waiter Request
10d. Frontdesk reviews waiter extension request and denies it
11d. System notifies waiter: "Extension request denied for Table [X]"
12d. System logs the denial with reason
13d. Use case ends

### Exception Flows

#### E1: Session Already Ended
5e. System detects session is already closed
6e. System displays: "Cannot extend closed session."
7e. Use case ends

#### E2: System Update Error
11f. System cannot update session end time
12f. System displays: "System error. Extension not saved. Please try again."
13f. Use case ends unsuccessfully

## Special Requirements
- Response time: Extension processing should complete within 2 seconds
- Notifications: Real-time alerts to relevant staff
- Validation: Multiple business rule checks before allowing extension
- Audit: Complete tracking of all extension activities

## Business Rules

### Duration Rules
- BR-SESSION-14: Maximum session duration is 4 hours for à la carte, 2 hours for buffet
- BR-SESSION-15: Extensions can only be granted during business hours
- BR-SESSION-16: Minimum extension increment is 15 minutes
- BR-SESSION-17: Maximum single extension is 90 minutes

### Permission Rules
- BR-SESSION-18: Frontdesk can directly extend any session
- BR-SESSION-19: Waiters can request extensions for any active session (all tables)
- BR-SESSION-20: Managers can override extension limits with reason
- BR-SESSION-21: Extensions cannot conflict with existing reservations

### Notification Rules
- BR-SESSION-22: Waiters must be notified when their extension requests are processed
- BR-SESSION-23: Multiple extensions require manager approval
