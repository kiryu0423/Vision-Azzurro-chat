import { Pencil } from "lucide-react" // ✏️ アイコン用
import { Badge } from "@/components/ui/badge"

type Room = {
  room_id: string
  display_name: string
  last_message_at?: string
  unread_count?: number
}

type RoomListProps = {
  rooms: Room[]
  onSelectRoom: (id: string, name: string) => void
}

export default function RoomList({ rooms, onSelectRoom }: RoomListProps) {
  const handleEdit = async (roomId: string) => {
    const newName = prompt("新しいグループ名を入力")
    if (!newName) return

    await fetch(`http://localhost:8081/rooms/${roomId}/name`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ display_name: newName }),
    })

    // 親に再取得させる（Sidebarで useEffect 依存など）
    location.reload() // 簡易対応（setRefreshRoomListでもOK）
  }

  return (
    <ul className="space-y-2">
      {rooms.map((room) => (
        <li
          key={room.room_id}
          className="flex items-center justify-between px-4 py-2 rounded hover:bg-blue-100 cursor-pointer"
          onClick={() => onSelectRoom(room.room_id, room.display_name)}
        >
          <span>{room.display_name}</span>

          <div className="flex gap-2 items-center">
            {(room.unread_count ?? 0) > 0 && (
              <Badge className="bg-red-500 text-white">{room.unread_count}</Badge>
            )}
            <button
              className="text-gray-500 hover:text-gray-700"
              onClick={(e) => {
                e.stopPropagation()
                handleEdit(room.room_id)
              }}
              title="名前を編集"
            >
              <Pencil size={16} />
            </button>
          </div>
        </li>
      ))}
    </ul>
  )
}
