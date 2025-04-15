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
          </div>
        </li>
      ))}
    </ul>
  )
}
