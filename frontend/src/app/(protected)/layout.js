import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function ProtectedLayout({ children }) {
  const cookieStore = await cookies();
  const session = cookieStore.get("session_id")?.value;
  
  if (!session) {
    redirect("/login");
  }
  
  const res = await fetch(`http://localhost:4001/hassession`, {
    method: "GET",
    headers: { Cookie: `session_id=${session}` },
    cache: "no-store",
  });
  
  console.log("fetched results ",res);
  if (!res.ok) {
    redirect("/login");
  }

  return <>{children}</>;
}