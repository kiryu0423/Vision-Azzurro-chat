import { Badge } from "@/components/ui/badge"

type Room = {
  room_id: string
  display_name: string
  is_group: boolean
  last_message_at?: string
  unread_count?: number
  last_message?: string
}

type RoomListProps = {
  rooms: Room[]
  onSelectRoom: (id: string, name: string, isGroup: boolean) => void
}

export default function RoomList({ rooms, onSelectRoom }: RoomListProps) {

  return (
    <ul className="space-y-2">
      {rooms.map((room) => (
        <li
          key={room.room_id}
          className="px-4 py-2 rounded hover:bg-blue-100 cursor-pointer"
          onClick={() => onSelectRoom(room.room_id, room.display_name, room.is_group)}
        >
          <div className="flex justify-between items-center gap-2">
          <span
            className="truncate block max-w-[160px]"
            title={room.display_name} // ✅ ホバーで全文表示
          >
            {room.display_name}
          </span>
            {(room.unread_count ?? 0) > 0 && (
              <Badge className="bg-red-500 text-white shrink-0">
                {room.unread_count}
              </Badge>
            )}
          </div>

          {room.last_message && (
            <div className="text-sm text-gray-600 truncate">{room.last_message}</div>
          )}
        </li>
      ))}
    </ul>
  )
}
