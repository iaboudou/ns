export async function fetchRegister(formData) {
  const response = await fetch(`/api/register`, {
    method: "POST",
    credentials: "include",
    body: formData,
  });

  const er = await response.json().catch(() => ({}));
  if (!response.ok) return [false, er.message || "Registration failed"];
  return [true, null];
}

export function ValidateInput(formData) {
  const { email, password, firstname, lastname, dob, gender, nickname, about } =
    Object.fromEntries(formData.entries());

  if (!email || !password || !firstname || !lastname || !dob || !gender)
    return ["all required fields", false];

  // 
  const birthDate = new Date(dob);
  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const m = today.getMonth() - birthDate.getMonth();
  if (m < 0 || (m === 0 && today.getDate() < birthDate.getDate())) {
    age--;
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
    return ["Invalid email address", false];
  if (password.length < 3 || password.length > 32)
    return ["Password must be 3–32 chars", false];
  if (!/^[a-zA-Z]{2,12}$/.test(firstname))
    return ["First name must be 2–12 letters", false];
  if (!/^[a-zA-Z]{2,12}$/.test(lastname))
    return ["Last name must be 2–12 letters", false];
  
  if (isNaN(birthDate.getTime())) return ["Invalid date", false];
  if (age <= 15 || age >= 200) return ["You must be between 16 and 199 years old", false];

  if (!["male", "female"].includes(gender.toLowerCase()))
    return ["Gender must be male or female", false];
  if (nickname && !/^[a-zA-Z\s]{2,30}$/.test(nickname))
    return ["Nickname must be 2–30 letters/spaces", false];
  if (about && (about.length < 1 || about.length > 70))
    return ["About me must be 1–70 chars", false];

  return [null, true];
}
