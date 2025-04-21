import { useState, useEffect, useRef } from "react"
import RoomList from "./RoomList"
import UserList from "./UserList"

// 型定義（Roomなど）は別途インポートするか定義してください
type Room = {
  room_id: string
  display_name: string
  is_group: boolean
  last_message_at?: string
  unread_count?: number
}

type SidebarProps = {
  onSelectRoom: (id: string, name: string, isGroup: boolean) => void
  userId: number
}

export default function Sidebar({ onSelectRoom, userId }: SidebarProps) {
  const [selectedUserIds, setSelectedUserIds] = useState<number[]>([])
  const [rooms, setRooms] = useState<Room[]>([])
  const [currentRoomId, setCurrentRoomId] = useState<string | null>(null)
  const [showGroupCreator, setShowGroupCreator] = useState(false)

  const httpApiUrl = import.meta.env.VITE_API_URL
  const wsUrl = httpApiUrl.replace(/^http/, "ws")

  // 初回のみルーム一覧を取得
  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    if (!token) return

    fetch(`${httpApiUrl}/rooms`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => res.json())
      .then((data) => {
        const sorted: Room[] = (data as Room[]).sort((a: Room, b: Room) =>
          new Date(b.last_message_at ?? 0).getTime() -
          new Date(a.last_message_at ?? 0).getTime()
        )
        setRooms(sorted)
      })
  }, [httpApiUrl])

  const currentRoomIdRef = useRef<string | null>(null)

  const handleSelectRoom = (id: string, name: string, isGroup: boolean) => {
    setCurrentRoomId(id)
    currentRoomIdRef.current = id
    setRooms(prevRooms =>
      prevRooms.map(room =>
        room.room_id === id ? { ...room, unread_count: 0 } : room
      )
    )
    onSelectRoom(id, name, isGroup)
  }

  // WebSocket通知でルーム更新（未読 + 最終メッセージ）
  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    const socket = new WebSocket(`${wsUrl}/ws-notify?token=${token}`)

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data)
    
      // ✅ 表示中のルームに対しては既読処理を送る
      if (data.room_id === currentRoomIdRef.current) {
        const token = localStorage.getItem("jwt_token")
        fetch(`${import.meta.env.VITE_API_URL}/rooms/${data.room_id}/read`, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        })
      }
    
      setRooms((prevRooms) => {
        const updatedRooms = prevRooms.map((room) =>
          room.room_id === data.room_id
            ? {
                ...room,
                last_message_at: data.created_at,
                last_message: data.content,
                unread_count:
                  data.sender_id === userId ||
                  data.room_id === currentRoomIdRef.current ||
                  data.from_self
                    ? room.unread_count
                    : (room.unread_count ?? 0) + 1,
              }
            : room
        )
    
        const movedToTop = updatedRooms.find((r) => r.room_id === data.room_id)
        const others = updatedRooms.filter((r) => r.room_id !== data.room_id)
    
        return movedToTop ? [movedToTop, ...others] : updatedRooms
      })
    }

    return () => {
      socket.close()
    }
  }, [userId, wsUrl])

  useEffect(() => {
    if (currentRoomId) {
      currentRoomIdRef.current = currentRoomId
    }
  }, [currentRoomId])

  const createOneOnOne = async (userId: number, userName: string) => {
    const token = localStorage.getItem("jwt_token")

    const res = await fetch(`${httpApiUrl}/rooms`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_ids: [userId],
        display_name: ""
      }),
    })

    const data = await res.json()
    if (res.ok && data.room_id) {
      handleSelectRoom(data.room_id, userName, false)
      setSelectedUserIds([])
      setShowGroupCreator(false)
    } else {
      alert(data.error || "個人チャット作成に失敗しました")
      setSelectedUserIds([])
      setShowGroupCreator(false)
    }
  }

  const createGroup = async () => {
    if (selectedUserIds.length < 2) {
      alert("2人以上選択してください")
      return
    }

    const token = localStorage.getItem("jwt_token")

    const res = await fetch(`${httpApiUrl}/rooms`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_ids: selectedUserIds,
        display_name: ""
      }),
    })

    const data = await res.json()
    if (res.ok && data.room_id) {
      handleSelectRoom(data.room_id, data.display_name || "新しいグループ", true)
      setSelectedUserIds([])
      setShowGroupCreator(false)
    } else {
      alert(data.error || "ルーム作成に失敗しました")
    }
  }

  return (
    <aside className="w-full h-screen bg-gray-100 p-4 overflow-y-auto">
      {showGroupCreator ? (
        <>
          <h3 className="text-lg font-bold mb-4">新規グループ作成</h3>

          <UserList
            selectedUserIds={selectedUserIds}
            setSelectedUserIds={setSelectedUserIds}
            onCreateOneOnOne={createOneOnOne}
          />

          <button className="mt-4 w-full" onClick={createGroup}>グループ作成</button>
          <button className="mt-2 w-full" onClick={() => setShowGroupCreator(false)}>← 戻る</button>
        </>
      ) : (
        <>
          <button className="mb-4 w-full" onClick={() => setShowGroupCreator(true)}>
            ＋ 新しいチャット
          </button>

          <h3 className="text-lg font-bold mb-2">チャット一覧</h3>
          <RoomList rooms={rooms} onSelectRoom={(id, name, isGroup) => handleSelectRoom(id, name, isGroup)} />
        </>
      )}
    </aside>
  )
}
