const { createClient } = require('@supabase/supabase-js');

const SUPABASE_URL = 'http://127.0.0.1:54321';
const SUPABASE_KEY = 'sb_publishable_ACJWlzQHlZjBrEguHvfOxg_3BJgxAaH';

const supabase = createClient(SUPABASE_URL, SUPABASE_KEY);

async function seed() {
  console.log('Seeding Data...');

  const email = 'staff_admin@test.com';
  const password = 'password123';
  
  let userId;
  const { data: authData, error: authError } = await supabase.auth.signUp({ email, password });
  
  if (authError) {
    const { data: loginData, error: loginError } = await supabase.auth.signInWithPassword({ email, password });
    if (loginData.user) {
        userId = loginData.user.id;
        console.log('User already exists: ' + userId);
    } else {
        console.error('Auth Error:', authError.message);
        return;
    }
  } else {
    userId = authData.user?.id;
    console.log('Created User: ' + userId);
  }

  if (!userId) return;

  const sql = 
"-- Users\n" +
"INSERT INTO users (id, email, role, full_name) VALUES ('" + userId + "', '" + email + "', 'frontdesk', 'Admin Staff') ON CONFLICT (id) DO NOTHING;\n" +
"\n" +
"-- Tables\n" +
"INSERT INTO dining_tables (id, table_number) VALUES ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'T-101') ON CONFLICT DO NOTHING;\n" +
"INSERT INTO dining_tables (id, table_number) VALUES ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'T-102') ON CONFLICT DO NOTHING;\n" +
"INSERT INTO dining_tables (id, table_number) VALUES ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a77', 'T-103') ON CONFLICT DO NOTHING;\n" +
"\n" +
"-- Session Types\n" +
"INSERT INTO dining_session_types (id, name, is_buffet, price, duration_minutes) VALUES ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Lunch Buffet', true, 19.99, 90) ON CONFLICT DO NOTHING;\n" +
"INSERT INTO dining_session_types (id, name, is_buffet, price, duration_minutes) VALUES ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a88', 'Standard Take Away', false, 0.00, 240) ON CONFLICT DO NOTHING;\n" +
"INSERT INTO dining_session_types (id, name, is_buffet, price, duration_minutes) VALUES ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a99', 'À La Carte', false, 0.00, 240) ON CONFLICT DO NOTHING;\n" +
"\n" +
"-- Menu\n" +
"INSERT INTO menu_categories (id, name, sort_order) VALUES ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'Main Course', 1) ON CONFLICT DO NOTHING;\n" +
"INSERT INTO menu_items (id, category_id, name, price, is_available) VALUES ('e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'Cheeseburger', 10.00, true) ON CONFLICT DO NOTHING;\n";

  console.log('\nRun this SQL:');
  console.log(sql);
}

seed();
