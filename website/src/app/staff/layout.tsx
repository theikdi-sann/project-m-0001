"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { supabase } from "@/lib/supabase";
import { Loader2, LayoutDashboard, Utensils, LogOut } from "lucide-react";
import Link from "next/link";

export default function StaffLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkUser = async () => {
      const { data } = await supabase.auth.getUser();
      if (!data.user || data.user.is_anonymous) {
        // Redirect to Staff Login
        router.push("/login");
      } else {
        setLoading(false);
      }
    };
    checkUser();
  }, [router]);

  const handleLogout = async () => {
    await supabase.auth.signOut();
    router.push("/");
  };

  if (loading) return <div className="flex h-screen items-center justify-center"><Loader2 className="animate-spin"/></div>;

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <aside className="w-64 bg-white shadow-md flex flex-col">
        <div className="p-6 border-b">
          <h1 className="text-xl font-bold text-blue-600">Staff Portal</h1>
        </div>
        <nav className="flex-1 p-4 space-y-2">
          <Link href="/staff/kitchen" className="flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 text-gray-700">
            <Utensils size={20} />
            <span>Kitchen Display</span>
          </Link>
          <Link href="/staff/tables" className="flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 text-gray-700">
            <LayoutDashboard size={20} />
            <span>Tables</span>
          </Link>
          <Link href="/staff/menu" className="flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 text-gray-700">
            <Utensils size={20} /> {/* Reusing icon or change to Filter/Menu */}
            <span>Menu Items</span>
          </Link>
        </nav>
        <div className="p-4 border-t">
          <button onClick={handleLogout} className="flex items-center space-x-3 p-3 w-full text-left text-red-600 hover:bg-red-50 rounded-lg">
            <LogOut size={20} />
            <span>Logout</span>
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 overflow-auto p-8">
        {children}
      </main>
    </div>
  );
}
