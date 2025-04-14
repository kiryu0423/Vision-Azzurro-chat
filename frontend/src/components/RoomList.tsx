import { useEffect, useState } from "react"

type Room = {
  room_id: string
  display_name: string
  last_message_at?: string
}

type RoomListProps = {
  onSelectRoom: (id: string, name: string) => void
  refreshTrigger?: boolean
}

export default function RoomList({ onSelectRoom, refreshTrigger }: RoomListProps) {
  const [rooms, setRooms] = useState<Room[]>([])

  useEffect(() => {
    fetch("http://localhost:8081/rooms", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => setRooms(data || []))
  }, [refreshTrigger]) // ✅ 正しく再実行されるようになる！

  return (
    <ul className="space-y-2">
      {rooms
        .sort((a, b) => new Date(b.last_message_at ?? 0).getTime() - new Date(a.last_message_at ?? 0).getTime())
        .map((room) => (
          <li
            key={room.room_id}
            onClick={() => onSelectRoom(room.room_id, room.display_name)}
            className="cursor-pointer px-4 py-2 rounded hover:bg-blue-100 w-full text-left"
          >
            {room.display_name}
          </li>
        ))}
    </ul>
  )
}
