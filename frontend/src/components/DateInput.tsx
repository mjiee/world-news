import { DateInput as MantineDateInput } from "@mantine/dates";
import { useTranslation } from "react-i18next";
import "dayjs/locale/en";
import "dayjs/locale/zh";

export interface DateInputProps {
  label?: string;
  placeholder?: string;
  value: Date | null;
  onChange: (date: Date | null) => void;
}

export function DateInput(props: DateInputProps) {
  const { t, i18n } = useTranslation();

  return (
    <MantineDateInput
      maxDate={new Date()}
      locale={i18n.language}
      label={props.label}
      placeholder={props.placeholder}
      valueFormat="YYYY-MM-DD"
      value={props.value}
      onChange={props.onChange}
    />
  );
}
