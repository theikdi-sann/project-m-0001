"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { supabase } from "@/lib/supabase";
import { Order } from "@/types/order";
import { MenuItem } from "@/types/menu";
import { ArrowLeft, Clock, CheckCircle, ChefHat, XCircle, Loader2 } from "lucide-react";

export default function OrdersPage() {
  const router = useRouter();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  // Poll Orders
  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const sessionId = localStorage.getItem("qr_session_id");
        if (!sessionId) {
            setLoading(false);
            return;
        }

        const { data } = await supabase.auth.getSession();
        const token = data.session?.access_token;
        if (!token) {
            setLoading(false);
            return;
        }

        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/orders?session_id=${sessionId}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (res.ok) {
          const data = await res.json();
          setOrders(data || []);
        }
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
    const interval = setInterval(fetchOrders, 5000);
    return () => clearInterval(interval);
  }, []);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "pending": return <Clock className="text-yellow-500" />;
      case "preparing": return <ChefHat className="text-blue-500" />;
      case "ready": return <CheckCircle className="text-green-500" />;
      case "served": return <CheckCircle className="text-gray-500" />;
      case "cancelled": return <XCircle className="text-red-500" />;
      default: return <Clock />;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-24">
      <div className="bg-white shadow-sm p-4 sticky top-0 z-10 flex items-center">
        <button onClick={() => router.back()} className="mr-4 p-2 hover:bg-gray-100 rounded-full">
          <ArrowLeft size={20} />
        </button>
        <h1 className="text-xl font-bold text-gray-800">Order Status</h1>
      </div>

      <div className="p-4 space-y-4">
        {loading && <div className="text-center py-10"><Loader2 className="animate-spin mx-auto"/></div>}
        
        {!loading && (!orders || orders.length === 0) && (
          <div className="text-center py-10 text-gray-500">No orders placed yet.</div>
        )}

        {orders?.map((order) => (
          <div key={order.ID} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100">
            <div className="flex justify-between items-center mb-3 border-b pb-2">
              <div className="flex items-center space-x-2">
                {getStatusIcon(order.Status)}
                <span className="font-semibold capitalize">{order.Status}</span>
              </div>
              <span className="text-xs text-gray-400">
                {new Date(order.CreatedAt).toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}
              </span>
            </div>
            
            <div className="space-y-2">
              {order.Items?.map((item) => (
                <div key={item.ID} className="flex justify-between text-sm">
                  <span>
                    <span className="font-bold mr-2">{item.Quantity}x</span>
                    Item ID: {item.MenuItemID.substring(0, 8)}... {/* Placeholder for Name */}
                  </span>
                  {item.Notes && <span className="text-xs text-gray-500 italic">({item.Notes})</span>}
                </div>
              ))}
            </div>
            
            <div className="mt-3 pt-2 border-t flex justify-between font-bold text-sm">
              <span>Total</span>
              <span>${order.TotalAmount}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
