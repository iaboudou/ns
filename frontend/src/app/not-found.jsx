import Link from "next/link";

export default function NotFound() {
    return (
        <div style={{
            minHeight: "100vh",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            background: "#f5f5f5",
            gap: "12px",
        }}>
            <h1 style={{ fontSize: "64px", fontWeight: "700", margin: 0, color: "#222" }}>404</h1>
            <p style={{ fontSize: "16px", color: "#555", margin: 0 }}>Page not found</p>
            <Link href="/" style={{ fontSize: "14px", color: "#333", fontWeight: "500" }}>
                Go home
            </Link>
        </div>
    );
}