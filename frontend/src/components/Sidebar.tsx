// src/components/Sidebar.tsx
import RoomList from "./RoomList"
import UserList from "./UserList"
import { Button } from "@/components/ui/button"

type SidebarProps = {
  onSelectRoom: (id: string, name: string) => void
}

export default function Sidebar({ onSelectRoom }: SidebarProps) {
  return (
    <aside className="w-full bg-gray-100 p-4 overflow-y-auto">
      <h3 className="text-lg font-bold mb-2">チャット一覧</h3>
      <RoomList onSelectRoom={onSelectRoom} />

      <h4 className="text-md font-semibold mt-6 mb-2">ユーザー一覧（新規チャット）</h4>
      <UserList />

      <Button className="mt-4 w-full bg-blue-500 hover:bg-blue-600 text-white">
        グループ作成
      </Button>
    </aside>
  )
}
