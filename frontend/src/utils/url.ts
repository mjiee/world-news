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
