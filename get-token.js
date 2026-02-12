const { createClient } = require('@supabase/supabase-js');

const SUPABASE_URL = 'http://127.0.0.1:54321';
const SUPABASE_KEY = 'sb_publishable_ACJWlzQHlZjBrEguHvfOxg_3BJgxAaH';

const supabase = createClient(SUPABASE_URL, SUPABASE_KEY);

async function main() {
  const email = 'tester_' + Date.now() + '@example.com';
  const password = 'password123';

  console.log('Creating user: ' + email + '...');

  const { data, error } = await supabase.auth.signUp({
    email,
    password,
  });

  if (error) {
    console.error('Error:', error.message);
    return;
  }

  const token = data.session?.access_token;
  const userId = data.user?.id;

  if (token && userId) {
    console.log('SUCCESS!');
    console.log('User ID: ' + userId);
    console.log('JWT Token: ' + token);
    console.log('IMPORTANT: Run this SQL to sync the user to public.users:');
    console.log("INSERT INTO users (id, email, role, full_name) VALUES ('" + userId + "', '" + email + "', 'frontdesk', 'Test User') ON CONFLICT DO NOTHING;");
  } else {
    console.log('User created but no session. Check email confirmation.');
  }
}

main();
