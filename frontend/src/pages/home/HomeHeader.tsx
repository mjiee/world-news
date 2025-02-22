import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import { useTranslation } from "react-i18next";
import { Button, Container, Avatar, Group, Modal } from "@mantine/core";
import { DateInput } from "@mantine/dates";
import { useDisclosure } from "@mantine/hooks";
import dayjs from "dayjs";
import "dayjs/locale/en";
import "dayjs/locale/zh";
import { LanguageSwitcher } from "@/components";
import { hasCrawlingTask, crawlingNews } from "@/services";
import styles from "@/assets/styles/header.module.css";
import appicon from "@/assets/images/appicon.png";

export function HomeHeader() {
  let navigate = useNavigate();
  const { t } = useTranslation("home");

  return (
    <header className={styles.header}>
      <Container size="md" className={styles.inner}>
        <Avatar src={appicon} variant="default" radius="sm" />
        <Group>
          <LanguageSwitcher />
          <FetchNewsButton />
          <Button onClick={() => navigate("/records")}>{t("header.button.records")}</Button>
          <Button onClick={() => navigate("/settings")}>{t("header.button.settings")}</Button>
        </Group>
      </Container>
    </header>
  );
}

// fetch news button
function FetchNewsButton() {
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
