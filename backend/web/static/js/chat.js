"use strict";

let currentRoomID = null;
let socket = null;
let roomState = []; // ← 各ルーム { room_id, display_name, lastMessageAt } を持つ
let lastRenderedDateStr = null;

// ② ユーザー一覧取得
async function loadUsers() {
    const res = await fetch("/users");
    const users = await res.json();

    const userList = document.getElementById("userList");
    userList.innerHTML = "";

    users.forEach(user => {
    const li = document.createElement("li");

    // ✅ チェックボックス＋個人チャットボタンの併用
    li.innerHTML = `
    <div class="flex items-center justify-between">
        <label class="flex items-center gap-2">
        <input type="checkbox" value="${user.id}" class="userCheckbox">
        <span>${user.name}</span>
        </label>
        <button onclick="createOneOnOne(${user.id}, '${user.name}')"
                class="text-sm px-2 py-1 bg-blue-400 text-white rounded hover:bg-blue-500">
        個人チャット
        </button>
    </div>
    `;


    userList.appendChild(li);
    });
}

// 個人チャット作成
async function createOneOnOne(userID, userName) {
    const res = await fetch("/rooms", {
    method: "POST",
    headers: {
        "Content-Type": "application/json"
    },
    credentials: "include",
    body: JSON.stringify({ user_ids: [userID] })
    });

    const data = await res.json();
    if (res.ok && data.room_id) {
        await refreshRooms();
        joinRoom(data.room_id, userName);
    } else {
    alert(data.error || "チャット作成に失敗しました");
    }
}

// ③ グループ作成
async function createGroup() {
    const checkboxes = document.querySelectorAll(".userCheckbox:checked");
    const userIDs = Array.from(checkboxes).map(cb => parseInt(cb.value));

    if (userIDs.length < 2) {
      alert("2人以上選択してください");
      return;
    }

    const res = await fetch("/rooms", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
      body: JSON.stringify({
        user_ids: userIDs,
        display_name: ""
      })
    });

    const data = await res.json();

    if (res.ok && data.room_id) {
      await refreshRooms(); // ✅ 最新のルーム一覧を再取得
      const roomName = data.display_name || "新しいグループ";
      joinRoom(data.room_id, roomName); // ✅ フォールバック付きで表示
    } else {
      alert(data.error || "グループ作成に失敗しました");
    }
}

// ② ルームを切り替える
async function joinRoom(roomID, roomName, element) {
    lastRenderedDateStr = null;
    currentRoomID = roomID;
    document.getElementById("roomTitle").textContent = "チャット相手: " + roomName;

    // 過去ログ取得
    const res = await fetch(`/messages/${roomID}`);
    const messages = await res.json();

    const chatLog = document.getElementById("chatLog");
    chatLog.innerHTML = "";

    // ✅ 日時付き表示
    messages.forEach(msg => {
    const messageTime = new Date(msg.created_at);

    // 日付部分だけ切り出し（例: "2025-04-13"）
    const currentDateStr = messageTime.toISOString().slice(0, 10);

    if (currentDateStr !== lastRenderedDateStr) {
        // 日付見出しを追加
        const dateLi = document.createElement("li");
        dateLi.textContent = `--- ${currentDateStr} ---`;
        dateLi.classList.add("text-xs", "text-gray-500", "text-center", "py-1", "border-b");
        chatLog.appendChild(dateLi);
        lastRenderedDateStr = currentDateStr;
    }

    // 時刻（UTCのまま or JST補正不要）
    const hhmm = messageTime.toISOString().substring(11, 16);

    const li = document.createElement("li");
    li.textContent = `[${hhmm}] ${msg.sender}: ${msg.content}`;
    chatLog.appendChild(li);
    });

    // WebSocket接続しなおし
    if (socket) socket.close();
    socket = new WebSocket("ws://localhost:8081/ws?room=" + roomID);

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        const chatLog = document.getElementById("chatLog");

        const msgTime = new Date(data.created_at);
        const jst = new Date(msgTime.getTime() + (9 * 60 * 60 * 1000));
        const currentDateStr = jst.toISOString().slice(0, 10); // "YYYY-MM-DD"

        // ✅ 日付ラベルは直前のメッセージと異なる日なら表示
        if (currentDateStr !== lastRenderedDateStr) {
            const dateLi = document.createElement("li");
            dateLi.textContent = `--- ${currentDateStr} ---`;
            dateLi.classList.add("text-xs", "text-gray-500", "text-center", "py-1", "border-b");
            chatLog.appendChild(dateLi);
            lastRenderedDateStr = currentDateStr;
        }

        const hhmm = jst.toISOString().substring(11, 16);

        const li = document.createElement("li");
        li.textContent = `[${hhmm}] ${data.sender}: ${data.content}`;
        chatLog.appendChild(li);

        // 並び替え用にルーム更新
        const room = roomState.find(r => r.room_id === currentRoomID);
        if (room) {
            room.last_message_at = data.created_at;
            renderRoomList();
        }

        scrollToBottom();
    };

    document.querySelectorAll("#roomList li").forEach(li => li.classList.remove("selected"));
    element?.classList?.add("selected");
}

// ③ メッセージ送信
function sendMessage() {
    const input = document.getElementById("messageInput");
    const msg = input.value;
    if (socket && msg.trim()) {
    socket.send(msg);
    input.value = "";
    }
}

// ルーム一覧更新
function renderRoomList() {
    const roomList = document.getElementById("roomList");
    roomList.innerHTML = "";

    // lastMessageAt で並べ替え（新着が上）
    roomState.sort((a, b) => {
        const timeA = new Date(a.last_message_at || 0);  // ← null対策
        const timeB = new Date(b.last_message_at || 0);
        return timeB - timeA;
    });

    roomState.forEach(room => {
        const li = document.createElement("li");
        li.textContent = room.display_name;
        li.style.cursor = "pointer";
        li.onclick = (e) => joinRoom(room.room_id, room.display_name, e.currentTarget);
        roomList.appendChild(li);
    });
}

// ルーム更新
async function refreshRooms() {
    const res = await fetch("/rooms", { credentials: "include" });
    const data = await res.json();
    roomState = Array.isArray(data) ? data : [];
    renderRoomList();
}


function scrollToBottom() {
    const chatLog = document.getElementById("chatLog");
    chatLog.scrollTop = chatLog.scrollHeight;
}

// ④ ログアウト
function logout() {
    fetch("/logout", { method: "POST", credentials: "include" })
    .then(() => window.location.href = "/login-page");
}

// 初期化
window.onload = () => {
    refreshRooms();
    loadUsers();
};
