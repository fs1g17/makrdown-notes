import axios, { AxiosError } from "axios";
import { NextRequest, NextResponse } from "next/server";

const PUBLIC_PATHS = ["/sign-in", "/sign-up"];

async function validateAuthToken(auth_token: string | undefined) {
  if (!auth_token) {
    return false;
  }

  try {
    await axios(`${process.env.BACKEND_BASE_URL}/me`, {
      method: "GET",
      headers: {
        Cookie: `auth_token=${auth_token}`,
      },
    });
    return true;
  } catch (error) {
    const err = error as AxiosError;
    if (err.response?.status === 401) {
      console.log("JWT validation failed", err);
    }
    // TODO: handle other errors here, e.g. 500
    return false;
  }
}

export default async function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;

  if (
    pathname.includes("/_next/")
  ) {
    return NextResponse.next();
  }

  const auth_token = req.cookies.get("auth_token");
  const isPublicPath = PUBLIC_PATHS.some((path) => pathname.startsWith(path));

  if (!auth_token && !isPublicPath) {
    return NextResponse.redirect(new URL("/sign-in", req.url));
  }

  if (!auth_token && isPublicPath) {
    return NextResponse.next();
  }

  try {
    const isValidAuthToken = await validateAuthToken(auth_token?.value);

    if (!isValidAuthToken) {
      return isPublicPath
        ? NextResponse.next()
        : NextResponse.redirect(new URL("/sign-in", req.url));
    }

    if (isPublicPath) {
      return NextResponse.redirect(new URL("/folders", req.url));
    }

    return NextResponse.next();
  } catch (error) {
    console.error("Middleware error:", error);
    return NextResponse.next();
  }
}

export const config = {
  matcher: [
    /*
     * Match all request paths EXCEPT for the ones starting with:
     * - web/api (API routes)
     * - web/_next/static (static files)
     * - web/_next/image (image optimization files)
     * Also excludes all paths ending with .png
     */
    "/((?!api|_next/image|favicon.ico|.*\\.png$).*)",
  ],
};