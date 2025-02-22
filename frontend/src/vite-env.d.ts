/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

interface ImportMetaEnv {
  readonly VITE_PLATFORM?: string; // "web" | "desktop"
  readonly VITE_SERVICE_HOST?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
