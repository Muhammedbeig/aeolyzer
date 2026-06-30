import { MessageSquare } from "lucide-react"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog"
import { SEARCH_RESULTS } from "./constants"

interface SearchDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SearchDialog({ open, onOpenChange }: SearchDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="overflow-hidden p-0 shadow-xl border-border bg-sidebar-bg sm:max-w-[600px] gap-0">
        <DialogTitle className="sr-only">Search History</DialogTitle>
        <Command className="bg-transparent [&_[cmdk-group-heading]]:text-sidebar-muted w-full flex flex-col">
          <CommandInput 
            placeholder="Search chats and projects" 
            className="border-b-[0.5px] border-black/10 dark:border-white/10 text-sidebar-text"
          />
          <CommandList className="max-h-[60vh] overflow-y-auto">
            <CommandEmpty className="py-6 text-center text-sm text-sidebar-muted">No results found.</CommandEmpty>
            <CommandGroup className="p-2 text-sidebar-text">
              {SEARCH_RESULTS.map((item, i) => (
                <CommandItem 
                  key={i} 
                  onSelect={() => onOpenChange(false)}
                  className="flex items-center justify-between px-3 py-3 cursor-pointer group rounded-lg data-[selected=true]:bg-sidebar-hover data-[selected=true]:text-sidebar-text text-sidebar-muted transition-colors"
                >
                  <div className="flex items-center gap-3 truncate">
                    <MessageSquare size={16} strokeWidth={1.5} className="flex-shrink-0" />
                    <span className="truncate text-[15px]">{item.title}</span>
                  </div>
                  <div className="flex items-center flex-shrink-0 text-[13px] text-sidebar-muted pl-4">
                    <span className="block group-data-[selected=true]:hidden">{item.time}</span>
                    <span className="hidden group-data-[selected=true]:block font-medium">Enter</span>
                  </div>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </DialogContent>
    </Dialog>
  )
}
