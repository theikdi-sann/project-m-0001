# Use Case: Manage Users

## Basic Information
- **Use Case ID**: UC-201
- **Primary Actor**: Admin
- **Secondary Actors**: System, Email Service
- **Priority**: High
- **Status**: Draft
- **Trigger**: Admin needs to manage user accounts in the system

## Flow of Events

### Preconditions
1. Admin is authenticated and logged into the system
2. Admin has appropriate permissions to manage users
3. System is operational and accessible
4. Admin navigates to User Management section

### Postconditions
1. User account is created/updated/deactivated as requested
2. Audit log records the management action
3. Affected user receives notification if required
4. User permissions are updated in real-time

### Main Success Scenario

#### Create New User
1. Admin clicks "Add New User" button
2. System displays user creation form with fields:
   - Username (required)
   - Email (required)
   - Temporary Password (system)
   - Role (dropdown: Admin, Manager, Frontdesk, Waiter, Kitchen)
   - Full Name (required)
   - Phone Number (required)
   - Status (Active/Inactive)
3. Admin fills in required information
4. Admin selects appropriate role and permissions
5. Admin clicks "Create User" button
6. System validates:
   - Username is unique
   - Email format is valid
   - Email is unique
   - Phone Number is unique
   - Required fields are populated
   - Role is valid
7. System creates new user account
8. System sends welcome email with login credentials
9. System displays success message: "User created successfully"
10. System logs the creation event in audit trail

#### Update Existing User
1. Admin searches for user in user list
2. System displays list of users matching search criteria
3. Admin selects user to edit
4. System displays user details in editable form
5. Admin modifies user information:
   - Role changes
   - Status updates
   - Personal information
   - Permission adjustments
6. Admin clicks "Save Changes" button
7. System validates changes
8. System updates user record
9. System displays success message: "User updated successfully"

#### Deactivate User
1. Admin selects user from user list
2. Admin clicks "Deactivate" button
3. System shows confirmation dialog: "Deactivate user [Username]? They will no longer be able to login."
4. Admin confirms deactivation
5. System:
   - Sets user status to "Inactive"
   - Invalidates any active sessions
   - Preserves user data for audit purposes
6. System displays success message: "User deactivated successfully"

### Alternative Flows

#### A1: Username Already Exists
6a. System detects username already exists
7a. System displays error: "Username already taken. Please choose another."
8a. Use case resumes at step 2

#### A2: Invalid Email Format
6b. System detects invalid email format
7b. System displays error: "Please enter a valid email address"
8b. Use case resumes at step 2

#### A3: :Email Already Exists
6c. System detects email already exists
7c. System displays error: "Email already taken. Please choose another."
8c. Use case resumes at step 2

#### A4: Phone Number Alredy Exists
6d. System detects phone number already exists
7d. System displays error: "Phone number already taken. Please choose another."
8d. Use case resumes at step 2

#### A3: Bulk User Operations
1c. Admin selects multiple users
2c. Admin chooses bulk action (Activate, Deactivate, Change Role)
3c. System shows bulk action confirmation
4c. System processes all selected users
5c. System shows bulk operation results

#### A4: Reset User Password
1d. Admin selects "Reset Password" for a user
2d. Admin set temporary password
3d. System sends password reset email to user
4d. System logs password reset event

### Exception Flows

#### E1: Database Connection Error
7e. System cannot save user data
8e. System displays: "Database error. Changes cannot be saved. Please try again later."
9e. Use case ends unsuccessfully

#### E2: Email Service Unavailable          
9f. System cannot send notification email
10f. System displays warning: "User created but notification email failed. Please notify user manually."
11f. Use case continues

#### E3: Permission Denied
5g. Admin tries to modify another Admin account without sufficient privileges
6g. System displays: "Insufficient permissions to modify admin accounts."
7g. Use case ends

## Special Requirements
- Response time: User operations should complete within 3 seconds
- Security: Password must be encrypted before storage
- Audit: All user management actions must be logged
- Performance: Support 1000+ user accounts

## Business Rules

### User Creation Rules
- BR-USER-01: Username must be unique across the system
- BR-USER-02: Email must be unique and valid format
- BR-USER-03: Passwords must be at least 8 characters
- BR-USER-04: New users are created in "Active" status by default

### Role Management Rules
- BR-USER-05: Only Admin can create other Admin accounts
- BR-USER-06: Admin cannot deactivate their own account
- BR-USER-07: Role changes take effect immediate logout and re-login
- BR-USER-08: Users can only be assigned one primary role

### Permission Rules
- BR-USER-09: Permissions are role-based, not individual
- BR-USER-10: Deactivated users cannot login but data is preserved
- BR-USER-11: User management actions are irreversible without Admin override

## User Interface Requirements

### User List View
- Search and filter capabilities
- Sort by username, role, status, last login
- Bulk action options
- Pagination for large user lists

### User Form
- Clear validation messages
- Role-based permission preview
- Password strength indicator
- Status toggle (Active/Inactive)

### Notifications
- Success/error messages for all operations
- Confirmation dialogs for destructive actions
- Email notifications for account changes

## Data Fields

### Required Fields
- Username (unique, 3-20 characters)
- Email (unique, valid format)
- Role (Admin, Manager, Frontdesk, Waiter, Kitchen)
- Full Name
- Status (Active, Inactive)

### Optional Fields
- Phone Number
- Department
- Notes
- Last Login Date
- Created Date

## Security Considerations
- Passwords are never displayed, only reset
- Session invalidation on role changes
- Audit trail for all user modifications
- Role-based access to user management features