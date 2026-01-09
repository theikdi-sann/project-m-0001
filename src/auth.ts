import { PrismaAdapter } from "@auth/prisma-adapter";
import NextAuth from "next-auth";
import Google from "next-auth/providers/google";
import { prisma } from "../lib/prisma";

export const { auth, handlers, signIn, signOut } = NextAuth({
  adapter: PrismaAdapter(prisma),
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
