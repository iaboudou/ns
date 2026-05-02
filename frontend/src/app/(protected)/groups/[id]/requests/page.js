import { redirect } from "next/navigation";
import { cookies } from "next/headers";
import { GetGroup } from "@/_lib/group";
import GroupRequests from "@/components/groups/requests/mainPage";

export default async function RequestsPage({ params }) {
  const { id } = await params;
  const cookieStore = await cookies();
  const cookie = cookieStore.toString();

  try {
    const group = await GetGroup(id, cookie);
    if (!group.isCreator) redirect(`/groups/${id}/posts`);
  } catch {
    redirect(`/groups/${id}/posts`);
  }

  return <GroupRequests id={id} />;
}
