// src/components/ChatArea.tsx
import { useEffect, useRef, useState } from "react"
import { Pencil } from "lucide-react"
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog"

type Message = {
  id: string
  sender_id: number
  sender: string
  content: string
  created_at: string
}

type ChatAreaProps = {
  roomId: string
  roomName: string
  userId: number
  isGroup: boolean
}

export default function ChatArea({ roomId, roomName, userId, isGroup }: ChatAreaProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState("")
  const socketRef = useRef<WebSocket | null>(null)
  const chatLogRef = useRef<HTMLUListElement>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [isEditingName, setIsEditingName] = useState(false)
  const [newRoomName, setNewRoomName] = useState(roomName)
  const [currentRoomName, setCurrentRoomName] = useState(roomName)
  const [members, setMembers] = useState<string[]>([])
  const lastMessageRef = useRef<HTMLLIElement | null>(null)

  // グループ名の最大文字数
  const maxGroupNameLength = 30

  // HTTPベースのAPI URL（例: https://backend.fly.dev）
  const httpApiUrl = import.meta.env.VITE_API_URL
  // WebSocketのURLに変換（https → wss, http → ws）
  const wsUrl = httpApiUrl.replace(/^http/, "ws")

  // メッセージ取得 + WebSocket接続
  useEffect(() => {
    if (!roomId) return

    // 過去ログ取得
    const token = localStorage.getItem("jwt_token")
    if (!token) return
    fetch(`${import.meta.env.VITE_API_URL}/messages/${roomId}?limit=30`, {
      headers: {
        "Authorization": `Bearer ${token}`,
      },
    })
      .then((res) => res.json())
      .then((data) => {
        setMessages(data || [])
        setHasMore(data.length === 30)
        markAsRead(roomId)
    })

    // WebSocket接続
    socketRef.current?.close()
    const ws = new WebSocket(`${wsUrl}/ws?room=${roomId}&token=${token}`)
    socketRef.current = ws

    // Websocket接続判定
    ws.onopen = () => {
      setIsConnected(true)
      console.log("WebSocket接続成功")
    }
    ws.onclose = () => {
      setIsConnected(false)
      console.log("WebSocket切断")
    }

    ws.onmessage = (event) => {
      const msg: Message & { from_self?: boolean } = JSON.parse(event.data)


      // from_self でなければ JST 補正
      if (!msg.from_self) {
        const date = new Date(msg.created_at)
        date.setHours(date.getHours() + 9)
        msg.created_at = date.toISOString()
      }

      setMessages((prev) => {
        const updated = [...prev, msg]

        // ✅ 最新のメッセージを受信後に既読更新
        fetch(`${import.meta.env.VITE_API_URL}/rooms/${roomId}/read`, {
          method: "POST",
          headers: { Authorization: `Bearer ${token}` },
        })

        return updated
      })
      scrollToBottom()
    }

    return () =>{
        ws.close()
        setIsConnected(false)
    }
  }, [roomId])

  // メッセージ送信
  const notifySocketRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    if (!roomId || !token) return
    const notifyWS = new WebSocket(`${wsUrl}/ws-notify?token=${token}`)
    notifySocketRef.current = notifyWS
    return () => notifyWS.close()
  }, [roomId])

  // ルーム名を編集
  const handleRoomNameUpdate = async () => {

    if (!newRoomName.trim()) {
      alert("グループ名を入力してください")
      return
    }
    if (newRoomName.length > 30) {
      alert("グループ名は30文字以内で入力してください")
      return
    }

    const token = localStorage.getItem("jwt_token")
    const res = await fetch(`${import.meta.env.VITE_API_URL}/rooms/${roomId}/name`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ display_name: newRoomName }),
    })

    if (res.ok) {
      setCurrentRoomName(newRoomName) // ← ここが重要
      setIsEditingName(false)
    } else {
      alert("グループ名の更新に失敗しました")
    }
  }
  useEffect(() => {
    setCurrentRoomName(roomName)
    setNewRoomName(roomName)
  }, [roomName])

  // グループメンバー取得
  const fetchMembers = async () => {
    const token = localStorage.getItem("jwt_token")
    const res = await fetch(`${import.meta.env.VITE_API_URL}/rooms/${roomId}/members`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      const data = await res.json()
      setMembers(data.map((u: { name: string }) => u.name))
    }
  }
  
  // グループから退会
  const handleLeaveGroup = async () => {
    if (!window.confirm("本当にグループを退会しますか？")) return
  
    const token = localStorage.getItem("jwt_token")
    const res = await fetch(`${import.meta.env.VITE_API_URL}/rooms/${roomId}/members/me`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    })
  
    if (res.ok) {
      window.location.reload() // ✅ 退会後にリスト再読込
    } else {
      alert("退会に失敗しました")
    }
  }
  
  
  // グループを削除
  const handleDeleteGroup = async () => {
    if (!window.confirm("このグループを完全に削除しますか？")) return
  
    const token = localStorage.getItem("jwt_token")
    const res = await fetch(`${import.meta.env.VITE_API_URL}/rooms/${roomId}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    })
  
    if (res.ok) {
      window.location.href = "/chat" // ✅ 削除後にトップに戻るなど
    } else {
      alert("削除に失敗しました")
    }
  }  

  // メッセージ送信
  const handleSend = () => {
    if (!input.trim()) return
  
    if (input.length > 1000) {
      alert("メッセージは1000文字以内で入力してください")
      return
    }
  
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      socketRef.current.send(input)
      setInput("")
  
      notifySocketRef.current?.send(
        JSON.stringify({
          room_id: roomId,
          sender_id: userId,
          created_at: new Date().toISOString(),
          from_self: true,
        })
      )
    }
  }

  // 古いメッセージ読み込み
  const [isLoading, setIsLoading] = useState(false)
  const [hasMore, setHasMore] = useState(true)

  const loadOlderMessages = async () => {
    if (!roomId || isLoading || !hasMore || messages.length === 0) return
  
    setIsLoading(true)
  
    const token = localStorage.getItem("jwt_token")
    const oldest = messages[0].created_at
    const res = await fetch(`${import.meta.env.VITE_API_URL}/messages/${roomId}?before=${oldest}&limit=30`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    const data = await res.json()
  
    const container = chatLogRef.current
    const prevHeight = container?.scrollHeight ?? 0
  
    if (data.length === 0) {
      setHasMore(false)
  
      // 最上部にピタッと止める
      setTimeout(() => {
        if (container) container.scrollTop = 0
      }, 10)
  
      setIsLoading(false)
      return
    }
  
    setMessages((prev) => [...data, ...prev])
  
    setTimeout(() => {
      if (container) {
        const newHeight = container.scrollHeight
        container.scrollTop = newHeight - prevHeight
      }
    }, 10)
  
    setIsLoading(false)
    markAsRead(roomId)
  }

  const markAsRead = (roomId: string) => {
    const token = localStorage.getItem("jwt_token")
    if (!token) return
    fetch(`${httpApiUrl}/rooms/${roomId}/read`, {
      method: "POST",
      headers: { Authorization: `Bearer ${token}` },
    })
  }

  useEffect(() => {
    const container = chatLogRef.current
    if (!container) return
  
    const handleScroll = () => {
      if (container.scrollTop < 50 && !isLoading && hasMore) {
        loadOlderMessages()
      }
    }
  
    container.addEventListener("scroll", handleScroll)
    return () => container.removeEventListener("scroll", handleScroll)
  }, [isLoading, hasMore, messages])
  


  // スクロール位置を最下部に
  useEffect(() => {
    if (lastMessageRef.current) {
      lastMessageRef.current.scrollIntoView({ behavior: "auto" })
    }
  }, [messages])
  const scrollToBottom = () => {
    setTimeout(() => {
      chatLogRef.current?.scrollTo({
        top: chatLogRef.current.scrollHeight,
        behavior: "smooth",
      })
    }, 10)
  }

  useEffect(() => {
    const token = localStorage.getItem("jwt_token")
    if (!roomId || !token) return
  
    fetch(`${import.meta.env.VITE_API_URL}/rooms/${roomId}/read`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
  }, [roomId])  

  const formatDate = (dateStr: string) =>
    new Date(dateStr).toISOString().slice(0, 10)

  const formatTime = (dateStr: string) =>
    new Date(dateStr).toISOString().slice(11, 16)

  let lastRenderedDate: string | null = null

  return (
    <main className="flex-1 min-w-[600px] flex flex-col h-screen p-4">
      <div className="flex items-center justify-between mb-2"> {/* justify-betweenで左右に配置 */}
      <div>
        {isEditingName ? (
          <>
            <div className="flex items-center gap-2">
              <input
                className="border px-2 py-1 rounded text-sm w-80"
                value={newRoomName}
                onChange={(e) => setNewRoomName(e.target.value)}
                maxLength={maxGroupNameLength + 10}
              />
              <button
                onClick={handleRoomNameUpdate}
                className=""
              >
                保存
              </button>
              <button
                onClick={() => setIsEditingName(false)}
                className=""
              >
                キャンセル
              </button>
            </div>

            {/* 文字数カウント（下に配置） */}
            <div
              className={`text-xs mt-1 ${newRoomName.length > maxGroupNameLength ? "text-red-500" : "text-gray-400"}`}
            >
              {newRoomName.length}/{maxGroupNameLength}文字
            </div>
          </>        
        ) : (
          <div className="flex items-center gap-2"> {/* グループ名と編集ボタンをまとめる */}
            <h3 className="text-xl font-bold">
              <span
              className="truncate block max-w-[400px]"
              title={currentRoomName} // ✅ ホバーで全文表示
              >
                {currentRoomName}
              </span>
            </h3>
            {isGroup && (
              <>
              <button
                className="rounded-md p-1"
                onClick={() => {
                  setNewRoomName(currentRoomName)
                  setIsEditingName(true)
                }}
              >
                <Pencil size={18} />
              </button>
        
              {/* ✅ ここに DialogTrigger を追加 */}
              <Dialog onOpenChange={(open) => open && fetchMembers()}>
                <DialogTrigger asChild>
                  <button>メンバー一覧</button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>メンバー一覧</DialogTitle>
                  </DialogHeader>
        
                  {/* メンバー一覧 */}
                  <ul className="text-sm list-disc list-inside space-y-1 my-2">
                    {members.map((name, i) => (
                      <li key={i}>{name}</li>
                    ))}
                  </ul>
        
                  {/* 区切り + 操作ボタン */}
                  <div className="border-t pt-4 mt-4 flex gap-2">
                    <button onClick={handleLeaveGroup}>
                      グループを退会する
                    </button>
                    <button onClick={handleDeleteGroup}>
                      グループを削除する
                    </button>
                  </div>
                </DialogContent>
              </Dialog>
            </>
            )}
          </div>
        )}
      </div>
      <span className={`text-sm ${isConnected ? "text-green-600" : "text-red-500"}`}>
        {isConnected ? "● 接続中" : "● 切断中"}
      </span>
    </div>

    <ul ref={chatLogRef} className="flex-1 flex flex-col overflow-y-auto border rounded p-2 space-y-1 bg-white">
    {messages.map((msg, index) => {
      const currentDate = formatDate(msg.created_at)
      const showDate = currentDate !== lastRenderedDate
      lastRenderedDate = currentDate

      const isLast = index === messages.length - 1 // ← 最後の要素か判定

      return (
        <div key={msg.id}>
          {showDate && (
            <li className="text-xs text-gray-500 text-center py-1">
              <div className="bg-gray-100 py-0.5 rounded-full inline-block px-2">
                --- {currentDate} ---
              </div>
            </li>
          )}
          <li
            className={`flex ${msg.sender_id === userId ? "justify-end" : "justify-start"}`}
            ref={isLast ? lastMessageRef : undefined} // ✅ 最後のメッセージに ref をつける
          >
            <div
              className={`text-sm p-2 rounded max-w-[70%] break-words whitespace-pre-wrap ${
                msg.sender_id === userId ? "bg-blue-200 text-right" : "bg-gray-100 text-left"
              }`}
            >
              <span>{msg.content}</span>
              <div className="text-xs text-gray-500 block mt-1">
                [{formatTime(msg.created_at)}] {msg.sender}
              </div>
            </div>
          </li>
        </div>
      )
    })}
    </ul>

    <div className="mt-2">
      <div className="flex gap-2 items-end">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault()
              handleSend()
            }
          }}
          rows={2}
          className="flex-1 border border-gray-300 rounded p-2 resize-none"
          placeholder="メッセージを入力（Shift+Enterで改行）"
        />
        <button onClick={handleSend}>送信</button>
      </div>
      
      {/* 文字数カウント（右下） */}
      <div className={`text-right text-xs mt-1 pr-1 ${
        input.length > 1000 ? "text-red-500" : "text-gray-400"
      }`}>
        {input.length}/1000文字
      </div>
    </div>
    </main>
  )
}
