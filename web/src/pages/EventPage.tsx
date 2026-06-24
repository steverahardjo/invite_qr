import { useParams } from "react-router-dom";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function EventPage() {
  const { id } = useParams<{ id: string }>();

  return (
    <div className="min-h-screen w-full flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        <Card>
          <CardHeader>
            <CardTitle className="text-center">Wedding Event</CardTitle>
          </CardHeader>
          <CardContent className="text-center space-y-4">
            <p className="text-muted-foreground">
              Select your invitation to view details.
            </p>
            <p className="text-sm text-muted-foreground">
              Event ID: <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{id}</code>
            </p>
            <p className="text-xs text-muted-foreground">
              Visit <code className="rounded bg-muted px-1 py-0.5 font-mono">/{id}/your-name</code> to view your invite
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
