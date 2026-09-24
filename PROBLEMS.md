# โจทย์ฝึก Live Coding — Golang Developer

ชุดโจทย์ 14 ข้อ ออกแบบจากหัวข้อใน [roadmap.sh/questions/golang](https://roadmap.sh/questions/golang)
เรียงจากง่ายไปยาก แต่ละข้อมี test เตรียมไว้แล้ว หน้าที่ของคุณคือเขียนโค้ดให้ test ผ่าน

## วิธีใช้

```bash
make test N=01        # รัน test ของข้อ 01 (โฟลเดอร์ exercises/01_*)
make solution N=01    # ดูว่าเฉลยผ่าน test อะไรบ้าง
make reset N=01       # ล้างโค้ดที่เขียนไป กลับเป็นโจทย์เริ่มต้น
make test-all         # รันทุกข้อ
```

> **กติกาแนะนำ:** จับเวลาตามที่แต่ละข้อกำหนด พูดความคิดออกมาดัง ๆ ขณะเขียน (แบบสัมภาษณ์จริง)
> อย่าเปิดเฉลย ([SOLUTIONS.md](SOLUTIONS.md) / `solutions/`) จนกว่าจะทำเสร็จหรือหมดเวลา

| ข้อ | หัวข้อ | ระดับ | เวลา | คำถาม roadmap ที่เกี่ยวข้อง |
|---|---|---|---|---|
| 01 | ตัวแปร, type conversion, string/rune | ง่าย | 10 นาที | 3, 4, 12, 14, 34, 41 |
| 02 | Pointer และ pointer receiver | ง่าย | 10 นาที | 5, 10, 11, 24, 43, 49 |
| 03 | Slice, map, sort | ง่าย | 15 นาที | 6, 13 |
| 04 | Struct, interface, polymorphism | ง่าย | 15 นาที | 17, 21, 23, 31, 42 |
| 05 | Error handling และ custom error | กลาง | 15 นาที | 22, 29, 33, 44 |
| 06 | JSON และ struct tags | กลาง | 15 นาที | 14, 19 |
| 07 | Worker pool (goroutine + channel) | กลาง | 20 นาที | 25, 30, 40, 45 |
| 08 | Race condition และ Mutex | กลาง | 10 นาที | 18, 48 |
| 09 | Pipeline, fan-in, กัน goroutine leak | ยาก | 25 นาที | 30, 37, 40 |
| 10 | Timeout ด้วย context | ยาก | 20 นาที | 25, 30, 37 |
| 11 | LRU Cache แบบ generics | ยาก | 25 นาที | 13, 46 |
| 12 | เขียน test หาบั๊ก | กลาง | 15 นาที | 8, 35 |
| 13 | ใช้ memory อย่างมีประสิทธิภาพ | กลาง | 15 นาที | 9, 15, 16, 32, 36, 37, 46 |
| 14 | REST API ด้วย net/http | ยาก | 30 นาที | 19, 22, 25 |

คำถามทฤษฎีทั้ง 50 ข้อพร้อมคำตอบภาษาไทยอยู่ใน [THEORY.md](THEORY.md)

---

## ข้อ 01 — ตัวแปร, Type Conversion, String/Rune

**ไฟล์:** `exercises/01_basics/basics.go` · **ระดับ:** ง่าย · **เวลา:** 10 นาที

เขียนฟังก์ชันต่อไปนี้

1. `SumStrings(items []string) (int, error)` แปลงข้อความทุกตัวเป็น `int` แล้วรวมกัน
   - ตัด space หัวท้ายก่อนแปลง เช่น `" 2 "` ให้ได้ `2`
   - ถ้าตัวไหนแปลงไม่ได้ ให้คืน error ที่มีข้อความ `index <i>` **และ wrap error เดิมไว้** (caller ต้องใช้ `errors.As` หา `*strconv.NumError` ได้)
2. `Average(nums []int) (float64, error)` หาค่าเฉลี่ย ถ้า slice ว่างให้คืน `ErrEmpty`
3. `Reverse(s string) string` กลับลำดับตัวอักษร ต้องรองรับภาษาไทยและ emoji
4. `CharCount(s string) (bytes, runes int)` คืนจำนวน byte และจำนวนตัวอักษร

**ตัวอย่าง**

```go
SumStrings([]string{"1", " 2 ", "39"}) // 42, nil
SumStrings([]string{"1", "2", "abc"})  // 0, "index 2: strconv.Atoi: ..."
Average([]int{1, 2})                   // 1.5, nil   (ไม่ใช่ 1!)
Reverse("Go🚀")                         // "🚀oG"
CharCount("สวัสดี")                     // 18, 6
```

<details><summary>💡 Hint</summary>

- `strconv.Atoi`, `strings.TrimSpace`, `fmt.Errorf("...: %w", err)`
- `int / int` ได้ `int` ต้องแปลงเป็น `float64` ก่อนหาร
- `string` คือ byte sequence (UTF-8) ใช้ `[]rune(s)` หรือ `utf8.RuneCountInString`

</details>

**คำถามต่อยอดที่ interviewer อาจถาม:** zero value ของ `string`, `slice`, `map`, `pointer` คืออะไร? / `var x int` ต่างจาก `x := 0` ยังไง?

---

## ข้อ 02 — Pointer และ Pointer Receiver

**ไฟล์:** `exercises/02_pointers/pointers.go` · **ระดับ:** ง่าย · **เวลา:** 10 นาที

⚠️ โค้ดตั้งต้นมี **บั๊กที่จงใจใส่ไว้ 2 จุด** ให้หาสาเหตุ อธิบาย แล้วแก้

1. `Swap(a, b *int)` สลับค่าที่ pointer สองตัวชี้อยู่
2. `(*Account).Deposit(amount int) error` ฝากเงิน (`amount <= 0` → `ErrInvalidAmount`)
   **บั๊ก:** ฝากแล้วยอดเงินไม่เปลี่ยน
3. `(*Account).Withdraw(amount int) error` ถอนเงิน
   `amount <= 0` → `ErrInvalidAmount`, เงินไม่พอ → `ErrInsufficientFunds`
4. `Transfer(from, to *Account, amount int) error` โอนเงิน ถ้าไม่สำเร็จ **ห้ามมีบัญชีไหนเปลี่ยน**
5. `ApplyBonus(accounts []Account, bonus int)` บวก bonus ให้ทุกบัญชี
   **บั๊ก:** รันแล้วไม่มีใครได้ bonus

<details><summary>💡 Hint</summary>

- Value receiver `(a Account)` ได้ **สำเนา** ของ struct
- `for _, a := range accounts` ตัว `a` ก็เป็น **สำเนา** เหมือนกัน

</details>

**คำถามต่อยอด:** เมื่อไหร่ควรใช้ pointer receiver? / ถ้า method หนึ่งเป็น pointer receiver ควรทำให้ทุก method เป็นแบบเดียวกันไหม? / return pointer ของตัวแปร local ได้ไหม (escape analysis)?

---

## ข้อ 03 — Slice, Map, Sort

**ไฟล์:** `exercises/03_collections/collections.go` · **ระดับ:** ง่าย · **เวลา:** 15 นาที

1. `WordFrequency(text string) map[string]int` นับความถี่ของคำ โดยไม่สนตัวพิมพ์เล็ก/ใหญ่
   "คำ" คือกลุ่มตัวอักษรหรือตัวเลขที่อยู่ติดกัน ตัวอื่นทั้งหมดถือเป็นตัวคั่น
2. `TopK(freq map[string]int, k int) []string` คืน k คำที่พบบ่อยที่สุด เรียงจากมากไปน้อย
   ถ้าจำนวนเท่ากันให้เรียง a→z
3. `Unique(nums []int) []int` ลบตัวซ้ำโดยคงลำดับเดิม
4. `Chunk(nums []int, size int) [][]int` แบ่ง slice เป็นก้อนละ `size` ตัว (`size <= 0` → `nil`)
   **เงื่อนไขพิเศษ:** ถ้า `append` เพิ่มลง chunk แรก ต้องไม่ไปทับข้อมูลของ chunk ที่สอง

**ตัวอย่าง**

```go
WordFrequency("Go is fun. go, GO! Is it?") // map[go:3 is:2 fun:1 it:1]
TopK(map[string]int{"go":3,"is":2,"fun":2,"it":1}, 3) // [go fun is]
Unique([]int{3, 1, 3, 2, 1, 5})           // [3 1 2 5]
Chunk([]int{1, 2, 3, 4, 5}, 2)            // [[1 2] [3 4] [5]]
```

<details><summary>💡 Hint</summary>

- `strings.FieldsFunc`, `unicode.IsLetter`, `unicode.IsDigit`
- การวนลูป map **ไม่มีลำดับ** ต้อง sort เอง (`sort.Slice` หรือ `slices.SortFunc`)
- `map[int]struct{}` ใช้เป็น set
- full slice expression: `s[low:high:max]`

</details>

**คำถามต่อยอด:** slice ต่างจาก array ยังไง? / `len` กับ `cap` คืออะไร? / map thread-safe ไหม?

---

## ข้อ 04 — Struct, Interface, Polymorphism

**ไฟล์:** `exercises/04_interfaces/shapes.go` · **ระดับ:** ง่าย · **เวลา:** 15 นาที

มี interface `Shape { Area() float64; Perimeter() float64 }` และ struct `Rectangle{Width, Height}`, `Circle{Radius}`

1. เขียน `Area()` และ `Perimeter()` ของทั้งสอง struct ให้ถูกต้อง
2. `TotalArea(shapes ...Shape) float64` รวมพื้นที่ รับกี่ตัวก็ได้ (variadic)
3. `Largest(shapes []Shape) (Shape, bool)` คืน shape ที่ใหญ่ที่สุด ถ้า slice ว่างให้คืน `nil, false`
4. `Describe(s Shape) string` ใช้ **type switch**
   - `Rectangle{2,3}` → `"rectangle 2x3"`
   - `Circle{1.5}` → `"circle r=1.5"`
   - type อื่น → `"unknown shape"`

<details><summary>💡 Hint</summary>

- Go implement interface แบบ implicit: แค่มี method ครบก็ใช้ได้
- `switch v := s.(type) { case Rectangle: ... }`
- `%g` ใน `fmt.Sprintf` จะตัดศูนย์ท้ายทศนิยมทิ้ง

</details>

**คำถามต่อยอด:** Go มี method overloading ไหม? ถ้าไม่มีใช้อะไรแทน? / `var _ Shape = Rectangle{}` มีไว้ทำอะไร? / interface ที่เป็น nil ต่างจาก interface ที่ข้างในเก็บ nil pointer ยังไง?

---

## ข้อ 05 — Error Handling และ Custom Error

**ไฟล์:** `exercises/05_errors/store.go` · **ระดับ:** กลาง · **เวลา:** 15 นาที

สร้าง in-memory user store

1. `(*ValidationError).Error()` ต้องคืนข้อความรูปแบบ `validation failed on <Field>: <Msg>`
2. `NewStore()` และออกแบบ field ของ `Store` เอง
3. `Create(name, email string) (User, error)` ID เริ่มที่ 1 แล้วเพิ่มทีละ 1
   - name ว่าง (หลัง trim) → `*ValidationError{Field: "name"}`
   - email ไม่มี `@` → `*ValidationError{Field: "email"}`
   - email ซ้ำ → error ที่ **wrap** `ErrDuplicateEmail`
4. `Get(id int) (User, error)` ไม่เจอให้คืน error ที่ **wrap** `ErrNotFound` และมี id อยู่ในข้อความ

test จะตรวจด้วย `errors.Is(err, ErrNotFound)` และ `errors.As(err, &ve)`

<details><summary>💡 Hint</summary>

- `fmt.Errorf("get user %d: %w", id, ErrNotFound)`
- คืน `&ValidationError{...}` (pointer) เพราะ method `Error()` เป็น pointer receiver
- ใช้ map เพิ่มอีกตัว index ด้วย email เพื่อตรวจซ้ำแบบ O(1)

</details>

**คำถามต่อยอด:** `errors.Is` กับ `errors.As` ต่างกันยังไง? / ทำไม Go ไม่มี exception? / `panic` ควรใช้เมื่อไหร่?

---

## ข้อ 06 — JSON และ Struct Tags

**ไฟล์:** `exercises/06_json/orders.go` · **ระดับ:** กลาง · **เวลา:** 15 นาที

รูปแบบ JSON ของ order:

```json
{
  "id": "A001",
  "customer_name": "Tah",
  "items": [{"sku": "KB-01", "qty": 2, "price": 10.5}],
  "note": "optional"
}
```

1. เพิ่ม **struct tag** ให้ `Item` และ `Order` ตรงกับ JSON ข้างบน ถ้า `note` ว่างต้อง**ไม่มี** key `"note"` ตอน marshal
2. `(Order).Total() float64` = ผลรวม `Qty * Price`
3. `ParseOrder(data []byte) (Order, error)` JSON ผิดรูปแบบ → error, ไม่มี `id` → `ErrMissingID`
4. `Summarize(data []byte) ([]byte, error)` รับ JSON array ของ orders แล้วคืน JSON:
   `{"count": 3, "revenue": 28.5, "by_customer": {"Tah": 23, "Bee": 5.5}}`

<details><summary>💡 Hint</summary>

- `` `json:"customer_name"` ``, `` `json:"note,omitempty"` ``
- `json.Unmarshal(data, &v)` และ `json.Marshal(v)`
- field ต้องขึ้นต้นด้วยตัวพิมพ์ใหญ่ (exported) json ถึงจะมองเห็น

</details>

**คำถามต่อยอด:** ถ้า JSON ใหญ่มาก (เช่น 1GB) จะอ่านยังไง? (`json.Decoder` แบบ stream) / unmarshal ลง `map[string]any` ตัวเลขจะเป็น type อะไร?

---

## ข้อ 07 — Worker Pool

**ไฟล์:** `exercises/07_workerpool/workerpool.go` · **ระดับ:** กลาง · **เวลา:** 20 นาที

`ParallelMap(nums []int, workers int, fn func(int) int) []int`

- เรียก `fn` กับทุกตัวใน `nums` โดยมี goroutine ทำงานพร้อมกัน **ไม่เกิน** `workers` ตัว
- ผลลัพธ์ต้อง**เรียงตามลำดับเดิม**
- ต้องผ่าน `go test -race`

test จะวัดว่า 8 งาน งานละ 50ms ใช้ 4 workers ต้องเสร็จเร็ว (ประมาณ 100ms) และต้องไม่มีงานรันพร้อมกันเกิน 4 ตัว

<details><summary>💡 Hint</summary>

- ส่ง **index** ของงานผ่าน channel แล้วให้ worker เขียนผลลง `results[i]`
- ใช้ `close(jobs)` บอก worker ว่าหมดงาน และ `sync.WaitGroup` รอทุก worker จบ

</details>

**คำถามต่อยอด:** goroutine ต่างจาก OS thread ยังไง? / buffered กับ unbuffered channel ต่างกันยังไง? / ถ้า `fn` อาจ panic จะทำยังไง? / ถ้าต้องคืน error ด้วยล่ะ (`errgroup`)?

---

## ข้อ 08 — Race Condition และ Mutex

**ไฟล์:** `exercises/08_mutex/counter.go` · **ระดับ:** กลาง · **เวลา:** 10 นาที

`Counter` ตัวนี้พังเมื่อถูกเรียกจากหลาย goroutine (`fatal error: concurrent map writes`)

1. แก้ให้ `Inc` และ `Value` ปลอดภัยเมื่อใช้พร้อมกัน
2. แก้ `Snapshot()` ให้คืน**สำเนา** ถ้าผู้เรียกแก้ map ที่ได้คืนไป ต้องไม่กระทบ Counter

ลองรัน `go test -race ./exercises/08_mutex` ทั้งก่อนและหลังแก้

<details><summary>💡 Hint</summary>

`sync.Mutex` + `defer c.mu.Unlock()` (zero value ของ Mutex ใช้ได้เลย ไม่ต้อง init)

</details>

**คำถามต่อยอด:** `sync.Mutex` กับ `sync.RWMutex` ต่างกันยังไง? / ใช้ `sync.Map` หรือ `atomic` แทนได้ไหม? / ทำไมห้าม copy struct ที่มี Mutex (`go vet` เตือนเรื่องอะไร)?

---

## ข้อ 09 — Pipeline, Fan-in และ Goroutine Leak

**ไฟล์:** `exercises/09_pipeline/pipeline.go` · **ระดับ:** ยาก · **เวลา:** 25 นาที

`Generate` และ `Square` ในโค้ดตั้งต้นให้ผลถูกต้อง แต่ถ้าผู้อ่าน**เลิกอ่านกลางทาง** goroutine จะค้างตลอดไป (goroutine leak)

1. แก้ `Generate(ctx, nums...)` และ `Square(ctx, in)` ให้หยุดทันทีเมื่อ `ctx` ถูก cancel
2. เขียน `Merge(ctx, cs ...<-chan int) <-chan int` (fan-in) รวมหลาย channel เป็นช่องเดียว
   ต้อง close channel ผลลัพธ์เมื่อ input **ทุกตัว** ถูก close แล้ว

```go
src := Generate(ctx, 1, 2, ..., 10)
// fan-out: 3 workers อ่านจาก src เดียวกัน → fan-in ด้วย Merge
out := Merge(ctx, Square(ctx, src), Square(ctx, src), Square(ctx, src))
// sum(out) == 385
```

<details><summary>💡 Hint</summary>

```go
select {
case out <- n:
case <-ctx.Done():
    return
}
```

ใน `Merge` ให้แยก goroutine ที่ทำ `wg.Wait(); close(out)` ออกมา

</details>

**คำถามต่อยอด:** จะตรวจหา goroutine leak ใน production ยังไง? (pprof goroutine profile, `goleak`) / ทำไมผู้ส่งควรเป็นคน close channel ไม่ใช่ผู้รับ?

---

## ข้อ 10 — Timeout ด้วย Context

**ไฟล์:** `exercises/10_timeout/timeout.go` · **ระดับ:** ยาก · **เวลา:** 20 นาที

1. `CallWithTimeout(ctx, d, fn) (string, error)` เรียก `fn` แต่รอไม่เกิน `d`
   - ถ้าเกินเวลาต้องคืน error ที่ `errors.Is(err, context.DeadlineExceeded)`
   - **ห้ามมี goroutine ค้าง** แม้ `fn` จะไม่สนใจ ctx และทำงานต่อจนเสร็จหลัง timeout ไปแล้วก็ตาม
2. `FirstSuccess(ctx, fns...) (string, error)` รันทุก `fn` พร้อมกัน
   - คืนผลของ**ตัวแรกที่สำเร็จ** แล้ว cancel ตัวที่เหลือ
   - ถ้าทุกตัวล้มเหลว คืน error ที่รวมทุก error (`errors.Is` ต้องเจอทุกตัว)
   - ไม่มี fn เลย → `ErrNoFuncs`

<details><summary>💡 Hint</summary>

- `context.WithTimeout` / `context.WithCancel` แล้ว `defer cancel()` เสมอ
- ส่งผลลัพธ์ผ่าน **buffered channel** เพื่อให้ goroutine ส่งได้แม้ไม่มีใครรออ่าน
- `errors.Join(errs...)`

</details>

**คำถามต่อยอด:** ทำไมต้อง `defer cancel()` ทุกครั้ง? / ควรส่ง `context.Context` เป็น parameter ตัวแรกเพราะอะไร? / ควรเก็บ ctx ไว้ใน struct ไหม?

---

## ข้อ 11 — LRU Cache (Generics)

**ไฟล์:** `exercises/11_lru/lru.go` · **ระดับ:** ยาก · **เวลา:** 25 นาที

เขียน `Cache[K comparable, V any]` ที่ `Get` และ `Put` เป็น **O(1)**

- `New[K, V](capacity int) *Cache[K, V]`
- `Get(key) (V, bool)` คืนค่าและทำให้ key นั้นเป็น "ใช้ล่าสุด"
- `Put(key, value)` เพิ่มหรืออัปเดตค่า (นับเป็นการใช้ล่าสุดด้วย) ถ้าเกิน capacity ให้ลบตัวที่**ไม่ได้ใช้นานที่สุด**
- `Len() int`

```go
c := New[string, int](2)
c.Put("a", 1); c.Put("b", 2)
c.Get("a")      // 1, true → a กลายเป็นใช้ล่าสุด
c.Put("c", 3)   // เกิน → เตะ b ออก
c.Get("b")      // 0, false
```

<details><summary>💡 Hint</summary>

`map[K]*list.Element` + `container/list` (doubly linked list) หน้าสุดคือตัวที่ใช้ล่าสุด ท้ายสุดคือตัวที่จะถูกเตะ
อย่าลืมลบ key ออกจาก map ตอนเตะด้วย

</details>

**คำถามต่อยอด:** ทำให้ thread-safe ยังไง? / เพิ่ม TTL ยังไง? / `comparable` กับ `any` ต่างกันยังไง?

---

## ข้อ 12 — เขียน Test หาบั๊ก

**ไฟล์:** `exercises/12_testing/palindrome.go` และ `palindrome_test.go` · **ระดับ:** กลาง · **เวลา:** 15 นาที

`IsPalindrome(s string) bool` ต้องคืน `true` ถ้าอ่านจากหน้าไปหลังและหลังไปหน้าได้เหมือนกัน
โดย**ไม่สน**ตัวพิมพ์เล็ก/ใหญ่ และ**ข้าม**ทุกอย่างที่ไม่ใช่ตัวอักษรหรือตัวเลข (รองรับ Unicode)

โค้ดที่ให้มามี**บั๊กอย่างน้อย 2 จุด** และ test ตั้งต้นผ่าน (เพราะมีแค่เคสเดียว)

1. เขียน **table-driven test** ด้วย `t.Run` อย่างน้อย 8 เคส ให้จับบั๊กได้
2. แก้ `IsPalindrome` ให้ test ผ่าน
3. เขียน `BenchmarkIsPalindrome` และ `ExampleIsPalindrome`
4. (โบนัส) เขียน `FuzzIsPalindrome`

```bash
go test -v ./exercises/12_testing
go test -bench=. -benchmem ./exercises/12_testing
go test -fuzz=FuzzIsPalindrome -fuzztime=10s ./exercises/12_testing
```

**คำถามต่อยอด:** จะ mock dependency (DB, HTTP) ใน Go ยังไง? / `t.Helper()`, `t.Cleanup()`, `t.Parallel()` ใช้ทำอะไร? / วัด coverage ยังไง?

---

## ข้อ 13 — ใช้ Memory อย่างมีประสิทธิภาพ

**ไฟล์:** `exercises/13_memory/memory.go` · **ระดับ:** กลาง · **เวลา:** 15 นาที

โค้ดตั้งต้นให้ผลถูกต้อง แต่ allocate memory เยอะเกินไป ให้ปรับจนผ่าน test ที่ใช้ `testing.AllocsPerRun`

1. `Squares(n int) []int` ต้อง allocate **1 ครั้ง**
2. `JoinInts(nums []int, sep string) string` ต้อง allocate **ไม่เกิน 3 ครั้ง** สำหรับ 100 ตัว (ของเดิมประมาณ 200 ครั้ง)
3. `Head(data []byte, n int) []byte` ผลลัพธ์ต้อง**ไม่แชร์ memory** กับ `data`
   (ลองคิดดูว่า ถ้า `data` ใหญ่ 100MB แล้วเราเก็บไว้แค่ 10 byte แรกด้วย `data[:10]` จะเกิดอะไรขึ้น)

```bash
go test -bench=. -benchmem ./exercises/13_memory   # ดู allocs/op
```

<details><summary>💡 Hint</summary>

- `make([]int, 0, n)`
- `strings.Builder` + `Grow` + `strconv.AppendInt`
- `make` แล้ว `copy`

</details>

**คำถามต่อยอด:** stack กับ heap ต่างกันยังไง? escape analysis คืออะไร (`go build -gcflags=-m`)? / จะ profile memory ด้วย pprof ยังไง? / `sync.Pool` ใช้เมื่อไหร่?

---

## ข้อ 14 — REST API ด้วย net/http

**ไฟล์:** `exercises/14_http/server.go` · **ระดับ:** ยาก · **เวลา:** 30 นาที

สร้าง Todo API แบบ in-memory โดยใช้แค่ standard library

| Method | Path | สำเร็จ | กรณี error |
|---|---|---|---|
| GET | `/todos` | 200 + array เรียงตาม id (ว่าง = `[]` ไม่ใช่ `null`) | — |
| POST | `/todos` body `{"title":"..."}` | 201 + todo ที่สร้าง (id เริ่มที่ 1) | JSON พัง / title ว่าง → 400 |
| GET | `/todos/{id}` | 200 + todo | ไม่เจอ → 404, id ไม่ใช่ตัวเลข → 400 |
| DELETE | `/todos/{id}` | 204 (ไม่มี body) | ไม่เจอ → 404, id ไม่ใช่ตัวเลข → 400 |

- ทุก error ตอบเป็น `{"error": "..."}`
- ทุก JSON response ต้องมี header `Content-Type: application/json`
- ต้องปลอดภัยเมื่อมี request เข้ามาพร้อมกัน

<details><summary>💡 Hint</summary>

- Go 1.22+ : `mux.HandleFunc("GET /todos/{id}", h)` แล้วอ่านด้วย `r.PathValue("id")`
- `json.NewDecoder(r.Body).Decode(&in)`
- `sync.RWMutex`
- nil slice ถูก encode เป็น `null` แต่ `make([]Todo, 0)` ได้ `[]`

</details>

**คำถามต่อยอด:** จะเพิ่ม middleware (logging, auth) ยังไง? / graceful shutdown ทำยังไง? / จะตั้ง timeout ของ `http.Server` ตัวไหนบ้าง? / จะแยก layer handler/service/repository ยังไง?
