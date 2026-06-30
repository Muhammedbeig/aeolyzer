export interface SidebarProps {
  activeTab: string
  onTabChange: (tab: string) => void
  isOpen: boolean
  onToggle: () => void
  onNewChat: () => void
  currentChatTitle?: string
  onOpenSettings: () => void
}
