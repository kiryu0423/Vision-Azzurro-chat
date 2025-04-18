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
  const [refreshRoomList, setRefreshRoomList] = useState(false)
  const [currentRoomId, setCurrentRoomId] = useState<string | null>(null)
  const [showGroupCreator, setShowGroupCreator] = useState(false)


  const createOneOnOne = async (userId: number, userName: string) => {
    const token = localStorage.getItem("jwt_token")

    const res = await fetch("http://localhost:8081/rooms", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_ids: [userId], // または selectedUserIds
        display_name: ""
      }),
    })

    const data = await res.json()
    if (res.ok && data.room_id) {
      setRefreshRoomList(prev => !prev)
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

    const res = await fetch("http://localhost:8081/rooms", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`,
      },
      body: JSON.stringify({
        user_ids: [userId], // または selectedUserIds
        display_name: ""
      }),
    })

    const data = await res.json()
    if (res.ok && data.room_id) {
      setRefreshRoomList(prev => !prev)
      handleSelectRoom(data.room_id, data.display_name || "新しいグループ", true)
      setSelectedUserIds([])
      setShowGroupCreator(false) // ✅ ← これを追加して戻る
    } else {
      alert(data.error || "ルーム作成に失敗しました")
    }
  }

  // ルーム一覧の取得
  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    if (!token) return

    fetch("http://localhost:8081/rooms", {
      headers: {
        "Authorization": `Bearer ${token}`,
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
  }, [refreshRoomList])


  // 定期的にルーム一覧を再取得（ポーリング）
  useEffect(() => {
    const interval = setInterval(() => {
      const token = localStorage.getItem("jwt_token")
      if (!token) return

      fetch("http://localhost:8081/rooms", {
        headers: {
          "Authorization": `Bearer ${token}`,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          setRooms(prevRooms => mergeRoomList(data, prevRooms))
        })
    }, 5000)

    return () => clearInterval(interval)
  }, [])

  const currentRoomIdRef = useRef<string | null>(null)

  const handleSelectRoom = (id: string, name: string, isGroup: boolean) => {
    setCurrentRoomId(id)
    currentRoomIdRef.current = id // ← ここで反映
    setRooms(prevRooms =>
      prevRooms.map(room =>
        room.room_id === id ? { ...room, unread_count: 0 } : room
      )
    )
    onSelectRoom(id, name, isGroup)
  }

  // 未読管理
  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    const socket = new WebSocket(`ws://localhost:8081/ws-notify?token=${token}`)
  
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data)
    
      setRooms((prevRooms) => {
        const updatedRooms = prevRooms.map((room) =>
          room.room_id === data.room_id
            ? {
                ...room,
                last_message_at: data.created_at,
                unread_count:
                  data.sender_id === userId ||
                  data.room_id === currentRoomIdRef.current ||
                  data.from_self
                    ? room.unread_count
                    : (room.unread_count ?? 0) + 1,
              }
            : room
        )
      
        // そのルームを先頭に移動（他の順序を変えない）
        const movedToTop = updatedRooms.find((r) => r.room_id === data.room_id)
        const others = updatedRooms.filter((r) => r.room_id !== data.room_id)
      
        return movedToTop ? [movedToTop, ...others] : updatedRooms
      })
    }
  
    return () => {
      socket.close()
    }
  }, [userId])

  useEffect(() => {
    if (currentRoomId) {
      currentRoomIdRef.current = currentRoomId
    }
  }, [currentRoomId])  

  // ルームの管理（通知など）
  const mergeRoomList = (fetched: Room[], current: Room[]) => {
    const currentMap = new Map(current.map(r => [r.room_id, r]))
  
    return fetched.map(fetchedRoom => {
      const currentRoom = currentMap.get(fetchedRoom.room_id)
  
      return {
        ...fetchedRoom,
        unread_count: currentRoom?.unread_count ?? 0, // ✅ ポーリングでは上書きしない
      }
    }).sort((a, b) =>
      new Date(b.last_message_at ?? 0).getTime() - new Date(a.last_message_at ?? 0).getTime()
    )
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

        <button
          className="mt-4 w-full"
          onClick={createGroup}
        >
          グループ作成
        </button>

        <button
          className="mt-2 w-full"
          onClick={() => setShowGroupCreator(false)}
        >
          ← 戻る
        </button>
      </>
    ) : (
      <>
        <button
          className="mb-4 w-full"
          onClick={() => setShowGroupCreator(true)}
        >
          ＋ 新しいチャット
        </button>

        <h3 className="text-lg font-bold mb-2">チャット一覧</h3>
        <RoomList rooms={rooms} onSelectRoom={(id, name, isGroup) => handleSelectRoom(id, name, isGroup)} />
      </>
    )}
  </aside>
  )
}
