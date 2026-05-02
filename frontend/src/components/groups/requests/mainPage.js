"use client";

import { useEffect, useState } from "react";
import styles from "@/components/groups/styles/singleGroup.module.css";
import RequestCard from "@/components/groups/requests/getjoinrequest";
import handleFetchData from "@/components/groups/utils/FetchData";
import { useRouter } from "next/navigation";

export default function GroupRequests({ id }) {
  const [requests, setRequests] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    handleFetchData(
      requests,
      setRequests,
      setLoading,
      setHasMore,
      "requests",
      id,
      router,
    );
  }, [id]);

  return (
    <>
      {requests.length === 0 ? (
        <p>There is no requests for now but you can invite people !</p>
      ) : (
        <>
          {requests.map((u) => (
            <RequestCard
              key={u.id}
              user={u}
              groupId={id}
              setRequests={setRequests}
            />
          ))}
        </>
      )}

      {hasMore && (
        <button
          className={styles.loadMoreBtn}
          onClick={() =>
            handleFetchData(
              requests,
              setRequests,
              setLoading,
              setHasMore,
              "requests",
              id,
              router,
            )
          }
          disabled={loading}
        >
          {loading ? "Loading..." : "Load More"}
        </button>
      )}
    </>
  );
}
