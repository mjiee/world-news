import { httpx } from "wailsjs/go/models";

// getPageNumbers to get page numbers
export function getPageNumber(page: httpx.Pagination): number {
  return Math.ceil((page.total ?? 0) / (page.limit ?? 1));
}

// getPageData to get page data
export function getPageData<T>(arr: T[], page: number, limit: number): T[] {
  if (!Array.isArray(arr)) return [];

  if (arr.length <= limit) return arr;

  const startIndex = (page - 1) * limit;
  const endIndex = startIndex + limit;

  return arr.slice(startIndex, endIndex);
}
