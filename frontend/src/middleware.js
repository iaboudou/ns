import { NextResponse } from "next/server";

// this is the middleware of the frontend where we redirect the user based on the session
export async function middleware(request) {
  // // get the session from the cookies
  const session = request.cookies.get("session_id")?.value;
  
  // get if the user is in auth page or no based on the path name
  const { pathname } = request.nextUrl;
  const isAuthPage = pathname === "/login" || pathname === "/register";

  // if no session exists
  if (!session && !isAuthPage) {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  // check session validation
  if (session) {
    try {
      const res = await fetch(`http://localhost:4001/hassession`, {
        method: "GET",
        headers: { Cookie: `session_id=${session}` },
      });

      // already connected
      if (res.ok && isAuthPage) {
        return NextResponse.redirect(new URL("/", request.url));
      }

      // session not valid
      if (!res.ok) {
        if (!isAuthPage) {
          const response = NextResponse.redirect(new URL("/login", request.url));
          response.cookies.delete("session_id");
          return response;
        } else {
          const response = NextResponse.next();
          response.cookies.delete("session_id");
          return response;
        }
      }
    } catch {
      return NextResponse.next();
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!_next|favicon.ico|.*\\..*).*)",
  ],
};