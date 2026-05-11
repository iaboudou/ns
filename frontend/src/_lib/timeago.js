export function timeAgo(datei) {
  if (!datei || datei === "now") datei = Date.now();

  let date = new Date(datei);
  if (date === "Invalid Date") date = new Date.now();
  
  const formatted = String(date.getDate()).padStart(2, "0") + "-" + String(date.getMonth() + 1).padStart(2, "0") + "-" + date.getFullYear() + " " + String(date.getHours()).padStart(2, "0") + ":" + String(date.getMinutes()).padStart(2, "0");
  return formatted;
}
