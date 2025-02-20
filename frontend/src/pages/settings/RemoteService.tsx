import { useState } from "react";
import toast from "react-hot-toast";
import { Button, Stack, Switch, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useField } from "@mantine/form";
import { useRemoteServiceStore } from "@/stores";
import { isWeb } from "@/utils/platform";
import { validateUrl } from "@/utils/validate";

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
      return t("remote_service.validate.invalid_host", { ns: "settings" });
    },
  });

  // save remote service
  const saveServiceHandler = () => {
    if (!newHost.validate()) return;

    saveService(checked, newHost.getValue());

    toast.success("ok");
  };

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
      <Button onClick={saveServiceHandler}>{t("button.save")}</Button>
    </Stack>
  );
}
