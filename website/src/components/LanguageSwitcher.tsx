"use client";

import { useLanguage } from "@/context/LanguageContext";
import { Globe } from "lucide-react";

export function LanguageSwitcher() {
  const { language, setLanguage } = useLanguage();

  return (
    <button
      onClick={() => setLanguage(language === "en" ? "my" : "en")}
      className="p-2 rounded-full hover:bg-gray-100 flex items-center space-x-1"
      aria-label="Switch Language"
    >
      <Globe size={20} />
      <span className="uppercase text-sm font-bold">{language}</span>
    </button>
  );
}
