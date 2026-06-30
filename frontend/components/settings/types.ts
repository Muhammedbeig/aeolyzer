export type Theme = "light" | "system" | "dark"

export interface SettingsProps {
  isOpen: boolean
  onClose: () => void
  theme: Theme
  onThemeChange: (theme: Theme) => void
}
