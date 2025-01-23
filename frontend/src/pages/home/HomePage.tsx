import { Link } from "react-router";
import { Button } from "@mantine/core";
import { HomeHeader } from "./HomeHeader";
import { CrawlingNewsRecords } from "./CrawlingNewsRecords";

// Application homepage
export function HomePage() {
  return (
    <>
      <HomeHeader />
      <CrawlingNewsRecords />
    </>
  );
}
