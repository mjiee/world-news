import { Link } from "react-router";
import { Button } from "@mantine/core";

// Application homepage
export function HomePage() {
  return (
    <div>
      <Link to="/settings">
        <Button>Settings</Button>
      </Link>
    </div>
  );
}
