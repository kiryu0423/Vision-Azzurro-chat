import { useState, useEffect } from "react"
import RoomList from "./RoomList"
import UserList from "./UserList"
import { Button } from "@/components/ui/button"

// 型定義（Roomなど）は別途インポートするか定義してください
type Room = {
  room_id: string
  display_name: string
  last_message_at?: string
  unread_count?: number
}

type SidebarProps = {
  onSelectRoom: (id: string, name: string) => void
  userId: number
}

export default function Sidebar({ onSelectRoom, userId }: SidebarProps) {

  const [selectedUserIds, setSelectedUserIds] = useState<number[]>([])
  const [rooms, setRooms] = useState<Room[]>([])
  const [refreshRoomList, setRefreshRoomList] = useState(false)

  const createOneOnOne = async (userId: number, userName: string) => {
    const res = await fetch("http://localhost:8081/rooms", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        user_ids: [userId],
        display_name: ""
      })
    })
  
    const data = await res.json()
    if (res.ok && data.room_id) {
      setRefreshRoomList(prev => !prev)
      handleSelectRoom(data.room_id, userName)
    } else {
      alert(data.error || "個人チャット作成に失敗しました")
    }
  }  

  const createGroup = async () => {
    if (selectedUserIds.length < 2) {
      alert("2人以上選択してください")
      return
    }
  
    const res = await fetch("http://localhost:8081/rooms", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        user_ids: selectedUserIds,
        display_name: ""
      }),
    })
  
    const data = await res.json()
    if (res.ok && data.room_id) {
      setRefreshRoomList(prev => !prev)
      handleSelectRoom(data.room_id, data.display_name || "新しいグループ")
    } else {
      alert(data.error || "ルーム作成に失敗しました")
    }
  }

  // ルーム一覧の取得
  useEffect(() => {
    fetch("http://localhost:8081/rooms", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => {
        const sorted = data.sort((a, b) =>
          new Date(b.last_message_at ?? 0).getTime() -
          new Date(a.last_message_at ?? 0).getTime()
        )
        setRooms(sorted)
      })
  }, [refreshRoomList])
  
  
  const handleSelectRoom = (id: string, name: string) => {
    setRooms(prevRooms =>
      prevRooms.map(room =>
        room.room_id === id ? { ...room, unread_count: 0 } : room
      )
    )
    onSelectRoom(id, name)
  }

  // 未読管理
  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8081/ws-notify")
  
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data)
    
      setRooms((prevRooms) => {
        const updated = prevRooms.map((room) =>
          room.room_id === data.room_id
            ? {
                ...room,
                last_message_at: data.created_at,
                unread_count:
                  data.sender_id === userId
                    ? room.unread_count // 自分のメッセージは未読にしない
                    : (room.unread_count ?? 0) + 1,
              }
            : room
        )
    
        return [...updated].sort((a, b) =>
          new Date(b.last_message_at ?? 0).getTime() -
          new Date(a.last_message_at ?? 0).getTime()
        )
      })
    }
  
    return () => {
      socket.close()
    }
  }, [userId])

  return (
    <aside className="w-64 h-screen bg-gray-100 p-4 overflow-y-auto">
      <h3 className="text-lg font-bold mb-2">チャット一覧</h3>
      <RoomList rooms={rooms} onSelectRoom={handleSelectRoom} />

      <h4 className="text-md font-semibold mt-6 mb-2">ユーザー一覧（新規チャット）</h4>
      <UserList
        selectedUserIds={selectedUserIds}
        setSelectedUserIds={setSelectedUserIds}
        onCreateOneOnOne={createOneOnOne}
      />

      <Button
        className="mt-4 w-full bg-blue-500 hover:bg-blue-600 text-white"
        onClick={createGroup}
      >
        グループ作成
      </Button>
    </aside>
  )
}
