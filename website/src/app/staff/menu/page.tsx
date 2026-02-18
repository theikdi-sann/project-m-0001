"use client";

import { useEffect, useState } from "react";
import { supabase } from "@/lib/supabase";
import { MenuItem, MenuCategory } from "@/types/menu";
import { Loader2, ToggleLeft, ToggleRight, Filter } from "lucide-react";
import { toast } from "sonner";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function StaffMenuPage() {
  const [categories, setCategories] = useState<MenuCategory[]>([]);
  const [activeCategory, setActiveCategory] = useState<string | null>(null);
  const [items, setItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);

  // Initial Load
  useEffect(() => {
    const loadCategories = async () => {
      const { data } = await supabase.auth.getSession();
      if (!data.session) return; // Layout handles redirect

      try {
        const res = await fetch(`${API_URL}/menu/categories`);
        if (res.ok) {
          const cats = await res.json();
          setCategories(cats);
          if (cats.length > 0) setActiveCategory(cats[0].ID);
        }
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    };
    loadCategories();
  }, []);

  // Fetch Items when category changes
  useEffect(() => {
    if (!activeCategory) return;

    const loadItems = async () => {
      const { data } = await supabase.auth.getSession();
      const token = data.session?.access_token;
      if (!token) return;

      try {
        const res = await fetch(`${API_URL}/menu/items?category_id=${activeCategory}&include_unavailable=true`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (res.ok) {
          setItems(await res.json());
        }
      } catch (e) {
        console.error(e);
      }
    };
    loadItems();
  }, [activeCategory]);

  const toggleAvailability = async (item: MenuItem) => {
    const { data } = await supabase.auth.getSession();
    const token = data.session?.access_token;
    if (!token) return;

    const newStatus = !item.IsAvailable;
    
    // Optimistic Update
    setItems(prev => prev.map(i => i.ID === item.ID ? { ...i, IsAvailable: newStatus } : i));

    try {
      const res = await fetch(`${API_URL}/menu/items/${item.ID}/availability`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ is_available: newStatus })
      });
      if (!res.ok) throw new Error("Failed");
      toast.success(`${item.Name} is now ${newStatus ? "Available" : "Unavailable"}`);
    } catch (e) {
      toast.error("Update failed");
      // Revert (fetch again or toggle back)
      setItems(prev => prev.map(i => i.ID === item.ID ? { ...i, IsAvailable: !newStatus } : i));
    }
  };

  if (loading) return <div className="flex h-full items-center justify-center"><Loader2 className="animate-spin"/></div>;

  return (
    <div className="h-full flex flex-col">
      <h1 className="text-2xl font-bold text-gray-800 mb-6 flex items-center">
        <Filter className="mr-2" /> Menu Management
      </h1>

      {/* Categories */}
      <div className="flex space-x-2 overflow-x-auto pb-4 mb-4">
        {categories.map(cat => (
          <button
            key={cat.ID}
            onClick={() => setActiveCategory(cat.ID)}
            className={`px-4 py-2 rounded-lg font-medium transition-colors ${
              activeCategory === cat.ID ? "bg-blue-600 text-white" : "bg-white text-gray-600 hover:bg-gray-50 border"
            }`}
          >
            {cat.Name}
          </button>
        ))}
      </div>

      {/* Items List */}
      <div className="flex-1 overflow-auto bg-white rounded-xl shadow p-4">
        <div className="space-y-2">
          {items.map(item => (
            <div key={item.ID} className={`flex justify-between items-center p-3 rounded-lg border ${!item.IsAvailable ? "bg-gray-50 opacity-75" : "bg-white"}`}>
              <div>
                <h3 className="font-semibold text-gray-900">{item.Name}</h3>
                <p className="text-sm text-gray-500">${item.Price}</p>
              </div>
              
              <button 
                onClick={() => toggleAvailability(item)}
                className={`flex items-center space-x-2 px-4 py-2 rounded-lg font-bold transition-colors ${
                  item.IsAvailable 
                    ? "bg-green-100 text-green-700 hover:bg-green-200" 
                    : "bg-red-100 text-red-700 hover:bg-red-200"
                }`}
              >
                {item.IsAvailable ? <ToggleRight size={24} /> : <ToggleLeft size={24} />}
                <span>{item.IsAvailable ? "In Stock" : "Sold Out"}</span>
              </button>
            </div>
          ))}
          {items.length === 0 && <div className="text-center text-gray-500 py-10">No items found.</div>}
        </div>
      </div>
    </div>
  );
}
