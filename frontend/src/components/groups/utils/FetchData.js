import { GetData } from "@/_lib/group";

export default async function handleFetchData(
  data,
  setData,
  setLoading,
  setHasmore,
  tab,
  groupId,
  router = null,
) {
  const lastEvent = data.length === 0 ? undefined : data.at(-1);
  const lastTime = lastEvent?.created_at;
  const lastId = lastEvent?.id;

  try {
    setLoading(true);
    const newData = await GetData(tab, groupId, lastTime, lastId);

    if (newData.length < 10) setHasmore(false);

    setData((prev) => {
      const map = new Map(prev.map((g) => [g.id, g]));
      newData.forEach((g) => map.set(g.id, g));
      return Array.from(map.values());
    });
  } catch (err) {
    if (err.message !== "not creator") alert(err.message);
    else router.replace(`/groups/${groupId}/posts`);
  } finally {
    setLoading(false);
  }
}
