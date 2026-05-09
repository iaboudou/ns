"use client";

import { useParams, useRouter } from "next/navigation";
import { useEffect, useMemo, useRef, useState } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";
import styles from "./chat.module.css";

export default function Page() {
  const params = useParams();
  const router = useRouter();
  const chatId = params.id;
  const [inputText, setInputText] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const [chatUser, setChatUser] = useState(null);
  const scrollRef = useRef(null);

  const { messages, sendMessage, port, hasMoreMap, onlineUsers } =
    useWebSocket();
  const hasMore = hasMoreMap[chatId] !== false; // default true
  const isChatUserOnline = chatUser ? onlineUsers.includes(chatUser.id) : false;

  // fetch user details for the header
  useEffect(() => {
    let ignore = false;

    const fetchUser = async () => {
      try {
        const resp = await fetch(
          `/api/getUsers?search=${encodeURIComponent(chatId)}`,
          {
            credentials: "include",
          },
        );
        const res = await resp.json();
        if (ignore) return;

        if (
          res.code === 200 &&
          Array.isArray(res.data) &&
          res.data.length > 0
        ) {
          // extra filter
          const user = res.data.find((u) => u.id === chatId);
          if (user) {
            setChatUser(user);
            return;
          }
        }

        router.replace("/chat");
      } catch {
        if (!ignore) router.replace("/chat");
      }
    };

    setChatUser(null);

    if (chatId) {
      fetchUser();
    } else {
      router.replace("/chat");
    }

    return () => {
      ignore = true;
    };
  }, [chatId, router]);

  useEffect(() => {
    if (port && chatUser?.id) {
      port.postMessage({
        type: "send",
        payload: {
          type: "mark_read",
          receiver_Id: chatUser.id,
        },
      });

      port.postMessage({
        type: "messages_read",
        receiver_Id: chatUser.id,
      });

      // load history
      port.postMessage({
        type: "send",
        payload: {
          type: "load_history",
          receiver_Id: chatUser.id,
          last_read_time: 0,
        },
      });
    }
  }, [port, chatUser?.id]);

  // filter messages for this specific conversation
  const conversationMessages = useMemo(() => {
    if (!chatUser?.id) return [];

    return messages
      .filter(
        (msg) =>
          msg.sender_id === chatUser.id || msg.receiver_id === chatUser.id,
      )
      .sort((a, b) => a.created_at - b.created_at);
  }, [messages, chatUser?.id]);

  const lastReadMsgId = useRef(null);

  // mark new messages as read if they arrive while we are actively in the chat
  useEffect(() => {
    if (port && chatUser?.id && conversationMessages.length > 0) {
      const lastMsg = conversationMessages[conversationMessages.length - 1];
      if (
        lastMsg.isNew &&
        lastMsg.sender_id === chatUser.id &&
        lastReadMsgId.current !== lastMsg.id
      ) {
        lastReadMsgId.current = lastMsg.id;
        port.postMessage({
          type: "send",
          payload: {
            type: "mark_read",
            receiver_Id: chatUser.id,
          },
        });

        port.postMessage({
          type: "messages_read",
          receiver_Id: chatUser.id,
        });
      }
    }
  }, [conversationMessages, port, chatUser?.id]);

  // scroll to bottom when messages change
  useEffect(() => {
    if (scrollRef.current) {
      const lastMsg = conversationMessages[conversationMessages.length - 1];
      if (lastMsg?.isNew || conversationMessages.length <= 10) {
        scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
      }
    }
  }, [conversationMessages]);

  const handleLoadMore = () => {
    if (!port || !chatUser?.id || conversationMessages.length === 0) return;
    const oldestTime = conversationMessages[0].created_at;
    port.postMessage({
      type: "send",
      payload: {
        type: "load_history",
        receiver_Id: chatUser.id,
        last_read_time: oldestTime,
      },
    });
  };

  const commonEmojis = [
    "👍",
    "😀",
    "😂",
    "🥰",
    "😎",
    "🤔",
    "😅",
    "🔥",
    "❤️",
    "🙏",
    "✨",
    "🎉",
  ];

  if (!chatUser) {
    return (
      <div className={styles.chatContainer}>
        <div className={styles.emptyState}>
          <h1>you did not select the chat yet</h1>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.chatContainer}>
      <div className={styles.chatHeader}>
        <div className={styles.headerIdentity}>
          <span className={styles.headerName}>
            {chatUser.firstname} {chatUser.lastname}
          </span>
          {chatUser.nickname && (
            <span className={styles.headerNick}> @{chatUser.nickname}</span>
          )}
        </div>
        <span
          className={`${styles.headerStatus} ${isChatUserOnline ? styles.headerStatusOnline : ""}`}
        >
          <span className={styles.headerStatusDot} />
          {isChatUserOnline ? "Online" : "Offline"}
        </span>
      </div>

      <div className={styles.messagesArea} ref={scrollRef}>
        {hasMore && conversationMessages.length >= 10 && (
          <button
            type="button"
            className={styles.loadMore}
            onClick={handleLoadMore}
          >
            Load older messages
          </button>
        )}
        {conversationMessages.map((msg) => (
          <div
            key={msg.id}
            className={`${styles.messageRow} ${msg.sender_id !== chatUser.id ? styles.sentRow : styles.receivedRow}`}
          >
            <div
              className={`${styles.messageBubble} ${msg.sender_id !== chatUser.id ? styles.sentBubble : styles.receivedBubble}`}
            >
              {msg.content}
            </div>
          </div>
        ))}
      </div>

      <div className={styles.inputArea}>
        <div className={styles.emojiWrapper}>
          <button
            type="button"
            className={styles.emojiToggle}
            onClick={() => setShowEmojis(!showEmojis)}
          >
            👍
          </button>
          {showEmojis && (
            <div className={styles.emojiPicker}>
              {commonEmojis.map((emoji) => (
                <button
                  type="button"
                  key={emoji}
                  className={styles.emojiItem}
                  onClick={() => {
                    setInputText((prev) => prev + emoji);
                    setShowEmojis(false);
                  }}
                >
                  {emoji}
                </button>
              ))}
            </div>
          )}
        </div>

        <input
          type="text"
          className={styles.inputField}
          placeholder="Type a message..."
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              const content = inputText.trim();
              if (content === "") return;
              sendMessage({
                type: "send",
                payload: {
                  type: "chat",
                  receiver_Id: chatUser.id,
                  content: content,
                },
              });
              setInputText("");
            }
          }}
        />
        <button
          type="button"
          className={styles.sendButton}
          onClick={(e) => {
            e.preventDefault();
            const content = inputText.trim();
            if (content === "") return;
            sendMessage({
              type: "send",
              payload: {
                type: "chat",
                receiver_Id: chatUser.id,
                content: content,
              },
            });
            setInputText("");
          }}
        >
          Send
        </button>
      </div>
    </div>
  );
}
