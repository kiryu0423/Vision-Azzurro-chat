// src/components/ChatArea.tsx
import { useEffect, useRef, useState } from "react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

type Message = {
  id: string
  sender: string
  content: string
  created_at: string
}

type ChatAreaProps = {
  roomId: string
  roomName: string
}

export default function ChatArea({ roomId, roomName }: ChatAreaProps) {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState("")
  const socketRef = useRef<WebSocket | null>(null)
  const chatLogRef = useRef<HTMLUListElement>(null)
  const [isConnected, setIsConnected] = useState(false)

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
    }

    return () =>{
        ws.close()
        setIsConnected(false)
    }
  }, [roomId])

  // メッセージ送信
  const handleSend = () => {
    if (socketRef.current?.readyState === WebSocket.OPEN && input.trim()) {
      socketRef.current.send(input)
      setInput("")
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
    <main className="flex-1 w-full flex flex-col p-4">
      <div className="flex items-center gap-2 mb-2">
        <h3 className="text-xl font-bold">チャット相手: {roomName}</h3>
        <span className={`text-sm ${isConnected ? "text-green-600" : "text-red-500"}`}>
            {isConnected ? "● 接続中" : "● 切断中"}
        </span>
      </div>

      <ul ref={chatLogRef} className="flex-1 overflow-y-auto border rounded p-2 space-y-1 bg-white">
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
              <li className="text-sm">
                [{formatTime(msg.created_at)}] {msg.sender}: {msg.content}
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
