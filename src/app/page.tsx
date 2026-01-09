import { auth } from "@/auth";
import LoginButton from "@/components/LoginButton";

export default async function Home() {
  const session = await auth();
  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
        <div>
          {session ? (
            <p>Welcome</p>
          ) : (
            <>
              {" "}
              <p>Login</p> <LoginButton />{" "}
            </>
          )}
        </div>
      </main>
    </div>
  );
}
