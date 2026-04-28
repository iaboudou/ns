import { NextResponse } from "next/server";

export async function middleware(request) {
  const session = request.cookies.get("session_id")?.value;
  // 
  const { pathname } = request.nextUrl;

  const isAuthPage = pathname === "/login" || pathname === "/register";
  const isProtected = pathname === "/" || pathname.startsWith("/post");

  // if no session exists
  if (!session && isProtected) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // check session validation
  if (session && (isAuthPage || isProtected)) {
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_BASE}/hassession`, {
        method: "GET",
        headers: { Cookie: `session_id=${session}` },
      });

      if (res.ok && isAuthPage) {
        // already connected
        return NextResponse.redirect(new URL("/", request.url));
      }

      if (!res.ok && isProtected) {
        // session not valid
        const response = NextResponse.redirect(new URL("/login", request.url));
        response.cookies.delete("session_id");
        return response;
      }

      if (!res.ok && isAuthPage) {
        // session not valid
        const response = NextResponse.next();
        response.cookies.delete("session_id");
        return response;
      }
    } catch {
      return NextResponse.next();
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/", "/login", "/register"],
};
