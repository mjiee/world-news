// Is web platform
export function isWeb(): boolean {
  return import.meta.env.VITE_PLATFORM === "web";
}

// Is desktop platform
export function isDesktop(): boolean {
  return import.meta.env.VITE_PLATFORM === "desktop";
}
