import { useState } from "react";
import toast from "react-hot-toast";
import { Button, Stack, Switch, TextInput } from "@mantine/core";
import { useTranslation } from "react-i18next";
import { useForm } from "@mantine/form";
import { useRemoteServiceStore } from "@/stores";
import { isWeb } from "@/utils/platform";
import { validateUrl } from "@/utils/url";

// remote service settings
export function RemoteService() {
  const { t } = useTranslation();

  const enable = useRemoteServiceStore((state) => state.enable);
  const host = useRemoteServiceStore((state) => state.host);
  const token = useRemoteServiceStore((state) => state.token);
  const saveService = useRemoteServiceStore((state) => state.saveService);

  const [checked, setChecked] = useState(enable);

  const serviceFrom = useForm({
    mode: "uncontrolled",
    initialValues: {
      host: host,
      token: token,
    },
    validate: {
      host: (value) => {
        if (validateUrl(value)) return null;
        return t("remote_service.validate.invalid_host", { ns: "settings" });
      },
      token: (value) => {
        if (value) return null;

        return t("remote_service.validate.invalid_token", { ns: "settings" });
      },
    },
  });

  // save remote service
  const saveServiceHandler = () => {
    if (!serviceFrom.validate()) return;

    saveService(checked, serviceFrom.getValues().host, serviceFrom.getValues().token);

    toast.success("ok");
  };

  return isWeb() ? (
    <p>{t("remote_service.web_not_support", { ns: "settings" })}</p>
  ) : (
    <Stack>
      <TextInput
        key={serviceFrom.key("host")}
        {...serviceFrom.getInputProps("host")}
        label={t("remote_service.lable.service_host", { ns: "settings" })}
      />
      <TextInput
        key={serviceFrom.key("token")}
        {...serviceFrom.getInputProps("token")}
        label={t("remote_service.lable.service_token", { ns: "settings" })}
        placeholder="0123456"
      />
      <Switch
        checked={checked}
        onChange={(event) => setChecked(event.currentTarget.checked)}
        label={t("remote_service.lable.enable_remote_service", { ns: "settings" })}
      />
      <Button onClick={saveServiceHandler}>{t("button.save")}</Button>
    </Stack>
  );
}
