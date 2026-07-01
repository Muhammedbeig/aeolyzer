import { FileText, Award, Users, List, Mic, Brain } from "lucide-react"

interface KnowledgeBaseProps {
  activeSection: string
}

export function AeolyzerKnowledgeBase({ activeSection }: KnowledgeBaseProps) {
  const renderContent = () => {
    switch (activeSection) {
      case "profile":
        return (
          <div className="space-y-6 max-w-3xl">
            <div>
              <h2 className="text-xl font-semibold mb-1">Knowledge Base Profile</h2>
              <p className="text-sm text-muted-foreground">Manage your main profile details, company information, and primary directives.</p>
            </div>
            
            <div className="grid gap-6">
              <div className="space-y-2">
                <label className="text-sm font-medium">Company/Agent Name</label>
                <input 
                  type="text" 
                  placeholder="e.g. AEOlyzer SEO Expert" 
                  className="w-full px-3 py-2 bg-transparent border border-black/10 dark:border-white/10 rounded-md focus:outline-none focus:ring-1 focus:ring-black/20 dark:focus:ring-white/20 transition-all" 
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Core Description</label>
                <textarea 
                  placeholder="Describe the primary purpose and scope of this knowledge base..." 
                  className="w-full min-h-[120px] px-3 py-2 bg-transparent border border-black/10 dark:border-white/10 rounded-md focus:outline-none focus:ring-1 focus:ring-black/20 dark:focus:ring-white/20 transition-all resize-y" 
                />
              </div>
              <div className="pt-4">
                <button className="px-4 py-2 bg-black dark:bg-white text-white dark:text-black rounded-md font-medium text-sm hover:opacity-90 transition-opacity">
                  Save Profile
                </button>
              </div>
            </div>
          </div>
        )
      case "eeat":
        return (
          <div className="space-y-6 max-w-3xl">
            <div>
              <h2 className="text-xl font-semibold mb-1">E-E-A-T Guidelines</h2>
              <p className="text-sm text-muted-foreground">Configure Experience, Expertise, Authoritativeness, and Trustworthiness signals.</p>
            </div>
            <div className="p-8 border border-black/10 dark:border-white/10 rounded-xl bg-black/[0.02] dark:bg-white/[0.02] flex flex-col items-center justify-center text-center space-y-3">
              <Award className="w-8 h-8 text-muted-foreground opacity-50" />
              <div>
                <h3 className="font-medium">No E-E-A-T Rules Defined</h3>
                <p className="text-sm text-muted-foreground max-w-sm mt-1">Add your specific authority markers and trust signals to ensure content meets Google&apos;s quality rater guidelines.</p>
              </div>
              <button className="mt-2 px-4 py-2 border border-black/10 dark:border-white/10 rounded-md text-sm font-medium hover:bg-black/5 dark:hover:bg-white/5 transition-colors">
                Add Guideline
              </button>
            </div>
          </div>
        )
      case "competitors":
        return (
          <div className="space-y-6 max-w-3xl">
            <div>
              <h2 className="text-xl font-semibold mb-1">Competitor Analysis</h2>
              <p className="text-sm text-muted-foreground">Track competitor websites and analyze their content strategies.</p>
            </div>
            <div className="flex gap-2">
              <input 
                type="text" 
                placeholder="https://competitor.com" 
                className="flex-1 px-3 py-2 bg-transparent border border-black/10 dark:border-white/10 rounded-md focus:outline-none focus:ring-1 focus:ring-black/20 dark:focus:ring-white/20 transition-all" 
              />
              <button className="px-4 py-2 bg-black dark:bg-white text-white dark:text-black rounded-md font-medium text-sm hover:opacity-90 transition-opacity">
                Add
              </button>
            </div>
            <div className="p-8 border border-black/10 dark:border-white/10 rounded-xl bg-black/[0.02] dark:bg-white/[0.02] flex flex-col items-center justify-center text-center space-y-3">
              <Users className="w-8 h-8 text-muted-foreground opacity-50" />
              <div>
                <p className="text-sm text-muted-foreground max-w-sm">No competitors added yet.</p>
              </div>
            </div>
          </div>
        )
      case "topics":
        return (
          <div className="space-y-6 max-w-3xl">
            <div>
              <h2 className="text-xl font-semibold mb-1">Topic Clusters</h2>
              <p className="text-sm text-muted-foreground">Manage your core content pillars and keyword clusters.</p>
            </div>
            <div className="p-8 border border-black/10 dark:border-white/10 rounded-xl bg-black/[0.02] dark:bg-white/[0.02] flex flex-col items-center justify-center text-center space-y-3">
              <List className="w-8 h-8 text-muted-foreground opacity-50" />
              <div>
                <h3 className="font-medium">Map your Content</h3>
                <p className="text-sm text-muted-foreground max-w-sm mt-1">Define the main topics your AI should focus on for semantic relevance.</p>
              </div>
              <button className="mt-2 px-4 py-2 border border-black/10 dark:border-white/10 rounded-md text-sm font-medium hover:bg-black/5 dark:hover:bg-white/5 transition-colors">
                Create Cluster
              </button>
            </div>
          </div>
        )
      case "tone":
        return (
          <div className="space-y-6 max-w-3xl">
            <div>
              <h2 className="text-xl font-semibold mb-1">Brand Tone & Voice</h2>
              <p className="text-sm text-muted-foreground">Set the personality, reading level, and style for generated content.</p>
            </div>
            <div className="grid gap-6">
              <div className="space-y-2">
                <label className="text-sm font-medium">Primary Tone</label>
                <select className="w-full px-3 py-2 bg-transparent border border-black/10 dark:border-white/10 rounded-md focus:outline-none focus:ring-1 focus:ring-black/20 dark:focus:ring-white/20 transition-all appearance-none cursor-pointer">
                  <option>Professional & Authoritative</option>
                  <option>Conversational & Friendly</option>
                  <option>Academic & Technical</option>
                  <option>Persuasive & Direct</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Custom Voice Instructions</label>
                <textarea 
                  placeholder="e.g. Always use active voice, avoid jargon, use short paragraphs..." 
                  className="w-full min-h-[120px] px-3 py-2 bg-transparent border border-black/10 dark:border-white/10 rounded-md focus:outline-none focus:ring-1 focus:ring-black/20 dark:focus:ring-white/20 transition-all resize-y" 
                />
              </div>
              <div className="pt-4">
                <button className="px-4 py-2 bg-black dark:bg-white text-white dark:text-black rounded-md font-medium text-sm hover:opacity-90 transition-opacity">
                  Save Tone Settings
                </button>
              </div>
            </div>
          </div>
        )
      case "memory":
        return (
          <div className="space-y-6 max-w-3xl">
            <div>
              <h2 className="text-xl font-semibold mb-1">Agent Memory</h2>
              <p className="text-sm text-muted-foreground">Manage facts, past interactions, and contextual knowledge the AI retains.</p>
            </div>
            <div className="p-8 border border-black/10 dark:border-white/10 rounded-xl bg-black/[0.02] dark:bg-white/[0.02] flex flex-col items-center justify-center text-center space-y-3">
              <Brain className="w-8 h-8 text-muted-foreground opacity-50" />
              <div>
                <h3 className="font-medium">Memory Vault is Empty</h3>
                <p className="text-sm text-muted-foreground max-w-sm mt-1">The agent will automatically save important facts here during conversations, or you can add them manually.</p>
              </div>
              <button className="mt-2 px-4 py-2 border border-black/10 dark:border-white/10 rounded-md text-sm font-medium hover:bg-black/5 dark:hover:bg-white/5 transition-colors">
                Add Fact
              </button>
            </div>
          </div>
        )
      default:
        return <div>Select a section from the sidebar.</div>
    }
  }

  return (
    <div className="flex-1 w-full p-6 md:p-8 overflow-y-auto custom-scrollbar font-outfit">
      <div className="max-w-5xl mx-auto">
        {renderContent()}
      </div>
    </div>
  )
}
