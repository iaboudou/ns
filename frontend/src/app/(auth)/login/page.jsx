"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { handleLogin, isValidInput } from "./help";
import { Eye, EyeOff } from "lucide-react";

import styles from "./page.module.css";
import Link from "next/link";

export default function Login() {
  const router = useRouter();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [showPwd, setShowPwd] = useState(false);

  useEffect(() => {
    document.cookie = "session_id=; Max-Age=0; path=/";
  }, []);

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
      <div className={styles.bgOrb1} aria-hidden="true" />
      <div className={styles.bgOrb2} aria-hidden="true" />

      <form className={styles.formLogin} onSubmit={onSubmit}>
        <div className={styles.logoBlock}>
          <h2 className={styles.heading}>Welcome !</h2>
          <p className={styles.subheading}>Login to you account</p>
        </div>

        <div className={styles.fieldGroup}>
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            placeholder="Enter your email"
            disabled={submitting}
          />
        </div>

        <div className={styles.fieldGroup}>
          <label htmlFor="password">Password</label>
          <div className={styles.fieldWrap}>
            <input
              type={showPwd ? "text" : "password"}
              id="password"
              name="password"
              placeholder="Enter you password"
              disabled={submitting}
            />

            <button
              type="button"
              className={styles.eyeBtn}
              onClick={() => setShowPwd((v) => !v)}
              aria-label={showPwd ? "Masquer" : "Afficher"}
            >
              {showPwd ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
        </div>

        {error && <p className={styles.error}>{error}</p>}

        <button type="submit" disabled={submitting}>
          {submitting ? <span className={styles.spinner} /> : "Se connecter"}
        </button>

        <div className={styles.footer}>
          No account yet ? <Link href="/register">Register</Link>
        </div>
      </form>
    </div>
  );
}
