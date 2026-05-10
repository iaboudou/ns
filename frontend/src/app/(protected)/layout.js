import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import Leftbar from "@/components/leftbar/leftbar";
import styles from "./page.module.css";
import { WebSocketProvider } from "@/lib/UseWebsocket";

export default async function ProtectedLayout({ children }) {
  const cookieStore = await cookies();
  const session = cookieStore.get("session_id")?.value;

  if (!session) {
    redirect("/login");
  }

  const res = await fetch(`http://backend:4001/hassession`, {
    method: "GET",
    headers: { Cookie: `session_id=${session}` },
    cache: "no-store",
  });

  if (!res.ok) {
    redirect("/login");
  }

  return (
    <WebSocketProvider>
      <div className={styles?.wrapper}>
        <nav className={styles.leftbar}>
          <Leftbar />
        </nav>

        <main id="mainpage" className={styles.main}>
          <div className={styles.posts}>{children}</div>
        </main>
      </div>
    </WebSocketProvider>
  );
}
