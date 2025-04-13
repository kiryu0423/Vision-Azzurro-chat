// src/pages/Chat.tsx
import { useState } from "react"
import Sidebar from "@/components/Sidebar"
import ChatArea from "@/components/ChatArea"

export default function Chat() {
  const [selectedRoomId, setSelectedRoomId] = useState("")
  const [selectedRoomName, setSelectedRoomName] = useState("")

  return (
    <div className="h-screen flex">
    {/* サイドバー幅修正 */}
      <div className="w-64 min-w-[240px] border-r">
        <Sidebar
        onSelectRoom={(id, name) => {
            setSelectedRoomId(id)
            setSelectedRoomName(name)
        }}
        />
      </div>

      <div className="flex-1 overflow-hidden">
        <ChatArea roomId={selectedRoomId} roomName={selectedRoomName} />
      </div>
    </div>
  )
}
