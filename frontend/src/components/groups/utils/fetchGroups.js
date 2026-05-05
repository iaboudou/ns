import { GetGroups } from "@/_lib/group";

export default async function handleFetchGroups({
  groups,
  setGroups,
  setLoading,
  setHasMore,
  tab,
  hasMore = true,
  search,
  isReset = false,
}) {
  if (!hasMore && !isReset) return;

  const lastGroup = isReset ? undefined : groups?.at(-1);
  const lastTime = lastGroup?.created_at;
  const lastId = lastGroup?.id;

  try {
    setLoading(true);
    const newGroups = await GetGroups(tab, search, lastId, lastTime);

    if (newGroups.length < 10) setHasMore(false);

    setGroups(
      isReset
        ? newGroups
        : (prev) => {
            const map = new Map(prev.map((g) => [g.id, g]));
            newGroups.forEach((g) => map.set(g.id, g));
            return Array.from(map.values());
          },
    );
  } catch (err) {
    alert(err);
  } finally {
    setLoading(false);
  }
}
