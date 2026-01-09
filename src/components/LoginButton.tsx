"use client";

import { login } from "@/lib/actions/auth";

export default function LoginButton() {
  return (
    <div>
      <button
        className="bg-blue-500, text-white, p-2, rounded"
        onClick={() => login()}
      >
        Sign In with Google
      </button>
    </div>
  );
}
