// src/components/ChatArea.tsx
import { useEffect, useRef, useState } from "react"
import { Pencil } from "lucide-react"

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
        fetch(`http://localhost:8081/rooms/${roomId}/read`, {
          method: "POST",
          credentials: "include",
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
    const notifyWS = new WebSocket("ws://localhost:8081/ws-notify")
    notifySocketRef.current = notifyWS
    return () => notifyWS.close()
  }, [])

  // ルーム名を編集
  const handleRoomNameUpdate = async () => {

    if (!newRoomName.trim()) {
      alert("グループ名を入力してください")
      return
    }
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
          created_at: new Date().toISOString(),
          from_self: true
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
    <main className="flex-1 min-w-[600px] flex flex-col h-screen p-4">
      <div className="flex items-center justify-between mb-2"> {/* justify-betweenで左右に配置 */}
      <div>
        {isEditingName ? (
          <>
            <input
              className="border px-2 py-1 rounded text-sm"
              value={newRoomName}
              onChange={(e) => setNewRoomName(e.target.value)}
            />
            <button onClick={handleRoomNameUpdate}> {/* 保存ボタン */}
              保存
            </button>
            <button onClick={() => setIsEditingName(false)}>
              キャンセル
            </button>
          </>
        ) : (
          <div className="flex items-center gap-2"> {/* グループ名と編集ボタンをまとめる */}
            <h3 className="text-xl font-bold">
              {isGroup ? `グループ: ${currentRoomName}` : `チャット相手: ${currentRoomName}`}
            </h3>
            {isGroup && (
              <button
                className="rounded-md p-1"
                onClick={() => {
                  setNewRoomName(currentRoomName)
                  setIsEditingName(true)
                }}
                >
                <Pencil size={18} />
              </button>
            )}
          </div>
        )}
      </div>
      <span className={`text-sm ${isConnected ? "text-green-600" : "text-red-500"}`}>
        {isConnected ? "● 接続中" : "● 切断中"}
      </span>
    </div>

    <ul ref={chatLogRef} className="flex-1 flex flex-col overflow-y-auto border rounded p-2 space-y-1 bg-white">
      {messages.map((msg) => {
        const currentDate = formatDate(msg.created_at);
        const showDate = currentDate !== lastRenderedDate;
        lastRenderedDate = currentDate;

        return (
          <div key={msg.id}>
            {showDate && (
              <li className="text-xs text-gray-500 text-center py-1"> {/* border-bを削除 */}
                <div className="bg-gray-100 py-0.5 rounded-full inline-block px-2"> {/* 背景色とpaddingを追加 */}
                  --- {currentDate} ---
                </div>
              </li>
            )}
            <li className={`flex ${msg.sender_id === userId ? "justify-end" : "justify-start"}`}>
              <div
                className={`text-sm p-2 rounded max-w-[70%] break-words whitespace-pre-wrap ${
                  msg.sender_id === userId
                    ? "bg-blue-200 text-right" // 少し明るい青
                    : "bg-gray-100 text-left"
                }`}
              >
                <span>{msg.content}</span> {/* メッセージ本文を先に */}
                <div className="text-xs text-gray-500 block mt-1"> {/* 下に移動、margin-topを追加 */}
                  [{formatTime(msg.created_at)}] {msg.sender}
                </div>
              </div>
            </li>
          </div>
        );
      })}
    </ul>

      <div className="mt-2 flex gap-2">
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
    </main>
  )
}
