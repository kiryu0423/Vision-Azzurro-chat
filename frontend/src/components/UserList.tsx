import { useEffect, useState } from "react"

type User = {
  id: number
  name: string
}

type UserListProps = {
  selectedUserIds: number[]
  setSelectedUserIds: (ids: number[]) => void
  onCreateOneOnOne: (userId: number, userName: string) => void
}

export default function UserList({
  selectedUserIds,
  setSelectedUserIds,
  onCreateOneOnOne,
}: UserListProps) {
  const [users, setUsers] = useState<User[]>([])

  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    if (!token) return // トークンが無い場合は何もしない（必要ならリダイレクト）

    fetch("${import.meta.env.VITE_API_URL}/users", {
      headers: {
        "Authorization": `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("unauthorized")
        return res.json()
      })
      .then((data) => setUsers(data || []))
      .catch((err) => {
        console.error("ユーザー取得エラー:", err)
        // 例: window.location.href = "/"
      })
  }, [])

  const handleCheck = (id: number, checked: boolean) => {
    if (checked) {
      setSelectedUserIds([...selectedUserIds, id])
    } else {
      setSelectedUserIds(selectedUserIds.filter((uid) => uid !== id))
    }
  }

  return (
    <ul className="space-y-2">
      {users.map((user) => (
        <li key={user.id} className="flex justify-between items-center">
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={selectedUserIds.includes(user.id)}
              onChange={(e) => handleCheck(user.id, e.target.checked)}
            />
            <span className="truncate block max-w-[120px]" title={user.name}>
              {user.name}
            </span>
          </label>
          <button
            onClick={() => onCreateOneOnOne(user.id, user.name)}
            className="text-xs px-2 py-0.5 rounded border border-gray-300 hover:bg-gray-100"
          >
            チャット
          </button>
        </li>
      ))}
    </ul>
  )
}
