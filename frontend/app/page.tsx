"use client"

import { useState, useCallback, useEffect } from "react"
import { useTheme } from "next-themes"
import { Menu, Plus } from "lucide-react"
import { AeolyzerSidebar } from "@/components/sidebar/aeolyzer-sidebar"
import { AeolyzerChatArea } from "@/components/chat/chat-area"
import { AeolyzerWelcome } from "@/components/chat/welcome-screen"
import { AeolyzerChatInput } from "@/components/chat/chat-input"
import { AeolyzerSettings } from "@/components/settings/aeolyzer-settings"

type Theme = "light" | "system" | "dark"

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  isStreaming?: boolean
}

// Sample responses for demo
const sampleResponses = [
  `Here's a honest look at where I genuinely shine:

**Writing & Communication**

I'm strong at adapting tone and style — whether that's a formal business report, a casual blog post, a persuasive essay, or creative fiction. I can help draft, edit, restructure, and polish almost any kind of text.

**Reasoning & Analysis**

Breaking down complex problems, weighing pros and cons, spotting logical gaps, and thinking through multi-step scenarios are things I do well. Give me a messy situation and I'll help you find clarity.

**Code & Technical Work**

I can write, debug, and explain code across many languages. I'm particularly helpful for understanding concepts, reviewing logic, and working through implementation challenges.

**Research & Synthesis**

I'm good at meeting people where they are — breaking down difficult concepts for beginners or going deep for experts, adjusting as needed.

**Where I'm less reliable:** Very recent news (though I can search the web), tasks requiring physical actions, and highly specialized professional advice (legal, medical, financial) where you should always consult a qualified human.

What are you working on? I can give you a better sense of how I can help with your specific situation.`,
  `That's a great question! Let me break it down for you.

**The Core Concept**

At its heart, this is about understanding the fundamental principles and applying them consistently. Once you grasp the basics, everything else builds naturally.

**Key Points to Remember**

1. Start with the foundation — don't skip steps
2. Practice regularly to build muscle memory
3. Learn from mistakes — they're your best teachers
4. Connect concepts together to see the bigger picture

**Practical Application**

The best way to truly understand something is to apply it in real-world scenarios. Theory is important, but hands-on experience is invaluable.

Would you like me to elaborate on any specific aspect?`,
  `I'd be happy to help you with that! Here's my approach:

**Understanding the Problem**

First, let's make sure we're solving the right problem. Often, the initial question is just the surface of a deeper challenge.

**Step-by-Step Solution**

I'll walk you through this methodically:

1. Identify the core requirements
2. Break down the complexity into manageable pieces
3. Address each piece systematically
4. Validate the solution against your original goals

**Important Considerations**

Keep in mind that context matters significantly here. What works in one situation may need adjustment in another.

Is there anything specific you'd like me to focus on?`
]

export default function AeolyzerChatbot() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [activeTab, setActiveTab] = useState("Home")
  const [messages, setMessages] = useState<Message[]>([])
  const [isGenerating, setIsGenerating] = useState(false)
  const [chatTitle, setChatTitle] = useState<string | undefined>()
  const [settingsOpen, setSettingsOpen] = useState(false)
  const { theme, setTheme } = useTheme()

  // Auto-close sidebar on mobile on initial load
  useEffect(() => {
    if (window.innerWidth < 768) {
      setSidebarOpen(false)
    }
  }, [])

  const handleNewChat = useCallback(() => {
    setMessages([])
    setChatTitle(undefined)
    setIsGenerating(false)
  }, [])

  const handleSendMessage = useCallback(async (content: string) => {
    // Add user message
    const userMessage: Message = {
      id: `user-${Date.now()}`,
      role: "user",
      content,
    }
    
    setMessages(prev => [...prev, userMessage])
    setIsGenerating(true)

    // Set chat title from first message
    if (!chatTitle) {
      const title = content.length > 30 ? content.slice(0, 30) + "..." : content
      setChatTitle(title)
    }

    // Simulate thinking delay
    await new Promise(resolve => setTimeout(resolve, 1000))

    // Select a random response
    const responseContent = sampleResponses[Math.floor(Math.random() * sampleResponses.length)]

    // Add assistant message with streaming
    const assistantMessage: Message = {
      id: `assistant-${Date.now()}`,
      role: "assistant",
      content: responseContent,
      isStreaming: true,
    }

    setMessages(prev => [...prev, assistantMessage])

    // Calculate approximate streaming time based on content length
    const streamingTime = Math.min(responseContent.length * 15, 5000)
    
    // After streaming is done, update the message to not be streaming
    setTimeout(() => {
      setMessages(prev => 
        prev.map(msg => 
          msg.id === assistantMessage.id 
            ? { ...msg, isStreaming: false }
            : msg
        )
      )
      setIsGenerating(false)
    }, streamingTime)
  }, [chatTitle])

  return (
    <div 
      className="flex h-screen overflow-hidden bg-background"
    >
      {/* Sidebar */}
      <AeolyzerSidebar 
        isOpen={sidebarOpen} 
        onToggle={() => setSidebarOpen(!sidebarOpen)}
        onNewChat={handleNewChat}
        currentChatTitle={chatTitle}
        activeTab={activeTab}
        onTabChange={setActiveTab}
        onOpenSettings={() => setSettingsOpen(true)}
      />

      {/* Vertical divider */}
      <div className="w-px bg-border" />

      
      {/* Main content */}
      <main className="flex-1 flex flex-col min-w-0 min-h-0 relative">
        {/* Mobile Header */}
        <div className="md:hidden flex items-center justify-between p-3 border-b border-border bg-background flex-shrink-0 z-10">
          <button 
            onClick={() => setSidebarOpen(true)}
            className="p-1.5 rounded-md text-muted-foreground hover:bg-muted"
            aria-label="Open sidebar"
          >
            <Menu size={20} strokeWidth={1.5} />
          </button>
          <span className="text-sm font-medium text-foreground">AEOlyzer</span>
          <button 
            onClick={handleNewChat}
            className="p-1.5 rounded-md text-muted-foreground hover:bg-muted"
            aria-label="New chat"
          >
            <Plus size={20} strokeWidth={1.5} />
          </button>
        </div>

        {/* Chat area or welcome screen */}
        {activeTab === 'Home' ? (
          <div className="flex-1 flex flex-col bg-background overflow-y-auto">
            <AeolyzerWelcome onSend={(msg) => { handleSendMessage(msg); setActiveTab('Agent'); }} isGenerating={isGenerating} />
          </div>
        ) : (
          <>
            <AeolyzerChatArea 
              messages={messages} 
              isGenerating={isGenerating}
              chatTitle={chatTitle}
              onSend={handleSendMessage}
            />

            {/* Input area - only shown when there are messages */}
            {messages.length > 0 && (
              <div 
                className="flex-shrink-0 px-4 pb-4 pt-2 bg-background"
              >
                <AeolyzerChatInput 
                  onSend={handleSendMessage}
                  isGenerating={isGenerating}
                  placeholder="Reply..."
                />
              </div>
            )}
          </>
        )}
      </main>

      {/* Settings Modal */}
      <AeolyzerSettings 
        isOpen={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        theme={(theme as Theme) || "system"}
        onThemeChange={(t) => setTheme(t)}
      />
    </div>
  )
}
