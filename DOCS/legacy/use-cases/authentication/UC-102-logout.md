# Use Case: Logout

## Basic Information
- **Use Case ID**: UC-102
- **Primary Actor**: Any Authenticated Staff Member
- **Secondary Actors**: System
- **Priority**: High
- **Status**: Draft
- **Trigger**: Staff member wants to end their session

## Flow of Events

### Preconditions
1. Staff member is currently authenticated and has an active session
2. Staff member is logged into the system
3. Staff member has appropriate permissions for their role

### Postconditions
1. User session is terminated
2. All temporary data is cleared
3. Audit log records logout event
4. User is redirected to login page
5. System resources are released
6. Sessions cannot be resumed without re-authentication

### Main Success Scenario
1. Staff member clicks "Logout" button/menu option from any page
2. System displays confirmation dialog: "Are you sure you want to logout?"
3. User confirms logout action
4. System invalidates user session token
5. System clears session data and cache
6. SYstem records logout event in audit log
5. System redirects user to login page
6. System displays success message: "You have been logged out successfully"

### Alternative Flows

## TO REMOVE
#### A1: Automatic Logout (Timeout)

- 1a. System detects session timeout based on role-specific rules:

   - Admin/Manager: 8 hours of inactivity

   - Frontdesk: 8 hours of inactivity

   - Waiter: 3 hours of inactivity (shared terminals)

   - Kitchen Staff: 12 hours of inactivity (stationary devices)

- 2a. System automatically initiates logout process
- 3a. System displays: "Session expired due to inactivity"
- 4a. System performs steps 4-6 from main success scenario
- 5a. Logout type recorded as: "Timeout"


#### A2: Logout from Multiple Devices
4b. If user was logged in from multiple devices:
   - System invalidates all sessions for that user
   - System notifies other devices of forced logout
6b. Use case continues normally

#### A3: User Cancels Logout
- 3c. User cancels the logout confirmation dialog
- 4c. System returns user to previous page
- 5c. Use case ends

### Exception Flows

## TO THINK AGAIN
#### E1: Session Already Expired
- 4d. System detects session already invalid/expired
- 5d. System clear local session data
- 6d. System redirects directly to login page
- 7d. System displays: "Your session has expired. Please login again."
- 8d. Use case ends

#CHANGE TEXT
#### E2: Network Issues During Logout
- 4e. System cannot reach authentication service during - logout
- 5e. System clears local session data
- 6e. System redirects to login page with warning: "Logout completed locally. Some services may still be active."
- 7e. Use case ends

#### E3: Database Connection Error
- 4f. System cannot update audit log due to database error
- 5f. System continues with logout process but logs error locally
- 6f. System console: "Logout completed with warning: Audit log update failed."

## Special Requirements
- Response time: Logout should complete within 1 second
- Security: Session tokens must be properly invalidated
- Audit: All logout events must be logged with timestamp and user ID

## TO REMOVE
## Business Rules
- BR-AUTH-01: Sessions automatically expire according role-specific session timeouts
- BR-AUTH-02: All logout events must be audited with type and timestamp
- BR-AUTH-03: Admin can force logout any user

## Security Rules
- BR-AUTH-04: Authentication tokens must be properly invalidated server-side

- BR-AUTH-05: Browser cache must be cleared for sensitive data

- BR-AUTH-06: Redirect to login page must include security headers

- BR-AUTH-07: Concurrent sessions from same user should be terminated

## Data Management Rules
- BR-AUTH-08: Temporary session data must be completely cleared

- BR-AUTH-09: User preferences may be preserved (if explicitly saved)

- BR-AUTH-10: Unsaved work should be cleared without preservation

## User Interface Requirements
- Logout option should be easily accessible from all pages
- Confirmation dialog for logout action
- Redirect to login page with appropriate messages