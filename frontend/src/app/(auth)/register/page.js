"use client";
import styles from "./page.module.css";
import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { fetchRegister, ValidateInput } from "./help";

export default function RegisterPage() {
  const router = useRouter();
  const [error, setError] = useState("");

  async function onSubmit(e) {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);

    const [errors, ok] = ValidateInput(formData);
    if (!ok) {
      const firstError = Object.values(errors)[0];
      setError(firstError);
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

        {error && <span className={styles.error}>{error}</span>}

        <div className={styles.row}>
          <div>
            <label>First Name</label>
            <input type="text" name="firstname" />
          </div>
          <div>
            <label>Last Name</label>
            <input type="text" name="lastname" />
          </div>
        </div>

        <label>Email</label>
        <input type="email" name="email" />

        <label>Password</label>
        <input type="password" name="password" />

        <div className={styles.row}>
          <div>
            <label>Date of Birth</label>
            <input type="date" name="dob" />
          </div>
          <div>
            <label>Gender</label>
            <label>
              <input type="radio" name="gender" value="male" /> Male
            </label>
            <label>
              <input type="radio" name="gender" value="female" /> Female
            </label>
          </div>
        </div>

        <label>
          Avatar <span>(optional)</span>
        </label>
        <input type="file" name="avatar" accept="image/*" />

        <label>
          Nickname <span>(optional)</span>
        </label>
        <input type="text" name="nickname" />

        <label>
          About Me <span>(optional)</span>
        </label>
        <textarea name="about" />

        <button type="submit">Register</button>

        <div>
          Already have an account? <Link href="/login">Login</Link>
        </div>
      </form>
    </div>
  );
}
