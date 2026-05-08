"use client"

import ChatWindow from "@/components/chat/ChatWindow";
import { useEffect } from "react";

export default function Page() {
  useEffect(() => {
    console.log("🟢 MOUNT CHAT PAGE");
    return () => {
      console.log("🔴 UNMOUNT CHAT PAGE");
    };
  }, []);

  return <ChatWindow isSelected={false} />;
}
