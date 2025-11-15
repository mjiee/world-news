import { useState } from "react";
import toast from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Button, Stack, Switch, TextInput, PasswordInput, Text } from "@mantine/core";
import { useForm } from "@mantine/form";
import { useRemoteServiceStore } from "@/stores";
import { isWeb } from "@/utils/platform";
import { validateUrl } from "@/utils/url";

// Remote Service Component
export default function RemoteService() {
  const { t } = useTranslation();
  const enable = useRemoteServiceStore((state) => state.enable);
  const host = useRemoteServiceStore((state) => state.host);
  const token = useRemoteServiceStore((state) => state.token);
  const saveService = useRemoteServiceStore((state) => state.saveService);

  const [checked, setChecked] = useState(enable);

  const form = useForm({
    initialValues: { host, token },
    validate: {
      host: (value) => (!validateUrl(value) ? t("remote_service.validate.invalid_host", { ns: "settings" }) : null),
      token: (value) => (!value ? t("remote_service.validate.invalid_token", { ns: "settings" }) : null),
    },
  });

  if (isWeb()) {
    return <Text c="dimmed">{t("remote_service.web_not_support", { ns: "settings" })}</Text>;
  }

  return (
    <form
      onSubmit={form.onSubmit((values) => {
        saveService(checked, values.host, values.token);
        toast.success("ok");
      })}
    >
      <Stack gap="md">
        <TextInput label={t("remote_service.lable.service_host", { ns: "settings" })} {...form.getInputProps("host")} />
        <PasswordInput
          label={t("remote_service.lable.service_token", { ns: "settings" })}
          placeholder="0123456"
          {...form.getInputProps("token")}
        />
        <Switch
          checked={checked}
          onChange={(e) => setChecked(e.currentTarget.checked)}
          label={t("remote_service.lable.enable_remote_service", { ns: "settings" })}
        />
        <Button type="submit">{t("button.save")}</Button>
      </Stack>
    </form>
  );
}
