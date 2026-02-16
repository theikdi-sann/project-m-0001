import { MenuItem } from "./menu";

export type OrderStatus = "pending" | "preparing" | "ready" | "served" | "cancelled";

export interface Order {
  ID: string;
  DiningSessionID: string;
  Status: OrderStatus;
  TotalAmount: string;
  CreatedAt: string;
  Items: OrderItem[]; // Note: My Repo ListBySessionID currently DOES NOT return items (impl detail).
}

export interface OrderItem {
  ID: string;
  MenuItemID: string;
  Quantity: number;
  Notes?: string;
  // TODO: Expand this in Backend to include Item Name?
}
