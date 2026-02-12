-- Enable Realtime for Orders
alter publication supabase_realtime add table orders;
alter publication supabase_realtime add table order_items;
