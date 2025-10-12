import { useEffect, useState } from "react";

// تابع ساخت رنگ رندوم برای هر کاربر
function randomColor() {
    const colors = ["#FF5733", "#33C3FF", "#75FF33", "#E833FF", "#FFC300", "#33FFD7"];
    return colors[Math.floor(Math.random() * colors.length)];
}

function App() {
    const [socket, setSocket] = useState(null);
    const [messages, setMessages] = useState([]);
    const [input, setInput] = useState("");
    const [username, setUsername] = useState("");
    const [color, setColor] = useState(randomColor());

    useEffect(() => {
        const ws = new WebSocket("ws://localhost:8080/ws");
        ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            setMessages((prev) => [...prev, msg]);
        };
        setSocket(ws);
        return () => ws.close();
    }, []);

    const sendMessage = () => {
        if (socket && input.trim() && username.trim()) {
            const msg = { user: username, text: input, color };
            socket.send(JSON.stringify(msg));
            setInput("");
        }
    };

    if (!username) {
        return (
            <div style={{ textAlign: "center", marginTop: "100px" }}>
                <h2>👋 خوش اومدی!</h2>
                <input
                    placeholder="یوزرنیم خودت رو بنویس..."
                    onKeyDown={(e) => e.key === "Enter" && e.target.value && setUsername(e.target.value)}
                    style={{
                        padding: "10px",
                        width: "250px",
                        borderRadius: "5px",
                        border: "1px solid #aaa",
                    }}
                />
                <p style={{ color: "#555", marginTop: "10px" }}>Enter بزن برای ورود</p>
            </div>
        );
    }

    return (
        <div style={{ maxWidth: 500, margin: "50px auto", textAlign: "center" }}>
            <h2>💬 چت Go + React</h2>
            <div
                style={{
                    border: "1px solid #ccc",
                    height: "350px",
                    overflowY: "auto",
                    padding: "10px",
                    marginBottom: "10px",
                    borderRadius: "10px",
                    background: "#f9f9f9",
                }}
            >
                {messages.map((msg, i) => (
                    <div key={i} style={{ textAlign: msg.user === username ? "right" : "left" }}>
                        <span style={{ color: msg.color, fontWeight: "bold" }}>{msg.user}: </span>
                        <span>{msg.text}</span>
                    </div>
                ))}
            </div>
            <div>
                <input
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && sendMessage()}
                    style={{
                        width: "70%",
                        padding: "10px",
                        borderRadius: "5px",
                        border: "1px solid #aaa",
                    }}
                    placeholder="پیام بنویس..."
                />
                <button
                    onClick={sendMessage}
                    style={{
                        marginLeft: "10px",
                        padding: "10px 20px",
                        border: "none",
                        backgroundColor: "#007bff",
                        color: "white",
                        borderRadius: "5px",
                        cursor: "pointer",
                    }}
                >
                    ارسال
                </button>
            </div>
        </div>
    );
}

export default App;
