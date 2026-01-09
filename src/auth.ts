import NextAuth from "next-auth";
import Google from "next-auth/providers/google";

export const { auth, handlers, signIn, signOut } = NextAuth({
  providers: [Google],
  callbacks: {
    async signIn({account, profile}){

        // Google Check
        if(account?.provider === "google") {
            return !!profile?.email_verified}

        // Other provider Check
        
        return false
    }
  },
});
