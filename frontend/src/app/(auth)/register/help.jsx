// help.js

export async function fetchRegister(formData) {
  const response = await fetch(`${"http://localhost:4001"}/api/register`, {
    method: "POST",
    body: formData,
  });

  const er = await response.text();
  if (!response.ok) return [false, er];
  return [true, null];
}


export function ValidateInput(formData) {
  const { email, password, firstname, lastname, dob, gender, nickname, about } =
    Object.fromEntries(formData.entries());

  if (!email || !password || !firstname || !lastname || !dob || !gender)
    return ["all required fields", false];


  const age = Math.floor(
    (Date.now() - new Date(dob).getTime()) / (1000 * 60 * 60 * 24 * 365)
  );

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
    return ["Invalid email address", false];
  if (password.length < 3 || password.length > 32)
    return ["Password must be 3–32 chars", false];
  if (!/^[a-zA-Z]{2,12}$/.test(firstname))
    return ["First name must be 2–12 letters", false];
  if (!/^[a-zA-Z]{2,12}$/.test(lastname))
    return ["Last name must be 2–12 letters", false];
  if (isNaN(new Date(dob).getTime())) return ["Invalid date", false];
  if (age < 15 || age > 200) return ["You can't use this website", false];
  if (!["male", "female"].includes(gender.toLowerCase()))
    return ["Gender must be male or female", false];
  if (nickname && !/^[a-zA-Z\s]{2,30}$/.test(nickname))
    return ["Nickname must be 2–30 letters/spaces", false];
  if (about && (about.length < 1 || about.length > 70))
    return ["About me must be 1–70 chars", false];

  return [null, true];
}
