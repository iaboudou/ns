export function timeAgo(dateStr) {
  const now = Date.now() + (60 * 60* 1000)
  const past = new Date(dateStr).getTime()
  
  let diff = (now - past) / 1000

  if (diff <= 0) return "now"

  diff = Math.floor(diff)

  if (diff < 60) return `${diff}s ago`
  if (diff < 3600) return `${Math.floor(diff / 60)}min ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  if (diff < 2592000) return `${Math.floor(diff / 86400)}d ago`

  let d =  new Date(dateStr).toLocaleDateString() 
  return d == "Invalid Date" ? "now" : d
}
