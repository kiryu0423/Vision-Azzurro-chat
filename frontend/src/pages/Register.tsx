// src/pages/Register.tsx
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useState } from "react"

export default function Register() {
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirm, setConfirm] = useState("")
  const [error, setError] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (password !== confirm) {
      setError("パスワードが一致しません")
      return
    }

    const res = await fetch("http://localhost:8081/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name, email, password }),
    })

    if (res.ok) {
      window.location.href = "/" // ログインページに戻す
    } else {
      const data = await res.json().catch(() => ({}))
      setError(data?.error || "登録に失敗しました")
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <form
        onSubmit={handleSubmit}
        className="bg-white shadow-md rounded px-8 pt-6 pb-8 w-96 space-y-4"
      >
        <h1 className="text-2xl font-bold text-center">新規登録</h1>

        <div>
          <label className="block text-sm font-medium mb-1">名前</label>
            <Input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="名前を入力"
                required
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">メールアドレス</label>
          <Input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="メールアドレスを入力"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">パスワード</label>
          <Input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="8文字以上"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">パスワード確認</label>
          <Input
            type="password"
            value={confirm}
            onChange={(e) => setConfirm(e.target.value)}
            placeholder="再度入力"
            required
          />
        </div>

        {error && <p className="text-red-500 text-sm text-center">{error}</p>}

        <Button type="submit" className="w-full">
          登録する
        </Button>
      </form>
    </div>
  )
}
