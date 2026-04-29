"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { handleLogin, isValidInput } from "./help";
import styles from "./page.module.css";
import Link from "next/link";

export default function Login() {
    const router = useRouter();
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState("");

    useEffect(() => {
        document.cookie = "session_id=; Max-Age=0; path=/";
    }, []);


    // submit event
    async function onSubmit(e) {
        e.preventDefault();
        setError("");
        setSubmitting(true);

        try {
            const formData = new FormData(e.currentTarget);
            const credentials = isValidInput(formData);

            if (!credentials) {
                setError("Invalid credential");
                return;
            }

            const ok = await handleLogin(credentials.email, credentials.password);
            if (!ok) {
                setError("Invalid credential");
                return;
            }
            router.replace("/");
        } catch {
            setError("Something went wrong");
        } finally {
            setSubmitting(false);
        }
    }

    return (
        <div className={styles.page}>
            <form className={styles.formLogin} onSubmit={onSubmit}>
                <h2>Sign In</h2>

                {error && <p className={styles.error}>{error}</p>}

                <label htmlFor="email">Email</label>
                <input type="email" id="email" name="email" disabled={submitting} />

                <label htmlFor="password">Password</label>
                <input type="password" id="password" name="password" disabled={submitting} />

                <button type="submit" disabled={submitting}>
                    {submitting ? "Signing in..." : "Sign In"}
                </button>

                <div>Don't have an account? <Link href="/register">Register</Link></div>
            </form>
        </div>
    );
}