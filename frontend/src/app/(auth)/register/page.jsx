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
      setError(error || "You missed something, adventurer");
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
        <h2>Swear the Oath</h2>

        <div className={styles.row}>
          <div>
            <label>Bloodline Name</label>
            <input
              type="text"
              name="firstname"
              placeholder="Your first name"
              required
              minLength={2}
              maxLength={12}
              pattern="^[a-zA-Z]{2,12}$"
              title="First name must be 2–12 letters"
            />
          </div>
          <div>
            <label>Given Name</label>
            <input
              type="text"
              name="lastname"
              placeholder="Your last name"
              required
              minLength={2}
              maxLength={12}
              pattern="^[a-zA-Z]{2,12}$"
              title="Last name must be 2–12 letters"
            />
          </div>
        </div>

        <div>
          <label>Sigil</label>
          <input
            type="email"
            name="email"
            placeholder="seimor@gmail.com"
            required
            pattern="^[^\s@]+@[^\s@]+\.[^\s@]+$"
            title="Valid email required"
          />
        </div>

        <div>
          <label>Seal</label>
          <div className={styles.fieldWrap}>
            <input
              type={showPwd ? "text" : "password"}
              name="password"
              placeholder="Your password sir ..."
              required
              minLength={3}
              maxLength={32}
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
            <label>Birth Omen</label>
            <input type="date" name="dob" required />
          </div>
          <div>
            <label>Vessel</label>
            <div className={styles.radioGroup}>
              <label>
                <input type="radio" name="gender" value="male" required /> Male
              </label>
              <label>
                <input type="radio" name="gender" value="female" required />{" "}
                Female
              </label>
            </div>
          </div>
        </div>

        <div>
          <label>
            Mark of Identity <span>(optional)</span>
          </label>
          <input type="file" name="avatar" accept="image/*" />
        </div>

        <div>
          <label>
            Alias <span>(optional)</span>
          </label>
          <input
            type="text"
            name="nickname"
            placeholder="Name whispered in the dark"
            minLength={2}
            maxLength={30}
            pattern="^[a-zA-Z\s]{2,30}$"
            title="Nickname must be 2–30 letters/spaces"
          />
        </div>

        <div>
          <label>
            Chronicle <span>(optional)</span>
          </label>
          <textarea
            name="about"
            placeholder="Tell your story... or remain silent."
            minLength={1}
            maxLength={70}
          />
        </div>

        {error && <span className={styles.error}>{error}</span>}

        <button type="submit">Join the guild</button>

        <div className={styles.footer}>
          Already bound ? <Link href="/login">Identify yourself</Link>
        </div>
      </form>

      <img src="/pics/scroll.png" alt="scroll" />
    </div>
  );
}
