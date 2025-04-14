// src/components/Sidebar.tsx
import { useState } from "react"
import RoomList from "./RoomList"
import UserList from "./UserList"
import { Button } from "@/components/ui/button"

type SidebarProps = {
  onSelectRoom: (id: string, name: string) => void
}

export default function Sidebar({ onSelectRoom }: SidebarProps) {

  const [selectedUserIds, setSelectedUserIds] = useState<number[]>([])
  const [refreshRoomList, setRefreshRoomList] = useState(false)

  const createOneOnOne = async (userId: number, userName: string) => {
    const res = await fetch("http://localhost:8081/rooms", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        user_ids: [userId], // 個人チャット
        display_name: ""
      })
    })
  
    const data = await res.json()
    if (res.ok && data.room_id) {
      setRefreshRoomList((prev) => !prev)
      onSelectRoom(data.room_id, userName)
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
      // 成功時の処理
      // ✅ ルーム作成後にRoomList再読み込み
      setRefreshRoomList((prev) => !prev)

      // ✅ 作成されたルームに即入る
      onSelectRoom(data.room_id, data.display_name || "新しいグループ")
    } else {
      alert(data.error || "ルーム作成に失敗しました")
    }
  }  

  return (
    <aside className="w-full bg-gray-100 p-4 overflow-y-auto">
      <h3 className="text-lg font-bold mb-2">チャット一覧</h3>
      <RoomList onSelectRoom={onSelectRoom} refreshTrigger={refreshRoomList} />

      <h4 className="text-md font-semibold mt-6 mb-2">ユーザー一覧（新規チャット）</h4>
      <UserList
        selectedUserIds={selectedUserIds}
        setSelectedUserIds={setSelectedUserIds}
        onCreateOneOnOne={createOneOnOne}
      />

      <Button
        className="mt-4 w-full bg-blue-500 hover:bg-blue-600 text-white" onClick={createGroup}>
        グループ作成
      </Button>
    </aside>
  )
}
