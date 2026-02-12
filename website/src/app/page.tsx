"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { supabase } from "@/lib/supabase";
import { Scan, AlertCircle, Loader2 } from "lucide-react";

function SessionLoader() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [status, setStatus] = useState<"idle" | "loading" | "error" | "success">("idle");
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    const sessionId = searchParams.get("session_id");

    if (sessionId) {
      validateAndJoinSession(sessionId);
    }
  }, [searchParams]);

  async function validateAndJoinSession(sessionId: string) {
    setStatus("loading");
    try {
      // 1. Validate Session with Go Backend
      const apiUrl = process.env.NEXT_PUBLIC_API_URL;
      const res = await fetch(`${apiUrl}/sessions/${sessionId}`);
      
      if (!res.ok) {
        throw new Error("Invalid or expired session.");
      }

      // 2. Anonymous Auth with Supabase (for RLS/Realtime security)
      const { error: authError } = await supabase.auth.signInAnonymously();
      if (authError) throw authError;

      // 3. Store Session & Redirect
      localStorage.setItem("qr_session_id", sessionId);
      setStatus("success");
      router.push("/menu");

    } catch (err: any) {
      console.error(err);
      setStatus("error");
      setErrorMessage(err.message || "Failed to join session");
    }
  }

  if (status === "loading") {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-4 text-center">
        <Loader2 className="w-12 h-12 text-blue-600 animate-spin mb-4" />
        <h2 className="text-xl font-semibold">Joining Table...</h2>
      </div>
    );
  }

  if (status === "error") {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen p-4 text-center text-red-600">
        <AlertCircle className="w-16 h-16 mb-4" />
        <h2 className="text-2xl font-bold mb-2">Oops!</h2>
        <p>{errorMessage}</p>
        <button 
          onClick={() => window.location.reload()}
          className="mt-6 px-6 py-2 bg-blue-600 text-white rounded-full"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4 text-center bg-gray-50 text-gray-900">
      <div className="bg-white p-8 rounded-2xl shadow-xl max-w-sm w-full">
        <Scan className="w-16 h-16 text-blue-600 mx-auto mb-6" />
        <h1 className="text-3xl font-bold mb-2">Welcome</h1>
        <p className="text-gray-500 mb-8">Scan the QR code on your table to start ordering.</p>
        
        <div className="p-4 bg-blue-50 rounded-lg text-sm text-blue-700">
          Waiting for QR Code...
        </div>
      </div>
    </div>
  );
}

export default function Home() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <SessionLoader />
    </Suspense>
  );
}
