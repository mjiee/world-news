import { Center, Loader, Stack } from "@mantine/core";

export function Loading() {
  return (
    <Center h={300}>
      <Stack align="center" gap="md">
        <Loader color="blue" size="xl" type="bars" />
      </Stack>
    </Center>
  );
}
