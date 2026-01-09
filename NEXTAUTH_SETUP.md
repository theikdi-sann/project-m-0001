# Deep Dive: NextAuth.js (Auth.js v5) Implementation Guide

## Introduction
This guide provides a comprehensive walkthrough for implementing authentication in Next.js 15+ (App Router) using **Auth.js v5**. It is designed to explain not just the *code*, but the *mechanics* behind it—the What, How, When, and Why.

---

## 1. Installation

**What:** Installing the core library. Note that v5 is currently in beta but is the standard for the App Router.
**Why:** Next.js 13+ introduced the App Router, which changed how request handling works (from Request/Response objects to Web Standard APIs). v5 is rewritten to align with these standards.

Run the following command in your `website` directory:
```bash
npm install next-auth@beta
```

---

## 2. Environment Variables

**What:** Storing secrets securely.
**Why:** Authentication relies on cryptographic signatures. The `AUTH_SECRET` is the private key used to sign and verify Session Tokens (JWTs). If this leaks, attackers can forge sessions.

Create or update `.env.local` in your `website` folder:

```env
# 1. The Master Key
# Used to encrypt cookies and tokens.
# Generate securely via command line: npx auth secret
AUTH_SECRET="your-generated-secret-key-here"

# 2. OAuth Provider Secrets (Example: GitHub)
# The ID is public-facing; the SECRET proves to GitHub that your app is making the request.
AUTH_GITHUB_ID="client_id_from_github"
AUTH_GITHUB_SECRET="client_secret_from_github"
```

---

## 3. Configuration (`auth.ts`)

**What:** The central brain of your authentication logic.
**How:** This file exports methods that other parts of your app will import.
**Why:** In v5, configuration is decoupled from the route handler. This allows you to import `auth` helper functions anywhere (Middleware, Server Actions, Components) without importing the entire HTTP request handling logic.

Create `website/src/auth.ts`:

```typescript
import NextAuth from "next-auth"
import GitHub from "next-auth/providers/github"

// What is happening here?
// 1. We initialize NextAuth with our config.
// 2. It returns a set of helper functions specialized for our app.
export const { 
  handlers, // The API handler (GET/POST) for the /api/auth route
  signIn,   // Helper to trigger sign-in flow (server-side)
  signOut,  // Helper to trigger sign-out flow (server-side)
  auth      // The universal method to fetch the current session
} = NextAuth({
  providers: [
    GitHub,
    // Add Google, Credentials, etc. here
  ],
  // OPTIONAL: Callbacks allow you to hook into the lifecycle
  callbacks: {
    async session({ session, token, user }) {
      // WHEN: This runs every time the session is checked.
      // WHY: To expose extra data (like User ID) to the client that isn't there by default.
      if (session.user && token.sub) {
        session.user.id = token.sub; 
      }
      return session;
    }
  }
})
```

---

## 4. The Route Handler

**What:** The bridge between the outside world (browser) and your internal `auth.ts` logic.
**When:** The browser needs to perform auth actions (login, logout, callback handling).
**How:** It catches any request to `/api/auth/*` and passes it to the `handlers` we exported in step 3.
**Why:** OAuth protocols require specific HTTP endpoints to receive "Callback" codes from providers like Google. This file sets up those endpoints automatically.

Create `website/src/app/api/auth/[...nextauth]/route.ts`:

```typescript
import { handlers } from "@/auth" // Importing from our central config

// We export GET and POST methods to handle HTTP requests.
// Next.js App Router automatically maps these to API endpoints.
export const { GET, POST } = handlers
```

---

## 5. Middleware (Edge Security)

**What:** A gatekeeper that runs *before* a request completes.
**When:** On every single page load or API call (defined by the matcher).
**Why:**
1.  **Performance:** Middleware runs on the Edge (closer to the user). It checks auth status *before* your expensive Server Components or Database queries even start.
2.  **Security:** It prevents unauthenticated users from even seeing the loading state of a protected page.

Create `website/src/middleware.ts`:

```typescript
// We export the 'auth' function as 'middleware'.
// This is a shorthand provided by NextAuth. It automatically:
// 1. Checks for a session.
// 2. Updates the session expiry (rolling sessions).
// 3. Allows you to define custom logic (like redirects).
export { auth as middleware } from "@/auth"

// The Matcher: Tells Next.js WHICH routes to run this middleware on.
// We use a negative lookahead regex here:
// "Run on everything EXCEPT: api routes, static files, images, and favicon."
export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
}
```

---

## 6. Accessing the Session

### A. In Server Components (Preferred)

**What:** Fetching user data directly on the server.
**Why:** Fast, secure, and SEO-friendly. No loading spinners needed because the page renders with the user data already present.

```tsx
import { auth } from "@/auth"

export default async function DashboardPage() {
  const session = await auth()
  
  if (!session) {
    // Handle unauthenticated state (e.g., redirect or show message)
    return <div>Access Denied</div>
  }
 
  return (
    <div>
      <h1>Welcome back, {session.user?.name}</h1>
      <p>Email: {session.user?.email}</p>
    </div>
  )
}
```

### B. In Client Components

**What:** Using React Context to access session data in the browser.
**When:** You need user data for interactivity (e.g., a "User Menu" dropdown) *after* the page has loaded.
**Why:** Requires wrapping your app in a `SessionProvider`.

**Step 1: The Provider (`src/app/layout.tsx`)**
```tsx
import { SessionProvider } from "next-auth/react"

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        {/* Why wrap? This makes session data available to all 'useSession' hooks deep in the tree */}
        <SessionProvider>
          {children}
        </SessionProvider>
      </body>
    </html>
  )
}
```

**Step 2: The Component**
```tsx
"use client" // Must be a client component
import { useSession, signIn, signOut } from "next-auth/react"

export default function UserButton() {
  // 'data' is the session object. 'status' is "loading" | "authenticated" | "unauthenticated"
  const { data: session, status } = useSession()

  if (status === "loading") return <p>Loading...</p>

  if (session) {
    return <button onClick={() => signOut()}>Sign out</button>
  }
  return <button onClick={() => signIn()}>Sign in</button>
}
```

---

## 7. Database Integration (Prisma Adapter)

**What:** Connecting NextAuth to your PostgreSQL database.
**Why:**
1.  **Persistence:** Without a DB, users are "logged in" via a JWT cookie. If they clear cookies, they are gone. A DB allows you to keep user accounts permanently.
2.  **User Profiles:** You can attach extra data (Roles, Subscription Status) to the User table.
3.  **Magic Links:** Email sign-in *requires* a database to store the temporary verification tokens.

**How to implement:**

1.  **Install Adapter:**
    ```bash
    npm install @auth/prisma-adapter prisma @prisma/client
    ```

2.  **Update `auth.ts`:**
    ```typescript
    import NextAuth from "next-auth"
    import { PrismaAdapter } from "@auth/prisma-adapter"
    import { PrismaClient } from "@prisma/client"
    import GitHub from "next-auth/providers/github"

    const prisma = new PrismaClient()

    export const { handlers, auth, signIn, signOut } = NextAuth({
      adapter: PrismaAdapter(prisma), // <--- The Bridge
      providers: [GitHub],
      // Strategy is crucial here:
      // "jwt": (Default) Session data is stored in the cookie. Database is used only for User/Account storage.
      // "database": Session token is stored in DB. Cookie only holds an ID.
      session: { strategy: "jwt" }, 
    })
    ```

3.  **Update `schema.prisma`:**
    You must add the specific models NextAuth expects.
    (Refer to the Auth.js docs for the exact schema: https://authjs.dev/getting-started/adapters/prisma)

---

## 8. Server Actions (The Modern Way)

**What:** Calling server-side logic directly from a form.
**Why:** Cleaner than creating an API route just to handle a button click. Works without JavaScript enabled.

```tsx
import { signIn, signOut } from "@/auth"
 
export function LoginButton() {
  return (
    <form
      action={async () => {
        "use server"
        await signIn("github")
      }}
    >
      <button type="submit">Sign in with GitHub</button>
    </form>
  )
}
```