# ซ้อม Live Coding สัมภาษณ์งาน Golang Developer

ชุดฝึกที่อ้างอิงหัวข้อจาก [roadmap.sh/questions/golang](https://roadmap.sh/questions/golang)

| ไฟล์/โฟลเดอร์ | เนื้อหา |
|---|---|
| [PROBLEMS.md](PROBLEMS.md) | โจทย์ live coding 14 ข้อ (ภาษาไทย) |
| [SOLUTIONS.md](SOLUTIONS.md) | เฉลย + คำอธิบาย + คำตอบคำถามต่อยอด |
| [THEORY.md](THEORY.md) | คำถามทฤษฎี 50 ข้อ พร้อมคำตอบภาษาไทย |
| [golang-interview.html](golang-interview.html) | ทั้ง 3 ไฟล์ข้างบนรวมเป็นหน้าเว็บเดียว (เปิดใน browser ได้เลย มีตัวจับเวลา) |
| `exercises/` | โค้ดตั้งต้นที่ต้องเขียนต่อ + test ที่เตรียมไว้ |
| `solutions/` | โค้ดเฉลยที่ผ่าน test ทุกข้อ |

## เริ่มต้น

ต้องใช้ Go 1.24 ขึ้นไป

```bash
make test N=01        # รัน test ข้อ 01 (ตอนแรกจะ FAIL — เขียนโค้ดให้ผ่าน)
make solution N=01    # รัน test ของเฉลย
make reset N=01       # ล้างสิ่งที่เขียนไป กลับเป็นโจทย์เริ่มต้น
make test-all         # รันทุกข้อ
make html             # สร้าง golang-interview.html ใหม่หลังแก้ไฟล์ .md
```

## เทคนิคตอนสัมภาษณ์ live coding

1. **ถามให้ชัดก่อนเขียน:** input/output, edge case (ว่าง, nil, ติดลบ, Unicode), ขนาดข้อมูล, ต้อง thread-safe ไหม
2. **พูดความคิดออกมา:** interviewer ให้คะแนนวิธีคิดพอ ๆ กับโค้ด
3. **เริ่มจากวิธีง่ายที่ถูกต้อง** แล้วค่อย optimize พร้อมบอก Big-O
4. **จัดการ error ทุกจุด** และไม่ `_` ทิ้ง error เฉย ๆ
5. **เขียน test สั้น ๆ** หรืออย่างน้อยไล่ตัวอย่างด้วยมือ และถ้ามี goroutine ให้พูดถึง `-race`
6. **ปิดท้ายด้วยสิ่งที่จะทำต่อถ้ามีเวลา** เช่น context, timeout, logging, metrics
