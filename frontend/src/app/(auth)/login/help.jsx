/**
 * Handle login request function.
 * 
 * @param {email} email 
 * @param {password} password 
 * @returns {Promise<boolean>}
 */
export async function handleLogin(email, password) {
  const response = await fetch(`/api/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ email, password }),
  });

  const data = await response.json().catch(() => ({}));
  if (data.user) localStorage.setItem("user", JSON.stringify(data.user));
  if (!response.ok) return false;
  return true;
}

/**
 * Validates user login input fields.
 *
 * @param {FormData} formData - Raw form data submitted by the user
 * @returns {{ email: string, password: string } | null}
 */
export function isValidInput(formData) {
  const email = String(formData.get("email") ?? "");
  const password = String(formData.get("password") ?? "");

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) return null;
  if (password.length < 3 || password.length > 32) return null;

  return { email, password };
}
