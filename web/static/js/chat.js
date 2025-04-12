"use strict";

let currentRoomID = null;
let socket = null;

// ① ルーム一覧取得
async function loadRooms() {
    const res = await fetch("/rooms");
    const data = await res.json();

    // nullチェック＋配列確認（← これでエラー防止）
    const rooms = Array.isArray(data) ? data : [];

    const roomList = document.getElementById("roomList");
    roomList.innerHTML = "";

    rooms.forEach(room => {
    const li = document.createElement("li");
    li.textContent = room.display_name || "未設定ルーム名";
    li.style.cursor = "pointer";
    li.onclick = () => {
        joinRoom(room.room_id, room.display_name);
    };
    roomList.appendChild(li);
    });
}

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
    joinRoom(data.room_id, data.display_name);
    } else {
    alert(data.error || "グループ作成に失敗しました");
    }
}


// ② ルームを切り替える
async function joinRoom(roomID, roomName) {
    currentRoomID = roomID;
    document.getElementById("roomTitle").textContent = "チャット相手: " + roomName;

    // 過去ログ取得
    const res = await fetch(`/messages/${roomID}`);
    const messages = await res.json();

    const chatLog = document.getElementById("chatLog");
    chatLog.innerHTML = "";

    // ✅ 日時付き表示
    messages.forEach(msg => {
    const li = document.createElement("li");

    const time = new Date(msg.created_at);
    const loghhmm = time.toLocaleTimeString('ja-JP', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
        timeZone: 'Asia/Tokyo'
    });

    li.textContent = `[${loghhmm}] ${msg.sender}: ${msg.content}`;
    chatLog.appendChild(li);
    });

    // WebSocket接続しなおし
    if (socket) socket.close();
    socket = new WebSocket("ws://localhost:8081/ws?room=" + roomID);

    socket.onmessage = (event) => {
    const chatLog = document.getElementById("chatLog");
    const li = document.createElement("li");

    // 日時送信
    const now = new Date();
    const hhmm = now.toLocaleTimeString('ja-JP', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false,
        timeZone: 'Asia/Tokyo'
    });
    li.textContent = `[${hhmm}] ${event.data}`;
    chatLog.appendChild(li);
    scrollToBottom();
    };


    document.querySelectorAll("#roomList li").forEach(li => li.classList.remove("selected"));
    event?.target?.classList?.add("selected")
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
    loadRooms();
    loadUsers();
};
