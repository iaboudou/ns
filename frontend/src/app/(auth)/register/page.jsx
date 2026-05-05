"use client";
import styles from "./page.module.css";
import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Eye, EyeOff } from "lucide-react";
import { fetchRegister, ValidateInput } from "./help";

export default function RegisterPage() {
  const router = useRouter();
  const [error, setError] = useState("");
  const [showPwd, setShowPwd] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);

    const [error, ok] = ValidateInput(formData);
    if (!ok) {
      setError(error || "Please fill in all required fields");
      return;
    }

    const [fetched, er] = await fetchRegister(formData);
    if (!fetched) {
      setError(er);
      return;
    }

    router.replace("/login");
  }

  return (
    <div className={styles.page}>
      <form className={styles.formRegister} onSubmit={onSubmit}>
        <h2>Create Account</h2>

        <div className={styles.row}>
          <div>
            <label>First Name</label>
            <input type="text" name="firstname" placeholder="John" />
          </div>
          <div>
            <label>Last Name</label>
            <input type="text" name="lastname" placeholder="Doe" />
          </div>
        </div>

        <div>
          <label>Email</label>
          <input type="email" name="email" placeholder="you@example.com" />
        </div>

        <div>
          <label>Password</label>
          <div className={styles.fieldWrap}>
            <input
              type={showPwd ? "text" : "password"}
              name="password"
              placeholder="••••••••"
            />

            <button
              type="button"
              className={styles.eyeBtn}
              onClick={() => setShowPwd((v) => !v)}
              aria-label={showPwd ? "Hide password" : "Show password"}
            >
              {showPwd ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
        </div>

        <div className={styles.row}>
          <div>
            <label>Date of Birth</label>
            <input type="date" name="dob" />
          </div>
          <div>
            <label>Gender</label>
            <div className={styles.radioGroup}>
              <label>
                <input type="radio" name="gender" value="male" /> Male
              </label>
              <label>
                <input type="radio" name="gender" value="female" /> Female
              </label>
            </div>
          </div>
        </div>

        <div>
          <label>
            Avatar <span>(optional)</span>
          </label>
          <input type="file" name="avatar" accept="image/*" />
        </div>

        <div>
          <label>
            Nickname <span>(optional)</span>
          </label>
          <input type="text" name="nickname" placeholder="coolname42" />
        </div>

        <div>
          <label>
            About Me <span>(optional)</span>
          </label>
          <textarea
            name="about"
            placeholder="Tell us a little about yourself…"
          />
        </div>

        {error && <span className={styles.error}>{error}</span>}

        <button type="submit">Register</button>

        <div className={styles.footer}>
          Already have an account? <Link href="/login">Login</Link>
        </div>
      </form>
    </div>
  );
}
