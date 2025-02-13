/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_PLATFORM?: string; // "web" | "desktop"
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
