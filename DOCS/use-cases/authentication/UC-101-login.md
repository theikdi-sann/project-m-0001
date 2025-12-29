# Use Case: Login

## Basic Information
- **Use Case ID**: UC-101
- **Primary Actor**: Staff Members (Frontdesk, Waiter, Kitchen, Manager, Admin)
- **Secondary Actors**: System, Authentication Service
- **Priority**: High
- **Status**: Draft
- **Trigger**: Staff member needs to access the system

## Flow of Events

### Preconditions
1. Staff member has valid system credentials (username/password)
2. System is operational and accessible
3. Staff member is on login page
4. Staff member's account is active

### Postconditions
1. Staff member is authenticated and granted system access
2. Staff member session is established with appropriate permissions
3. Staff member is redirected to role-specific dashboard
4. Audit log records login attempt
5. Session timer starts based on role-specific timeout rules

### Main Success Scenario
1. Staff member navigates to login page
2. System displays login form with username and password fields
3. Staff member enters:
   - Username
   - Password
4. Staff member clicks "Login" button
5. System validates credentials against user database
6. System verifies user account is active
7. System creates user session with:
   - Staff member ID and role
   - Role-based permissions
   - Session timeout based on role
   - Timestamp and session ID
8. System redirects user to role-specific dashboard:
   - Admin → System Administration Dashboard
   - Manager → Business Operations Dashboard
   - Frontdesk → Session Management Dashboard
   - Waiter → Order & Customer Service Dashboard
   - Kitchen → Kitchen Order Management 
9. System displays welcome message: "Welcome back, [Username]"
10. System logs successful login in audit trail

### Alternative Flows

#### A1: Invalid Credentials
* 5a. System detects invalid username/password combination
* 6a. System displays error: "Invalid username or password"
* 7a. System increments failed login counter for username
* 8a. Use case resumes at step 2
<!-- 
#### A2: Account Locked
* 6b. System detects user account is locked (too many failed attempts)
* 7b. System displays error: "Account temporarily locked. Please contact administrator."
* 8b. Use case ends -->

#### A3: Password Expired
* 6c. System detects temporary password has expired
* 7c. System redirects to password change page
* 8c. Staff member must change password before proceeding
* 9c. Use case continues with password change workflow

# TO REMOVE
#### A4: Role-based Session Timeout
* 7d. System applies role-specific session timeout:
   - Admin/Manager: 8 hours of inactivity

   - Frontdesk: 8 hours of inactivity

   - Waiter: 3 hours of inactivity (shared terminals)

   - Kitchen Staff: 12 hours of inactivity (stationary devices)
#### A5: Role-based Redirection
* 8e. System redirects based on user role:
   - Admin → /admin/dashboard
   - Manager → /manager/dashboard  
   - Frontdesk → /frontdesk/sessions
   - Waiter → /waiter/orders
   - Kitchen → /kitchen/orders

### Exception Flows

#### E1: System Unavailable
* 5f. Authentication service is unavailable
* 6f. System displays: "System temporarily unavailable. Please try again later."
* 7f. Use case ends

#### E2: Database Connection Error
* 5g. System cannot connect to user database
* 6g. System displays: "Database connection error. Contact administrator."
* 7g. Use case ends

## Special Requirements
- Response time: Authentication should complete within 2 seconds
- Security: Passwords must be encrypted in transit and at rest
# TO REMOVE
- Session Management: Role-specific timeout rules must be enforced
- Concurrent users: Support 50+ simultaneous logins

## Business Rules
- BR-AUTH-01: Maximum 5 failed login attempts before temporary lockout
- BR-AUTH-02: Passwords must meet complexity requirements (minimum 9 characters)
- BR-AUTH-03: Role-specific session timeout must be enforced 
- BR-AUTH-04: Different dashboards and permissions for different staff roles

## Security Rules
- BR-AUTH-05: Session tokens must be securely generated and validated

- BR-AUTH-06: All login attempts (successful and failed) must be logged

- BR-AUTH-07: Passwords must be stored using secure hashing algorithms

- BR-AUTH-08: Session data must be cleared upon logout

## User Interface Requirements
Login form should be simple, clean, and intuitive

Show system version and copyright information in footer

Include "Forgot Password" link for password recovery

Display restaurant branding and logo

Responsive design for different device sizes

Clear error messages with appropriate styling

Loading indicator during authentication process

## Accessibility Requirements
Form fields must have proper labels for screen readers

Keyboard navigation support (Tab to navigate, Enter to submit)

High contrast mode support

Error messages must be announced to screen readers

Form should be usable without a mouse

## Integration Points
User Management System: Validates credentials and role permissions

Audit Logging System: Records all login attempts and sessions

Session Management: Creates and manages user sessions

Role-Based Access Control: Determines dashboard redirects and permissions

## Notes
Customer access is handled separately through QR code system (UC-103)

Session timeouts are role-specific to balance security and usability

Failed login attempts are tracked per username to prevent brute force attacks

All authentication events are logged for security auditing purposes

