# 🌸 I LOVE SIC BO 🌸
(Sòng Bài Tàii Xỉuuu Mẹ Yêuuu - ĐÁNH TÀI RA TÀI, ĐÁNH XỈU RA XỈU)

```text
  ___   _       ___ __     __ _____   ____ ___  ____     ____   ___  
 |_ _| | |     / _ \\ \   / /| ____| / ___|_ _|/ ___|   | __ ) / _ \ 
  | |  | |    | | | |\ \ / / |  _|   \___ \| | | |       |  _ \| | | |
  | |  | |___ | |_| | \ V /  | |___   ___) | | | |___    | |_) | |_| |
 |___| |_____| \___/   \_/   |_____| |____/___| \____|   |____/ \___/ 
```

Chào mừng bạn đến với **I LOVE SIC BO** — Sòng bài Sic Bo (Tài Xỉu) phong cách TUI (Terminal User Interface) cực kỳ "premium" với tông màu hồng cánh sen chủ đạo! 🎀

Đây là một dự án Go hoàn chỉnh cho phép bạn host một sòng bài trong mạng LAN và rủ bạn bè cùng tham gia sát phạt bằng giao diện ASCII Art cực kỳ dễ thương.

---

## ✨ Tính năng nổi bật

- 🎲 **Giao diện TUI Pink:** Xây dựng bằng `Bubble Tea` & `Lipgloss`, màu hồng lung linh, font chữ ASCII nghệ thuật.
- 📡 **Tự động tìm sòng:** Sử dụng **mDNS** giúp các máy con tự tìm thấy máy chủ trong mạng LAN mà không cần nhập IP.
- 🎰 **8 Loại cược:** Tài/Xỉu, Chẵn/Lẻ, Bộ ba, Đôi, Cặp, Tổng chính xác... đầy đủ như sòng chuyên nghiệp.
- 📈 **Bảng vị (History):** Theo dõi lịch sử ván đấu bằng các icon 🌸 (Tài) và 🍀 (Xỉu).
- 💾 **Persistence:** Lưu trữ số dư và lịch sử giao dịch vào SQLite (`casino.db`).
- ⚡ **Animation:** Hiệu ứng xúc xắc xoay xoay đẹp mắt trước khi công bố kết quả.

---

## 🛠 Yêu cầu hệ thống

- **Ngôn ngữ:** Go version 1.25+ (hoặc mới nhất).
- **Hệ điều hành:** macOS, Linux hoặc Windows (khuyến khích dùng Terminal hỗ trợ 256 màu/TrueColor).

---

## 🚀 Hướng dẫn cài đặt & Chạy thử

### 1. Tải mã nguồn & Cài đặt thư viện
```bash
git clone https://github.com/panadolextra91/i-love-sicbo.git
cd i-love-sicbo
go mod tidy
```

### 2. Chạy máy cái (Server)
Người host sòng bài cần chạy lệnh này:
```bash
go run cmd/server/main.go
```
*Ghi chú: Lần đầu chạy sẽ tự tạo file `casino.db` để lưu dữ liệu.*

### 3. Tham gia chơi (Client)
Tất cả người chơi (kể cả chủ sòng) chạy lệnh này:
```bash
go run cmd/client/main.go
```
*Client sẽ tự động tìm thấy Server và kết nối. Nếu mDNS lỗi, nó sẽ tự fallback về localhost.*

---

## 🎮 Cách chơi & Phím tắt

Hệ thống điều khiển theo phong cách "Hacker" bằng phím tắt:

- **[1]**: Cược **Tài/Xỉu** (Nhấn nhiều lần để đổi cửa).
- **[2]**: Cược **Chẵn/Lẻ**.
- **[3]**: **Bộ ba bất kỳ**.
- **[4]**: **Bộ ba chính xác** (111 - 666).
- **[5]**: **Đôi chính xác**.
- **[6]**: **Cặp số khác nhau**.
- **[7]**: **Tổng chính xác** (4 - 17).
- **[8]**: **Một số cụ thể** (Đoán số lần xuất hiện 1-3).
- **[Tab]**: Chuyển đổi giá trị chọn (Ví dụ: Đổi từ Tài sang Xỉu, hoặc đổi số tổng).
- **[+] / [-]**: Tăng hoặc giảm tiền cược (Bước nhảy 100).
- **[Enter]**: **CHỐT HẠ** đặt cược.
- **[Q] / [Ctrl+C]**: Thoát sòng.

---

## 🏗 Công nghệ sử dụng

- **Core:** Go (Golang)
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Networking:** Gorilla WebSocket & Hashicorp mDNS
- **Database:** SQLite3
- **Design:** ASCII Art & Pink Aesthetic Styling

---

## 📜 Luật chơi Sic Bo (Sơ lược)

1. **Tài (Big):** Tổng điểm 11-17. Thua nếu ra bộ ba (`1-1-1`, `2-2-2`...).
2. **Xỉu (Small):** Tổng điểm 4-10. Thua nếu ra bộ ba.
3. **Chẵn/Lẻ:** Tổng điểm chẵn hoặc lẻ. Thua nếu ra bộ ba.
4. **Bộ ba chính xác:** Thắng đậm nhất với tỷ lệ `180:1`.

---

Chúc mẹ và anh em đồng đạo có những giây phút giải trí cực hồng và cực vui tại **I LOVE SIC BO**! 🐹💖🌸
