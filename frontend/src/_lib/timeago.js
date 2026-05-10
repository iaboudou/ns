export function timeAgo(date) {
  if (!date) return "now";

  let past;
  if (typeof date === "number") {
    past = date < 10000000000 ? date * 1000 : date;
  } else {
    past = new Date(date).getTime();
  }

  const now = Date.now();
  let diff = (now - past) / 1000;

  if (diff <= 0) return "now";
  diff = Math.floor(diff);

  if (diff < 60) return `${diff}s ago`;
  if (diff < 3600) return `${Math.floor(diff / 60)}min ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  if (diff < 2592000) return `${Math.floor(diff / 86400)}d ago`;

  let d = new Date(past).toLocaleDateString();
  return d === "Invalid Date" ? "now" : d;
}
