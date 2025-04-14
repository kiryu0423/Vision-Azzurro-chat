// src/pages/Chat.tsx
import { useEffect, useState } from "react"
import Sidebar from "@/components/Sidebar"
import ChatArea from "@/components/ChatArea"

export default function Chat() {
  const [selectedRoomId, setSelectedRoomId] = useState("")
  const [selectedRoomName, setSelectedRoomName] = useState("")
  const [userId, setUserId] = useState<number | null>(null)

  useEffect(() => {
    fetch("http://localhost:8081/me", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => {
        if (data?.user_id) setUserId(data.user_id)
      })
      .catch(() => {
        window.location.href = "/"
      })
  }, [])

  if (!userId) return <div>読み込み中...</div>

  return (
    <div className="h-screen flex">
    {/* サイドバー幅修正 */}
      <div className="w-64 min-w-[240px] border-r">
        <Sidebar
        userId={userId}
        onSelectRoom={(id, name) => {
            setSelectedRoomId(id)
            setSelectedRoomName(name)
        }}
        />
      </div>

      <div className="flex-1 overflow-hidden">
        <ChatArea roomId={selectedRoomId} roomName={selectedRoomName} userId={userId}/>
      </div>
    </div>
  )
}
