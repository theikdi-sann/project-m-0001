"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useCart } from "@/context/CartContext";
import { supabase } from "@/lib/supabase";
import { ArrowLeft, Trash2, Loader2, CheckCircle } from "lucide-react";
import { toast } from "sonner";

export default function CartPage() {
  const router = useRouter();
  const { cart, removeFromCart, clearCart, totalPrice } = useCart();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [orderSuccess, setOrderSuccess] = useState(false);

  const placeOrder = async () => {
    setIsSubmitting(true);
    try {
      const sessionId = localStorage.getItem("qr_session_id");
      if (!sessionId) throw new Error("No active session");

      const { data: sessionData } = await supabase.auth.getSession();
      const token = sessionData.session?.access_token;

      if (!token) throw new Error("Authentication failed");

      const items = cart.map((item) => ({
        menu_item_id: item.menuItem.ID,
        quantity: item.quantity,
        notes: item.notes || "",
      }));

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/orders`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          dining_session_id: sessionId,
          items,
        }),
      });

      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || "Failed to place order");
      }

      clearCart();
      setOrderSuccess(true);
      toast.success("Order placed successfully!");
      // Optional: Redirect to Order Status page after 2s
      setTimeout(() => router.push("/orders"), 2000);

    } catch (err: any) {
      toast.error(err.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (orderSuccess) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-4 text-center text-green-600">
        <CheckCircle className="w-16 h-16 mb-4" />
        <h2 className="text-2xl font-bold">Order Placed!</h2>
        <p className="text-gray-600 mt-2">The kitchen has received your order.</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-24">
      {/* Header */}
      <div className="bg-white shadow-sm p-4 sticky top-0 z-10 flex items-center">
        <button onClick={() => router.back()} className="mr-4 p-2 hover:bg-gray-100 rounded-full">
          <ArrowLeft size={20} />
        </button>
        <h1 className="text-xl font-bold text-gray-800">Your Order</h1>
      </div>

      <div className="p-4 space-y-4">
        {cart.length === 0 ? (
          <div className="text-center py-10 text-gray-500">Your cart is empty.</div>
        ) : (
          cart.map((item) => (
            <div key={item.menuItem.ID} className="bg-white p-4 rounded-xl shadow-sm">
              <div className="flex justify-between items-start">
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="bg-blue-100 text-blue-800 text-xs font-bold px-2 py-1 rounded">
                      {item.quantity}x
                    </span>
                    <h3 className="font-semibold text-gray-900">{item.menuItem.Name}</h3>
                  </div>
                  <div className="text-gray-500 text-sm mt-1">${item.menuItem.Price}</div>
                </div>
                <button 
                  onClick={() => removeFromCart(item.menuItem.ID)}
                  className="text-red-400 hover:text-red-600 p-2"
                >
                  <Trash2 size={18} />
                </button>
              </div>
              {/* Notes Input */}
              {/* Ideally update context, but skipping for MVP brevity */}
            </div>
          ))
        )}
      </div>

      {cart.length > 0 && (
        <div className="fixed bottom-0 left-0 right-0 bg-white border-t p-4 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.1)]">
          <div className="flex justify-between mb-4 text-lg font-bold">
            <span>Total</span>
            <span>${totalPrice.toFixed(2)}</span>
          </div>
          <button
            onClick={placeOrder}
            disabled={isSubmitting}
            className="w-full bg-green-600 text-white p-4 rounded-xl shadow-lg font-bold flex justify-center items-center disabled:opacity-50"
          >
            {isSubmitting ? <Loader2 className="animate-spin" /> : "Place Order"}
          </button>
        </div>
      )}
    </div>
  );
}
