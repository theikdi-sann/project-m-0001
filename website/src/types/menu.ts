export interface MenuCategory {
  ID: string;
  Name: string;
  SortOrder: number;
}

export interface MenuItem {
  ID: string;
  CategoryID: string;
  Name: string;
  Description?: string;
  Price: string; // Decimal comes as string from backend usually
  ImageURL?: string;
  IsAvailable: boolean;
}

export interface CartItem {
  menuItem: MenuItem;
  quantity: number;
  notes?: string;
}
