import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { WeddingInvite } from "@/components/wedding/WeddingInvite";

export function DebugPage() {
  const [id, setId] = useState("wedding-2025");
  const [user, setUser] = useState("guest");

  const previewUrl = `/${id}/${user}`;

  return (
    <div className="min-h-screen w-full p-4 sm:p-8">
      <div className="max-w-5xl mx-auto space-y-8">
        <Card>
          <CardHeader>
            <CardTitle>Invite Template Debug</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Event ID</label>
                <Input
                  value={id}
                  onChange={(e) => setId(e.target.value)}
                  placeholder="event-id"
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Guest Name</label>
                <Input
                  value={user}
                  onChange={(e) => setUser(e.target.value)}
                  placeholder="guest-name"
                />
              </div>
            </div>

            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <span>Preview URL:</span>
              <code className="rounded bg-muted px-2 py-0.5 font-mono text-xs">
                {previewUrl}
              </code>
              <Button
                variant="outline"
                size="sm"
                onClick={() => window.open(previewUrl, "_blank")}
              >
                Open
              </Button>
            </div>
          </CardContent>
        </Card>

        <div className="border rounded-lg overflow-hidden">
          <div className="bg-muted px-4 py-2 text-xs text-muted-foreground font-mono border-b">
            Preview
          </div>
          <WeddingInvite id={id} user={user} />
        </div>
      </div>
    </div>
  );
}
