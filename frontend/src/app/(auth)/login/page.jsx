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
        setError("The ritual is incomplete");
        return;
      }

      const ok = await handleLogin(credentials.email, credentials.password);
      if (!ok) {
        setError("The seal rejects you");
        return;
      }

      router.replace("/");
    } catch {
      setError("The ritual has failed");
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
          <h2 className={styles.heading}>Return, Wanderer</h2>
          <p className={styles.subheading}>
            Enter your credentials to reclaim your path
          </p>
        </div>

        <div className={styles.fieldGroup}>
          <label htmlFor="email">Ethereal Sigil</label>
          <input
            type="email"
            id="email"
            name="email"
            placeholder="seimor@gmail.com"
            disabled={submitting}
            required
            pattern="^[^\s@]+@[^\s@]+\.[^\s@]+$"
            title="Valid email required"
          />
        </div>

        <div className={styles.fieldGroup}>
          <label htmlFor="password">Sacred Seal</label>
          <div className={styles.fieldWrap}>
            <input
              type={showPwd ? "text" : "password"}
              id="password"
              name="password"
              placeholder="Your password sir ..."
              disabled={submitting}
              required
              minLength={3}
              maxLength={32}
            />

            <button
              type="button"
              className={styles.eyeBtn}
              onClick={() => setShowPwd((v) => !v)}
              aria-label={showPwd ? "Hide" : "Show"}
            >
              {showPwd ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
        </div>

        {error && <p className={styles.error}>{error}</p>}

        <button type="submit" disabled={submitting}>
          {submitting ? <span className={styles.spinner} /> : "Enter the Gate"}
        </button>

        <div className={styles.footer}>
          No pact yet ? <Link href="/register">Forge one</Link>
        </div>
      </form>

      <img src="/static/sword.png" alt="sword" />
    </div>
  );
}
