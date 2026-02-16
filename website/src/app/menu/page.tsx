"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import useSWR from "swr";
import { MenuCategory, MenuItem } from "@/types/menu";
import { useCart } from "@/context/CartContext";
import { Plus, Minus, Loader2, Receipt } from "lucide-react";
import { toast } from "sonner";
import { useLanguage } from "@/context/LanguageContext";
import { LanguageSwitcher } from "@/components/LanguageSwitcher";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

const fetcher = (url: string) => fetch(url).then((res) => res.json());

export default function MenuPage() {
  const router = useRouter();
  const { t } = useLanguage();
  const { addToCart, removeFromCart, getQuantity, totalItems, totalPrice } = useCart();
  const [activeCategory, setActiveCategory] = useState<string | null>(null);

  // 1. Fetch Categories
  const { data: categories, error: catError } = useSWR<MenuCategory[]>(
    `${API_URL}/menu/categories`,
    fetcher
  );

  // 2. Set Default Category
  useEffect(() => {
    if (categories && categories.length > 0 && !activeCategory) {
      setActiveCategory(categories[0].ID);
    }
  }, [categories, activeCategory]);

  // 3. Fetch Items for Active Category
  const { data: items, isLoading: itemsLoading } = useSWR<MenuItem[]>(
    activeCategory ? `${API_URL}/menu/items?category_id=${activeCategory}` : null,
    fetcher
  );

  // Check Session
  useEffect(() => {
    const sessionId = localStorage.getItem("qr_session_id");
    if (!sessionId) {
      router.push("/");
      return;
    }

    // Validate Session Status
    fetch(`${API_URL}/sessions/${sessionId}`)
      .then((res) => {
        if (!res.ok) throw new Error("Session invalid");
        return res.json();
      })
      .then((session) => {
        if (session.Status !== "active") {
          localStorage.removeItem("qr_session_id");
          toast.error(t("session_ended"));
          router.push("/");
        }
        if (session.ExpiresAt && new Date(session.ExpiresAt) < new Date()) {
           localStorage.removeItem("qr_session_id");
           toast.error(t("session_expired"));
           router.push("/");
        }
      })
      .catch(() => {
        localStorage.removeItem("qr_session_id");
        router.push("/");
      });
  }, [router, t]);

  if (catError) return <div className="p-4 text-red-500">Failed to load menu.</div>;
  if (!categories) return <div className="p-4 flex justify-center"><Loader2 className="animate-spin"/></div>;

  return (
    <div className="min-h-screen bg-gray-50 pb-24">
      {/* Header */}
      <div className="bg-white shadow-sm p-4 sticky top-0 z-10">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-xl font-bold text-gray-800">{t("menu")}</h1>
          <div className="flex items-center space-x-2">
            <LanguageSwitcher />
            <button 
              onClick={() => router.push("/orders")}
              className="p-2 bg-gray-100 rounded-full hover:bg-gray-200 text-gray-600"
            >
              <Receipt size={20} />
            </button>
          </div>
        </div>
        {/* Category Tabs */}
        <div className="flex space-x-4 overflow-x-auto pb-2 no-scrollbar">
          {categories.map((cat) => (
            <button
              key={cat.ID}
              onClick={() => setActiveCategory(cat.ID)}
              className={`whitespace-nowrap px-4 py-2 rounded-full text-sm font-medium transition-colors ${
                activeCategory === cat.ID
                  ? "bg-blue-600 text-white shadow-md"
                  : "bg-gray-100 text-gray-600 hover:bg-gray-200"
              }`}
            >
              {cat.Name}
            </button>
          ))}
        </div>
      </div>

      {/* Items Grid */}
      <div className="p-4 grid grid-cols-1 md:grid-cols-2 gap-4">
        {itemsLoading ? (
          <div className="col-span-full text-center py-10"><Loader2 className="animate-spin mx-auto"/></div>
        ) : items?.map((item) => (
          <div key={item.ID} className="bg-white rounded-xl p-4 shadow-sm flex justify-between items-center">
            <div className="flex-1">
              <h3 className="font-semibold text-gray-900">{item.Name}</h3>
              <p className="text-sm text-gray-500 line-clamp-2">{item.Description}</p>
              <div className="text-blue-600 font-bold mt-2">${item.Price}</div>
            </div>
            
            <div className="flex flex-col items-center ml-4 space-y-2">
              {getQuantity(item.ID) > 0 ? (
                <div className="flex flex-col items-center bg-gray-100 rounded-lg p-1">
                  <button onClick={() => addToCart(item)} className="p-1 hover:bg-white rounded"><Plus size={16}/></button>
                  <span className="font-semibold text-sm py-1">{getQuantity(item.ID)}</span>
                  <button onClick={() => removeFromCart(item.ID)} className="p-1 hover:bg-white rounded"><Minus size={16}/></button>
                </div>
              ) : (
                <button
                  onClick={() => addToCart(item)}
                  className="bg-blue-100 text-blue-600 p-2 rounded-full hover:bg-blue-200"
                >
                  <Plus size={20} />
                </button>
              )}
            </div>
          </div>
        ))}
        {items?.length === 0 && <div className="text-center text-gray-500 col-span-full py-10">No items in this category.</div>}
      </div>

      {/* Floating Cart Button */}
      {totalItems > 0 && (
        <div className="fixed bottom-4 left-4 right-4 max-w-md mx-auto">
          <button 
            className="w-full bg-blue-600 text-white p-4 rounded-xl shadow-lg flex justify-between items-center"
            onClick={() => router.push("/cart")} 
          >
            <div className="flex items-center space-x-3">
              <div className="bg-blue-800 px-3 py-1 rounded-full text-xs font-bold">{totalItems}</div>
              <span className="font-medium">{t("view_order")}</span>
            </div>
            <span className="font-bold">${totalPrice.toFixed(2)}</span>
          </button>
        </div>
      )}
    </div>
  );
}
