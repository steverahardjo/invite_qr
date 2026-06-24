import { useParams, Navigate } from "react-router-dom";
import { WeddingInvite } from "@/components/wedding/WeddingInvite";

export function InvitePage() {
  const { id, user } = useParams<{ id: string; user: string }>();

  if (!id || !user) {
    return <Navigate to="/" replace />;
  }

  return <WeddingInvite id={id} user={user} />;
}
