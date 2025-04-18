// src/pages/Chat.tsx
import { useEffect, useState } from "react"
import Sidebar from "@/components/Sidebar"
import ChatArea from "@/components/ChatArea"

export default function Chat() {
  const [selectedRoomId, setSelectedRoomId] = useState("")
  const [selectedRoomName, setSelectedRoomName] = useState("")
  const [userId, setUserId] = useState<number | null>(null)
  const [selectedIsGroup, setSelectedIsGroup] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    if (!token) {
      window.location.href = "/"
      return
    }

    fetch("http://localhost:8081/me", {
      headers: {
        "Authorization": `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("unauthorized")
        return res.json()
      })
      .then((data) => {
        if (data?.user_id) setUserId(data.user_id)
        else throw new Error("user not found")
      })
      .catch(() => {
        localStorage.removeItem("jwt_token")
        window.location.href = "/"
      })
  }, [])

  if (!userId) return <div>読み込み中...</div>

  return (
    <div className="h-screen flex">
      <div className="w-[300px] min-w-[300px] border-r">
        <Sidebar
          userId={userId}
          onSelectRoom={(id, name, isGroup) => {
            setSelectedRoomId(id)
            setSelectedRoomName(name)
            setSelectedIsGroup(isGroup)
          }}
        />
      </div>

      <div className="flex-1 overflow-hidden">
        <ChatArea
          roomId={selectedRoomId}
          roomName={selectedRoomName}
          userId={userId}
          isGroup={selectedIsGroup}
        />
      </div>
    </div>
  )
}
