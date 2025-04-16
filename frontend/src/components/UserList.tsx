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
    fetch("http://localhost:8081/users", { credentials: "include" })
      .then((res) => res.json())
      .then((data) => setUsers(data || []))
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
            <span>{user.name}</span>
          </label>
          <button
            onClick={() => onCreateOneOnOne(user.id, user.name)}
            className="text-sm px-2 py-1 rounded"
          >
            個人チャット
          </button>
        </li>
      ))}
    </ul>
  )
}
