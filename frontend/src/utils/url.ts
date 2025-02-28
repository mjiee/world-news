import { useServiceHost } from "@/stores";
import { isWeb } from "./platform";

// validateUrl validate url
export const validateUrl = (value: string | undefined) => {
  if (!value) return false;

  try {
    new URL(value);
    return true;
  } catch (_) {
    return false;
  }
};

// getHost get host
export const getHost = (value: string | undefined) => {
  if (!value) return "";

  try {
    const url = new URL(value);
    return url.host;
  } catch (_) {
    return value;
  }
};

// setHost set host to url
export function setHost(url: string): string {
  let host = useServiceHost();

  if (isWeb()) host = import.meta.env.VITE_SERVICE_HOST;

  if (!host) return url;

  return new URL(url, host).href;
}
