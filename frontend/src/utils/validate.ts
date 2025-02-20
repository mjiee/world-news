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
