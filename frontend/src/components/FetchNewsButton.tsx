import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Button, Group, Modal } from "@mantine/core";
import { DateInput } from "@mantine/dates";
import { useDisclosure } from "@mantine/hooks";
import dayjs from "dayjs";
import "dayjs/locale/en";
import "dayjs/locale/zh";
import { hasCrawlingTask, crawlingNews } from "@/services";

// fetch news button
export function FetchNewsButton() {
  const { t, i18n } = useTranslation();
  const [opened, { open, close }] = useDisclosure(false);
  const [startTime, setStartTime] = useState<Date | null>(null);
  const [disabled, setDisabled] = useState<boolean>(false);

  // click ok handler
  const clickOkHandler = () => {
    crawlingNews({ startTime: startTime ? dayjs(startTime).format("YYYY-MM-DD HH:mm:ss") : "" });
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
        <DateInput
          maxDate={new Date()}
          locale={i18n.language}
          label={t("header.label.start_time", { ns: "home" })}
          valueFormat="YYYY-MM-DD"
          value={startTime}
          onChange={setStartTime}
        />
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
