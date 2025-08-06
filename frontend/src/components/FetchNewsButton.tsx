import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Button, Group, Modal } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { hasCrawlingTask, crawlingNews } from "@/services";
import { DateInput } from "./DateInput";

// fetch news button
export function FetchNewsButton() {
  const { t, i18n } = useTranslation();
  const [opened, { open, close }] = useDisclosure(false);
  const [startTime, setStartTime] = useState<string | null>(null);
  const [disabled, setDisabled] = useState<boolean>(false);

  // click ok handler
  const clickOkHandler = () => {
    crawlingNews({ startTime: startTime ?? "" });
    setStartTime(null);
    close();
  };

  // click cancel handler
  const clickCancelHandler = () => {
    setStartTime(null);
    close();
  };

  // check can fetch news
  const checkCanFetchNews = async () => {
    const processingTask = await hasCrawlingTask();

    if (processingTask) setDisabled(true);
    else setDisabled(false);
  };

  useEffect(() => {
    checkCanFetchNews();
  }, [opened]);

  return (
    <>
      <Modal opened={opened} onClose={close} withCloseButton={false}>
        <DateInput label={t("header.label.start_time", { ns: "home" })} onChange={setStartTime} />
        <Group justify="flex-end" mt="md">
          <Button disabled={disabled} type="submit" onClick={clickOkHandler}>
            {t("button.ok")}
          </Button>
          <Button onClick={clickCancelHandler} variant="default">
            {t("button.cancel")}
          </Button>
        </Group>
      </Modal>
      <Button onClick={open}>{t("header.button.fetch_news", { ns: "home" })}</Button>
    </>
  );
}
