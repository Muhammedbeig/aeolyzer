export const SETTINGS_TABS = [
  "General",
  "Account",
  "Privacy",
  "Billing",
  "Capabilities",
  "Connectors",
  "Aeolyzer Code",
] as const;

export type SettingsTab = typeof SETTINGS_TABS[number];
