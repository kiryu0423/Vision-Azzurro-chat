// src/pages/Login.tsx

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useState } from "react"

export default function Login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const res = await fetch("${import.meta.env.VITE_API_URL}/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    })

    if (res.ok) {
      const data = await res.json()

      // ✅ JWTを localStorage に保存
      localStorage.setItem("jwt_token", data.token)

      // ✅ チャット画面に遷移
      window.location.href = "/chat"
    } else {
      setError("ログインに失敗しました")
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <form
        onSubmit={handleSubmit}
        className="bg-white shadow-md rounded px-8 pt-6 pb-8 w-96 space-y-4"
      >
        <h1 className="text-2xl font-bold text-center">ログイン</h1>

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
            placeholder="パスワードを入力"
            required
          />
        </div>

        {error && <p className="text-red-500 text-sm text-center">{error}</p>}

        <Button type="submit" className="w-full">
          ログイン
        </Button>
      </form>
    </div>
  )
}
