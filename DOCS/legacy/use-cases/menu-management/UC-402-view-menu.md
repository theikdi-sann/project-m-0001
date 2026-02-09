# Use Case: View Menu

## Basic Information
- **Use Case ID**: UC-402
- **Primary Actor**: Customer, Waiter, Frontdesk, Kitchen Staff, Manager, Admin
- **Secondary Actors**: System
- **Priority**: High
- **Status**: Reviewed
- **Trigger**: User needs to browse or search the restaurant menu

## Flow of Events

### Preconditions
1. System is operational and accessible
2. Menu data is available in the system
3. User has appropriate access permissions
4. For customers: Active session exists or browsing allowed
5. Menu items are properly categorized with customizable options

### Postconditions
1. User can view menu items appropriate to their role and session type
2. Menu displays accurate pricing including customizable options
3. User can proceed to order items if authorized
4. Customizable options and pricing are clearly visible
5. Buffet pricing is clearly distinguished from à la carte pricing

### Main Success Scenario

#### Customer View - Buffet Session
1. Customer accesses menu via QR code or table device
2. System detects buffet session and displays **Buffet Menu View**
3. System shows menu items with clear buffet pricing indicators:
    ```
   [ITEM CARD - Buffet VIEW]
   🍜 Pho Bo - $0.00
   Traditional Vietnamese beef noodle soup
   *** Customizable Options Available ***
   [Customize & Add]
   ```
#### Option Selection Process (Modal Approach)
1. Customer clicks "Customize & Add" on any menu item
2. System opens a **full-screen modal** on mobile/tablet
3. Modal displays all option groups with current selections:
   ```
   [OPTION SELECTION MODAL]
   Customize Pho Bo
   ┌─────────────────────────────────────┐
   │ Spice Level (Required)              │
   │ ○ Mild    ○ Medium    ● Spicy    ○ Extra Spicy (+$1) │
   │                                     │
   │ Noodle Type (Required)              │
   │ ○ Rice    ● Egg (+$1.50)    ○ Udon (+$2) │
   │                                     │
   │ Extra Toppings (Optional)           │
   │ ☑ Extra Beef Brisket (+$3)          │
   │ ☐ Extra Meatballs (+$2.50)          │
   │ ☐ Extra Bean Sprouts (+$1)          │
   │                                     │
   │ Total: $19.50                       │
   │ [Add to Order] [Cancel]             │
   └─────────────────────────────────────┘
   ```
4. Customer browses categories or uses search
5. System filters to show only available items
6. Customer can:
   - View item with prcie tag and customizable options include or not tag
   - Click "Customize & Add" to select options and add to order
   - See option prices in customizable option view
   - System adds customized item to order redirect to previous menu page

#### Customer View - À La Carte Session
1. Customer accesses menu via QR code or table device
2. System detects à la carte session and displays **À La Carte Menu View**
3. System shows menu items with à la carte pricing:
   ```
   [ITEM CARD - À LA CARTE VIEW]
   🍜 Pho Bo - $14.00
   Traditional Vietnamese beef noodle soup
   *** Customizable Options Available ***
   [Customize & Add]
   ```

#### Option Selection Process (Modal Approach)
1. Customer clicks "Customize & Add" on any menu item
2. System opens a **full-screen modal** on mobile/tablet
3. Modal displays all option groups with current selections:
   ```
   [OPTION SELECTION MODAL]
   Customize Pho Bo
   ┌─────────────────────────────────────┐
   │ Spice Level (Required)              │
   │ ○ Mild    ○ Medium    ● Spicy    ○ Extra Spicy (+$1) │
   │                                     │
   │ Noodle Type (Required)              │
   │ ○ Rice    ● Egg (+$1.50)    ○ Udon (+$2) │
   │                                     │
   │ Extra Toppings (Optional)           │
   │ ☑ Extra Beef Brisket (+$3)          │
   │ ☐ Extra Meatballs (+$2.50)          │
   │ ☐ Extra Bean Sprouts (+$1)          │
   │                                     │
   │ Total: $19.50                       │
   │ [Add to Order] [Cancel]             │
   └─────────────────────────────────────┘
   ```
4. Customer makes selections with real-time price updates
5. Customer clicks "Add to Order" to confirm
6. System adds customized item to order redirect to previous menu page

#### Staff View (Internal System)
1. Staff member navigates to menu section
2. System displays comprehensive menu view with:
   - All items (including unavailable)
   - Cost and profit margin information
   - Option configuration details
   - Sales performance metrics
   - Management functions for authorized roles

## Alternative Flows

### A1: Buffet Customer Views Extra Charge Items
3a. System shows only items with additional charges in buffet:
   ```
   [EXTRA CHARGE ITEMS]
   🍣 Premium Sushi - +$8.00
   🥩 Grilled Salmon - +$10.00
   🍷 Premium Wine - +$5.00/glass
   ```

### A2: Customer Filters by Dietary Preferences
4b. Customer has dietary restrictions
5b. Customer applies filters:
   - ☑ Vegetarian Only
   - ☑ Gluten-Free
   - ☐ Spicy Items
6b. System shows only matching items with dietary icons

### A3: Quick Add for Simple Items
6c. Customer selects item with no customizable options
7c. System shows quick add button without modal:
   ```
   [SIMPLE ITEM CARD]
   🥤 Vietnamese Iced Coffee - $5.00
   Sweet creamy iced coffee
   [Add to Order]  // No customization needed
   ```

### A4: Menu Pagination
7d. Customer scrolls to bottom of 20-item page
8d. System shows pagination controls:
   ```
   Page 1 of 3 [1] [2] [3] [Next ›]
   Showing 1-20 of 45 items
   ```
9d. Customer clicks next page to load next 20 items

## Exception Flows

### E1: Item Becomes Unavailable
5e. System detects item is no longer available during viewing
6e. System removes item from display or shows "Temporarily Unavailable"
7e. Customer sees updated menu without the item

### E1: Customization Not Available
4f. Customer tries to customize item that kitchen can no longer prepare
5f. System shows: "Customization not available for this item right now"
6f. Offers option to add basic version or choose different item

### E3: Network Connectivity Issues
3g. System cannot load menu data
4g. System shows cached menu with "Last Updated" timestamp
5g. Displays warning: "Menu may not be current"

## Special Requirements
- **Performance**: Load 20 menu items within 2 seconds
- **Mobile Optimization**: Full-screen modals for option selection on mobile
- **Pricing Clarity**: Real-time price updates during customization
- **Buffet Transparency**: Clear included vs extra charge distinction
- **Pagination**: 20 items per page with clear navigation

## Business Rules

### Display Rules
- **BR-MENU-21**: Options and prices visible fater clicking "customize and add"
- **BR-MENU-22**: Buffet customers see clear "EXTRA CHARGE" indicators
- **BR-MENU-23**: Price ranges shown for items with customizable options
- **BR-MENU-24**: Required options must be selected before adding to order
- **BR-MENU-25**: Dietary and allergen icons must be prominent

### Buffet Pricing Rules
- **BR-MENU-26**: Buffet menu shows only items available for buffet
- **BR-MENU-27**: Extra charged items show "Extra Charged" badge prominently
- **BR-MENU-28**: Extra charge items show price clearly

### User Experience Rules
- **BR-MENU-29**: Maximum 20 items loaded per page for performance
- **BR-MENU-30**: Option selection uses options page (route) on mobile
- **BR-MENU-31**: Real-time price calculation during customization
- **BR-MENU-32**: Quick add available for items without options

## User Interface Requirements

### Buffet Pricing Indicators
```
EXTRA CHARGE ITEMS:
[🟡 +$8.00] Badge + Yellow color  
"➕ $8.00 EXTRA" text
```

### Option Display in Main View
```
Options visible immediately:
• Spice Level: Mild/Medium/Spicy/Extra Spicy (+$1)
• Noodle Type: Rice/Egg (+$1.50)/Udon (+$2)
• Extra Toppings: +$1-$3.50

Price Range Display:
"Estimated Total: $14.00 - $21.00"
```

### Mobile Modal Design
- Full-screen overlay on mobile devices
- Large touch targets for options
- Sticky total price at bottom
- Clear action buttons
- Smooth animations

### Pagination Design
```
Page 1 of 3
← Previous [1] [2] [3] Next →
Showing 1-20 of 45 items
```

## Performance Requirements
- Load initial 20 menu items: ≤ 2 seconds
- Option modal opening: ≤ 1 second  
- Real-time price calculation: ≤ 100ms
- Pagination loading: ≤ 1.5 seconds
- Support 100+ concurrent users

## Integration Points
- **Session Management**: Determines buffet vs à la carte view
- **Menu Management**: Provides item data and customizable options
- **Order Management**: Receives customized orders
- **Inventory System**: Checks item availability in real-time
- **Payment System**: Calculates accurate pricing

## Implementation Notes

### Buffet Menu View Strategy
- **Separate Buffet Menu**: Yes, distinct from à la carte menu
- **Clear Pricing**: "EXTRA $X.XX" badges
- **Upsell**: Extra charge items prominently displayed but clearly marked

### Mobile UX Strategy
- **Full-screen Modals**: Better for complex option selection
- **Immediate Option Visibility**: No hidden costs or surprises
- **Progressive Disclosure**: Show basics first, details on tap
- **Touch-Friendly**: Large buttons, minimal typing

### Performance Strategy
- **Pagination**: 20 items per page balances speed and usability
- **Lazy Loading**: Load images as needed
- **Caching**: Frequently viewed items cached locally
- **Progressive Enhancement**: Basic functionality without JavaScript
