"use client";

import { useEffect, useState } from "react";
import { supabase } from "@/lib/supabase";
import { Order } from "@/types/order";
import { MenuItem } from "@/types/menu";
import { Loader2, ChefHat, CheckCircle, Clock } from "lucide-react";
import { toast } from "sonner";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

export default function KitchenPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [menuMap, setMenuMap] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);

  const fetchOrders = async () => {
    const { data } = await supabase.auth.getSession();
    const token = data.session?.access_token;
    if (!token) return;

    try {
      const res = await fetch(`${API_URL}/orders/active`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) {
        setOrders(await res.json());
      }
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const init = async () => {
        // 1. Fetch Menu for Names
        try {
            const catRes = await fetch(`${API_URL}/menu/categories`);
            if (catRes.ok) {
                const cats = await catRes.json();
                const map: Record<string, string> = {};
                for (const cat of cats) {
                    const itemRes = await fetch(`${API_URL}/menu/items?category_id=${cat.ID}`);
                    if (itemRes.ok) {
                        const items: MenuItem[] = await itemRes.json();
                        items.forEach(i => map[i.ID] = i.Name);
                    }
                }
                setMenuMap(map);
            }
        } catch (e) { console.error(e); }

        // 2. Fetch Orders
        await fetchOrders();
    };

    init();

    // Realtime Subscription
    const channel = supabase
      .channel('kitchen-kds')
      .on(
        'postgres_changes',
        { event: '*', schema: 'public', table: 'orders' },
        (payload) => {
          console.log('Realtime update:', payload);
          fetchOrders(); 
        }
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, []);

  const updateStatus = async (orderId: string, newStatus: string) => {
    const { data } = await supabase.auth.getSession();
    const token = data.session?.access_token;
    if (!token) return;

    // Optimistic Update
    setOrders(prev => prev.map(o => o.ID === orderId ? { ...o, Status: newStatus as any } : o));

    try {
      const res = await fetch(`${API_URL}/orders/${orderId}/status`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ status: newStatus })
      });
      if (!res.ok) throw new Error("Failed to update");
      toast.success(`Order marked as ${newStatus}`);
    } catch (e) {
      toast.error("Update failed");
      fetchOrders(); // Revert
    }
  };

  if (loading) return <div className="flex h-full items-center justify-center"><Loader2 className="animate-spin"/></div>;

  const pending = orders.filter(o => o.Status === "pending");
  const preparing = orders.filter(o => o.Status === "preparing");
  const ready = orders.filter(o => o.Status === "ready");

  return (
    <div className="grid grid-cols-3 gap-4 h-full">
      {/* Pending Column */}
      <div className="bg-white rounded-xl shadow p-4 flex flex-col">
        <h2 className="text-lg font-bold text-gray-700 mb-4 flex items-center">
          <Clock className="mr-2 text-yellow-500" /> Pending ({pending.length})
        </h2>
        <div className="flex-1 overflow-auto space-y-3">
          {pending.map(order => (
            <OrderCard key={order.ID} order={order} menuMap={menuMap} onAction={() => updateStatus(order.ID, "preparing")} actionLabel="Start Cooking" />
          ))}
        </div>
      </div>

      {/* Preparing Column */}
      <div className="bg-white rounded-xl shadow p-4 flex flex-col">
        <h2 className="text-lg font-bold text-gray-700 mb-4 flex items-center">
          <ChefHat className="mr-2 text-blue-500" /> Preparing ({preparing.length})
        </h2>
        <div className="flex-1 overflow-auto space-y-3">
          {preparing.map(order => (
            <OrderCard key={order.ID} order={order} menuMap={menuMap} onAction={() => updateStatus(order.ID, "ready")} actionLabel="Mark Ready" />
          ))}
        </div>
      </div>

      {/* Ready Column */}
      <div className="bg-white rounded-xl shadow p-4 flex flex-col">
        <h2 className="text-lg font-bold text-gray-700 mb-4 flex items-center">
          <CheckCircle className="mr-2 text-green-500" /> Ready ({ready.length})
        </h2>
        <div className="flex-1 overflow-auto space-y-3">
          {ready.map(order => (
            <OrderCard key={order.ID} order={order} menuMap={menuMap} onAction={() => updateStatus(order.ID, "served")} actionLabel="Served" variant="green" />
          ))}
        </div>
      </div>
    </div>
  );
}

function OrderCard({ order, menuMap, onAction, actionLabel, variant = "blue" }: { order: Order, menuMap: Record<string, string>, onAction: () => void, actionLabel: string, variant?: string }) {
  return (
    <div className="border rounded-lg p-3 hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start mb-2">
        <span className="font-mono text-xs text-gray-500">#{order.ID.substring(0, 6)}</span>
        <span className="text-xs text-gray-400">
          {new Date(order.CreatedAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
        </span>
      </div>
      
      <div className="space-y-1 mb-3">
        {order.Items?.map(item => (
          <div key={item.ID} className="text-sm">
            <span className="font-bold mr-2">{item.Quantity}x</span>
            <span>{menuMap[item.MenuItemID] || "Item " + item.MenuItemID.substring(0, 4)}</span>
            {item.Notes && <div className="text-xs text-red-500 ml-6 italic">{item.Notes}</div>}
          </div>
        ))}
      </div>

      <button 
        onClick={onAction}
        className={`w-full py-2 rounded font-semibold text-sm text-white ${variant === "green" ? "bg-green-600 hover:bg-green-700" : "bg-blue-600 hover:bg-blue-700"}`}
      >
        {actionLabel}
      </button>
    </div>
  );
}
