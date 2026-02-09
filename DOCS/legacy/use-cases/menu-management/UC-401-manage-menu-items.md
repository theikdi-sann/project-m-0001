# Use Case: Manage Menu Items

## Basic Information
- **Use Case ID**: UC-401
- **Primary Actor**: Admin, Manager
- **Secondary Actors**: System, Kitchen Staff
- **Priority**: High
- **Status**: Reviewed
- **Trigger**: Need to add, update, or remove items from the restaurant menu

## Flow of Events

### Preconditions
1. User is authenticated and logged into the system
2. User has Admin or Manager role permissions
3. System is operational and accessible
4. User navigates to Menu Management section

### Postconditions
1. Menu items are added, updated, or removed as requested
2. Menu changes are reflected across the system (after manual refresh)
3. Kitchen is notified of significant menu changes
4. Audit log records all menu modifications
5. Customers see updated menu on their next refresh

### Main Success Scenario

#### Add New Menu Item
1. User clicks "Add New Item" button
2. System displays comprehensive menu item form with sections:

   **Basic Information:**
   - Item Name (required, unique)
   - Description (optional, 200 characters max)
   - Category (select from predefined list - one level only)
   - Preparation Time (minutes)
   - Chef's Notes (internal, 150 characters max)

   **Pricing & Availability:**
   - À La Carte Price (required)
   - Available for Buffet: ☐ Yes / ☐ No

   **If Available for Buffet Selected:**
   - Buffet Pricing Model:
     - ○ Always Included (free in all buffet types)
     - ○ Always Extra Charge: $______ (extra charge in all buffet types)
     - ○ Tier-Based Pricing (configure per buffet type)

   **If Tier-Based Pricing Selected:**
   ```
   Buffet Tier Pricing Matrix:
   ┌─────────────────┬────────────┬────────────────────┐
   │ Buffet Type     │  Included? │ Extra Charge Price │
   ├─────────────────┼────────────┼────────────────────┤
   │ Standard Buffet │ ☐          │ $ [8.00]          │
   │ Premium Buffet  │ ☐ ✓        │ $ [0.00]          │
   │ VIP Buffet      │ ☐ ✓        │ $ [0.00]          │
   └─────────────────┴────────────┴────────────────────┘
   ```

   **Customizable Options (New Section):**
   - ☐ This item has customizable options
   
   **If Customizable Options Enabled:**
   ```
   Option Groups:
   ┌──────────────────┬─────────────┬──────────┬─────────────┐
   │ Option Group     │ Required?   │ Type     │ Options     │
   ├──────────────────┼─────────────┼──────────┼─────────────┤
   │ Spice Level      │ ☐ Yes       │ Single   │ Mild (+$0)  │
   │                  │             │          │ Medium (+$0)│
   │                  │             │          │ Hot (+$0)   │
   │                  │             │          │ Extra Hot   │
   │                  │             │          │ (+$1)       │
   ├──────────────────┼─────────────┼──────────┼─────────────┤
   │ Noodle Type      │ ☐ Yes       │ Single   │ Rice (+$0)  │
   │                  │             │          │ Egg (+$0)   │
   │                  │             │          │ Udon (+$1)  │
   │                  │             │          │ Soba (+$1)  │
   ├──────────────────┼─────────────┼──────────┼─────────────┤
   │ Doneness         │ ☐ Yes       │ Single   │ Rare (+$0)  │
   │                  │             │          │ Medium (+$0)│
   │                  │             │          │ Well Done   │
   │                  │             │          │ (+$0)       │
   └──────────────────┴─────────────┴──────────┴─────────────┘
   
   [Add Option Group] [Edit Options] [Remove Option Group]
   ```

   **Dietary & Allergen Information:**
   - Dietary & Allergen Tags (multi-select, combined):
     - 🌱 Vegetarian
     - 🌿 Vegan
     - 🌾 Gluten-Free
     - 🥛 Dairy-Free
     - 🥜 Nut-Free
     - ⚡ Spicy
     - 🔥 Extra Spicy
     - 🍃 Healthy Choice
     - 🐟 Contains Fish
     - 🦐 Contains Shellfish
     - 🥚 Contains Eggs
     - 🌰 Contains Nuts
     - 🧴 Contains Dairy
     - 🌾 Contains Gluten
     - 🫘 Contains Soy
     - 🌱 Contains Sesame
     - ☪️ Halal
     - ✡️ Kosher
     - 🍖 Contains Pork

   **Operational:**
   - Status: ☐ Available / ☐ Unavailable
   - Ingredients List (text area, 500 characters max)
   - Kitchen Notes (internal, 200 characters max)
   - Item Image (optional upload)

3. User fills in required information
4. User configures customizable options if needed
5. User clicks "Save Item" button
6. System validates:
   - Item name is unique
   - Required fields are populated
   - À la carte price is provided
   - If available for buffet, buffet pricing model is selected
   - If tier-based pricing, at least one tier must be configured
   - Extra charge prices must be valid numbers
   - Customizable options have valid configuration
7. System creates new menu item with unique ID
8. System displays success message: "Menu item added successfully"
9. System logs creation in audit trail

#### Update Existing Menu Item
1. User searches for menu item using search/filter options
2. System displays list of menu items matching criteria
3. User selects item to edit
4. System displays item details in editable form with current values
5. User modifies information including buffet tier pricing and customizable options
6. User clicks "Save Changes" button
7. System validates changes
8. System updates menu item record
9. If availability status changed to "Unavailable", system:
   - Removes item from customer view
   - Notifies staff of unavailability
10. System displays success message: "Menu item updated successfully"

#### Delete Menu Item
1. User selects menu item from list
2. User clicks "Delete" button
3. System shows confirmation dialog with impact analysis:
   "Delete [Item Name]? This will affect:
   - Remove from all future orders
   - Preserve in historical orders and reports
   - [X] active orders contain this item"
4. User confirms deletion
5. System performs soft delete (archives item)
6. System updates active orders containing the item
7. System displays success message: "Menu item deleted successfully"

## Alternative Flows

### A1: Add Customizable Option Group
1a. User clicks "Add Option Group" button
2a. System displays option group form:
   - Option Group Name (e.g., "Spice Level", "Noodle Type", "Doneness")
   - Required: ☐ Yes / ☐ No
   - Selection Type: ○ Single Choice ○ Multiple Choices
   - Options List (with add/remove):
     - Option Name
     - Additional Price (can be $0)
     - Available: ☐ Yes / ☐ No
3a. User configures the option group
4a. User saves the option group
5a. System adds it to the menu item's customizable options

### A2: Edit Existing Option Group
1b. User clicks "Edit" on an existing option group
2b. System displays the option group in editable form
3b. User modifies options, prices, or availability
4b. User saves changes
5b. System updates the option group configuration

### A3: Bulk Menu Operations
1c. User selects multiple menu items
2c. User chooses bulk action:
   - Change Category
   - Update Prices (percentage or fixed amount)
   - Set Availability Status
   - Apply Dietary Tags
   - Update Buffet Tier Pricing
3c. System shows bulk operation interface
4c. User specifies changes for all selected items
5c. System processes all updates
6c. System shows operation results with success/failure count

## Exception Flows

### E1: Item Name Already Exists
6d. System detects item name already exists
7d. System displays error: "Item name already exists. Please choose a different name."
8d. Use case resumes at step 2

### E2: Invalid Option Configuration
6e. System detects option configuration errors:
   - Option group with no options
   - Required option group with no default selection
   - Invalid additional prices
6e. System displays specific option error message
7e. Use case resumes at option configuration step

### E3: Database Connection Error
7f. System cannot save menu data
8f. System displays: "Database error. Item not saved. Please try again."
9f. Use case ends unsuccessfully

## Special Requirements
- **Response Time**: Menu operations should complete within 3 seconds
- **Option Flexibility**: Support for various option types (single/multiple choice, with/without pricing)
- **Data Validation**: Comprehensive validation for complex option configurations
- **Performance**: Support 500+ menu items with customizable options efficiently

## Business Rules

### Menu Item Rules
- **BR-MENU-01**: Menu item names must be unique within the system
- **BR-MENU-02**: All items must have an à la carte price
- **BR-MENU-03**: Items can be available for both buffet and à la carte
- **BR-MENU-04**: Price changes don't affect active or historical orders

### Buffet Pricing Rules
- **BR-MENU-05**: Buffet tier pricing overrides always included/extra charge settings
- **BR-MENU-06**: If item is included in buffet tier, extra charge price is ignored for that tier
- **BR-MENU-07**: Items must have consistent pricing logic across all service types

### Customizable Options Rules
- **BR-MENU-17**: Option groups can be marked as required or optional
- **BR-MENU-18**: Options can have additional prices that affect the final cost
- **BR-MENU-19**: Option selections must be visible to kitchen staff
- **BR-MENU-20**: Option changes don't affect existing orders

### Availability Rules
- **BR-MENU-08**: Unavailable items are hidden from customer view
- **BR-MENU-09**: Staff can see unavailable items with clear indicators
- **BR-MENU-10**: Items cannot be deleted if in active orders

## User Interface Requirements

### Menu Management Interface
- Simple one-level category selection
- Intuitive option group management
- Drag-and-drop option ordering
- Real-time price calculation preview
- Bulk edit capabilities

### Predefined Categories (One Level Only)
```
- Appetizers
- Soups & Salads
- Main Courses
- Noodles & Rice
- Desserts
- Beverages
- Chef's Specials
- Kids Menu
```

### Common Option Group Templates
```
Spice Level: [Mild, Medium, Hot, Extra Hot]
Doneness: [Rare, Medium Rare, Medium, Medium Well, Well Done]
Noodle Type: [Rice Noodle, Egg Noodle, Udon, Soba]
Sauce Choice: [Black Bean, Oyster, Soy, Teriyaki]
Toppings: [Extra Cheese, Extra Meat, Vegetables, Herbs]
Temperature: [Hot, Cold, Room Temperature]
```

## Database Structure
```sql
menu_items
- id, name, description, category, a_la_carte_price, 
  available_for_buffet, buffet_pricing_model, always_extra_charge_price,
  status, preparation_time, image_url, has_customizable_options

buffet_tiers
- id, tier_name, base_price, status

buffet_item_pricing
- id, menu_item_id, buffet_tier_id, included, extra_charge_price

option_groups
- id, menu_item_id, group_name, required, selection_type, display_order

option_choices
- id, option_group_id, choice_name, additional_price, available
```

## Integration Points
- **Order Management**: Handles customizable options in order creation
- **Kitchen Display**: Shows option selections clearly for preparation
- **Billing System**: Calculates final prices including option charges
- **Customer Facing**: Displays option selection interface

## Notes
- One-level categories simplify menu organization
- Customizable options handle complex item variations like spice levels, doneness, noodle types
- Option pricing is additive to base item price
- Kitchen receives clear instructions about option selections
- Options can be required or optional based on business needs