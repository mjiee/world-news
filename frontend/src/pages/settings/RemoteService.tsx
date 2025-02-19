import { useState } from "react";
import toast from "react-hot-toast";
import { Button, Stack, Switch, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useField } from "@mantine/form";
import { useRemoteServiceStore } from "@/stores";
import { isWeb } from "@/utils/platform";

// remote service settings
export function RemoteService() {
  const { t } = useTranslation();

  const enable = useRemoteServiceStore((state) => state.enable);
  const host = useRemoteServiceStore((state) => state.host);
  const saveService = useRemoteServiceStore((state) => state.saveService);

  const [checked, setChecked] = useState(enable);

  const newHost = useField({
    initialValue: host,
    validateOnChange: true,
    validate: (value) => {
      if (validateUrl(value)) return null;
      return t("remote_service.validate.host_invalid", { ns: "settings" });
    },
  });

  return isWeb() ? (
    <p>{t("remote_service.web_not_support", { ns: "settings" })}</p>
  ) : (
    <Stack>
      <TextInput {...newHost.getInputProps()} label={t("remote_service.lable.service_host", { ns: "settings" })} />
      <Switch
        checked={checked}
        onChange={(event) => setChecked(event.currentTarget.checked)}
        label={t("remote_service.lable.enable_remote_service", { ns: "settings" })}
      />
      <Button
        onClick={() => {
          if (!newHost.validate()) return;

          saveService(checked, newHost.getValue());

          toast.success("ok");
        }}
      >
        {t("button.save")}
      </Button>
    </Stack>
  );
}

// validateUrl validate url
const validateUrl = (value: string | undefined) => {
  if (!value) return false;

  try {
    new URL(value);
    return true;
  } catch (_) {
    return false;
  }
};
