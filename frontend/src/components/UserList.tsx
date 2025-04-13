// src/components/UserList.tsx
import { useEffect, useState } from "react"

type User = {
  id: number
  name: string
}

export default function UserList() {
  const [users, setUsers] = useState<User[]>([])

  useEffect(() => {
    fetch("http://localhost:8081/users", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => setUsers(data || []))
  }, [])

  return (
    <ul className="space-y-2">
      {users.map((user) => (
        <li key={user.id} className="flex justify-between items-center">
          <label className="flex items-center gap-2">
            <input type="checkbox" value={user.id} className="userCheckbox" />
            <span>{user.name}</span>
          </label>
          <button
            onClick={() => {/* createOneOnOne(user.id) */}}
            className="text-sm px-2 py-1 bg-blue-400 text-white rounded hover:bg-blue-500"
          >
            個人チャット
          </button>
        </li>
      ))}
    </ul>
  )
}
