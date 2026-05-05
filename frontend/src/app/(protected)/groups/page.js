"use client";

import { notFound, usePathname } from "next/navigation";

export default function Groups() {
  const path = usePathname();
  if (path === "/groups") notFound();
  return <></>;
}
