"use client";

import { useParams } from "next/navigation";
import styles from "@/app/(protected)/chat/[id]/chat.module.css";
import groupChatStyles from "./group-chat-overrides.module.css";
import { useState, useEffect, useRef } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function GroupChat() {
  const params = useParams();
  const [inputText, setInputText] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const [myID, setMyID] = useState(null);
  const scrollRef = useRef(null);

  const { messages, sendMessage, port, hasMore } = useWebSocket();

  // fetch my personnal info=
  useEffect(() => {
    const fetchMe = async () => {
      try {
        const resp = await fetch(`/api/getpersonalinfo`, {
          credentials: "include",
        });
        const res = await resp.json().catch(() => ({}));
        //
        if (res.user) {
          setMyID(res.user.id);
        }
      } catch {}
    };
    fetchMe();
  }, []);


  // fetch messages in case the user reconnect
  useEffect(() => {
    if (port && params.id) {
      port.postMessage({
        type: "send",
        payload: {
          type: "load_group_history",
          receiver_Id: params.id,
          last_read_time: 0,
        },
      });
    }
  }, [port, params.id]);

  // filter messages for this specific group
  const conversationMessages = messages
    .filter(
      (msg) =>
        msg.group_id === params.id ||
        (msg.is_group && msg.receiver_Id === params.id),
    )
    .sort((a, b) => a.created_at - b.created_at);

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
    if (conversationMessages.length === 0) return;
    const oldestTime = conversationMessages[0].created_at;
    port.postMessage({
      type: "send",
      payload: {
        type: "load_group_history",
        receiver_Id: params.id,
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

  return (
    <div
      className={`${styles.chatContainer} ${groupChatStyles.chatContainer}`}
    >
      <div
        className={`${styles.messagesArea} ${groupChatStyles.messagesArea}`}
        ref={scrollRef}
      >
        {hasMore && (
          <button className={styles.loadMore} onClick={handleLoadMore}>
            Load older messages
          </button>
        )}
        {conversationMessages.map((msg) => {
          const isMe = msg.sender_id === myID;
          return (
            <div
              key={msg.id}
              className={`${styles.messageRow} ${isMe ? styles.sentRow : styles.receivedRow}`}
            >
              <div
                className={`${styles.messageBubble} ${isMe ? styles.sentBubble : styles.receivedBubble}`}
              >
                {!isMe && msg.sender_fullname && (
                  <div className={styles.senderName}>{msg.sender_fullname}</div>
                )}
                {msg.content}
              </div>
            </div>
          );
        })}
      </div>

      <div className={styles.inputArea}>
        <div className={styles.emojiWrapper}>
          <button
            className={styles.emojiToggle}
            onClick={() => setShowEmojis(!showEmojis)}
          >
            😀
          </button>
          {showEmojis && (
            <div className={styles.emojiPicker}>
              {commonEmojis.map((emoji) => (
                <span
                  key={emoji}
                  className={styles.emojiItem}
                  onClick={() => {
                    setInputText((prev) => prev + emoji);
                    setShowEmojis(false);
                  }}
                >
                  {emoji}
                </span>
              ))}
            </div>
          )}
        </div>

        <input
          type="text"
          className={styles.inputField}
          placeholder="Type a group message..."
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          maxLength={200}
          required
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              const content = inputText.trim();
              if (content === "") return;
              sendMessage({
                type: "send",
                payload: {
                  type: "group_chat",
                  receiver_Id: params.id,
                  content: content,
                },
              });
              setInputText("");
            }
          }}
        />
        <button
          className={styles.sendButton}
          onClick={(e) => {
            e.preventDefault();
            const content = inputText.trim();
            if (content === "") return;
            sendMessage({
              type: "send",
              payload: {
                type: "group_chat",
                receiver_Id: params.id,
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
