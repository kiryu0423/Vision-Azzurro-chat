// src/components/ChatArea.tsx
import { useEffect, useRef, useState } from "react"
import { Pencil } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

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
}

export default function ChatArea({ roomId, roomName, userId }: ChatAreaProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState("")
  const socketRef = useRef<WebSocket | null>(null)
  const chatLogRef = useRef<HTMLUListElement>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [isEditingName, setIsEditingName] = useState(false)
  const [newRoomName, setNewRoomName] = useState(roomName)
  const [currentRoomName, setCurrentRoomName] = useState(roomName)

  // 日付ラベルの表示用
  const lastDateRef = useRef<string | null>(null)

  // メッセージ取得 + WebSocket接続
  useEffect(() => {
    if (!roomId) return

    // 過去ログ取得
    fetch(`http://localhost:8081/messages/${roomId}`, { credentials: "include" })
      .then((res) => res.json())
      .then((data) => {
        setMessages(data || [])
        scrollToBottom()
        lastDateRef.current = null
      })

    // ルームの既読更新
    fetch(`http://localhost:8081/rooms/${roomId}/read`, {
      method: "POST",
      credentials: "include",
    })

    // WebSocket接続
    socketRef.current?.close()
    const ws = new WebSocket(`ws://localhost:8081/ws?room=${roomId}`)
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
      const msg: Message = JSON.parse(event.data)
      setMessages((prev) => [...prev, msg])
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
    const notifyWS = new WebSocket("ws://localhost:8081/ws-notify")
    notifySocketRef.current = notifyWS
    return () => notifyWS.close()
  }, [])

  // ルーム名を編集
  const handleRoomNameUpdate = async () => {
    const res = await fetch(`http://localhost:8081/rooms/${roomId}/name`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
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

  const handleSend = () => {
    if (socketRef.current?.readyState === WebSocket.OPEN && input.trim()) {
      socketRef.current.send(input)
      setInput("")

      notifySocketRef.current?.send(
        JSON.stringify({
          room_id: roomId,
          sender_id: userId,
          created_at: new Date().toLocaleString(),
        })
      )
    }
  }

  const scrollToBottom = () => {
    setTimeout(() => {
      chatLogRef.current?.scrollTo({
        top: chatLogRef.current.scrollHeight,
        behavior: "smooth",
      })
    }, 100)
  }

  const formatDate = (dateStr: string) =>
    new Date(dateStr).toISOString().slice(0, 10)

  const formatTime = (dateStr: string) =>
    new Date(dateStr).toISOString().slice(11, 16)

  let lastRenderedDate: string | null = null

  return (
    <main className="flex-1 w-full flex flex-col h-screen p-4">
      <div className="flex items-center gap-2 mb-2">
      {isEditingName ? (
          <>
            <input
              className="border px-2 py-1 rounded text-sm"
              value={newRoomName}
              onChange={(e) => setNewRoomName(e.target.value)}
            />
            <Button size="sm" onClick={handleRoomNameUpdate}>
              保存
            </Button>
            <Button variant="outline" size="sm" onClick={() => setIsEditingName(false)}>
              キャンセル
            </Button>
          </>
      ) : (
          <>
            <h3 className="text-xl font-bold">チャット相手: {currentRoomName}</h3>
            <button
              className="text-gray-500 hover:text-gray-700"
              onClick={() => {
                setNewRoomName(currentRoomName)
                setIsEditingName(true)
              }}
            >
              <Pencil size={18} />
            </button>
          </>
      )}
        <span className={`text-sm ${isConnected ? "text-green-600" : "text-red-500"}`}>
          {isConnected ? "● 接続中" : "● 切断中"}
        </span>
      </div>
  
      <ul ref={chatLogRef} className="flex-1 flex flex-col overflow-y-auto border rounded p-2 space-y-1 bg-white">
        {messages.map((msg) => {
          const currentDate = formatDate(msg.created_at)
          const showDate = currentDate !== lastRenderedDate
          lastRenderedDate = currentDate
  
          return (
            <div key={msg.id}>
              {showDate && (
                <li className="text-xs text-gray-500 text-center py-1 border-b">
                  --- {currentDate} ---
                </li>
              )}
              <li
                className={`text-sm p-2 rounded max-w-[70%] break-words ${
                  msg.sender_id === userId
                    ? "bg-blue-100 self-end text-right"
                    : "bg-gray-100 self-start text-left"
                }`}
              >
                <span className="text-xs text-gray-500 block">
                  [{formatTime(msg.created_at)}] {msg.sender}
                </span>
                <span>{msg.content}</span>
              </li>
            </div>
          )
        })}
      </ul>
  
      <div className="mt-2 flex gap-2">
        <Input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="メッセージを入力..."
          onKeyDown={(e) => e.key === "Enter" && handleSend()}
        />
        <Button onClick={handleSend}>送信</Button>
      </div>
    </main>
  )
}
