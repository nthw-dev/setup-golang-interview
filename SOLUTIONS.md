# เฉลยและคำอธิบาย

> โค้ดเฉลยฉบับเต็มอยู่ในโฟลเดอร์ `solutions/` (ทุกข้อผ่าน `go test -race ./solutions/...`)
> ไฟล์นี้อธิบาย **แนวคิด**, **จุดที่คนมักพลาด** และ **คำตอบของคำถามต่อยอด**

---

## ข้อ 01 — ตัวแปร, Type Conversion, String/Rune

**โค้ดเต็ม:** [solutions/01_basics/basics.go](solutions/01_basics/basics.go)
<!-- code: solutions/01_basics/basics.go -->

### แนวคิด

```go
n, err := strconv.Atoi(strings.TrimSpace(s))
if err != nil {
    return 0, fmt.Errorf("index %d: %w", i, err) // %w = wrap
}
```

- **`%w` กับ `%v` ต่างกัน:** `%w` เก็บ error เดิมไว้ข้างใน ผู้เรียกจึงใช้ `errors.Is` หรือ `errors.As` ดึงออกมาได้ ส่วน `%v` แปลงเป็นข้อความอย่างเดียว ข้อมูล type หายไป
- **Integer division:** `3 / 2 == 1` ต้องแปลงก่อนหาร `float64(sum) / float64(len(nums))`
  Go ไม่แปลง type ให้อัตโนมัติ (ไม่มี implicit conversion) แม้แต่ระหว่าง `int` กับ `int64`
- **string คือ byte ไม่ใช่ตัวอักษร:** string ของ Go เป็น byte sequence แบบ UTF-8 ตัวอักษรไทยใช้ 3 byte ต่อตัว
  - `len("สวัสดี")` = 18 (byte)
  - `utf8.RuneCountInString("สวัสดี")` = 6 (rune)
  - `s[i]` ได้ **byte** แต่ `for _, r := range s` ได้ **rune**
  - `Reverse` ต้องแปลงเป็น `[]rune` ก่อน ถ้าสลับทีละ byte ตัวอักษรหลาย byte จะพัง

### คำถามต่อยอด

- **Zero value:** `int`=0, `float64`=0, `bool`=false, `string`="", ส่วน pointer, slice, map, channel, func และ interface = `nil`
  nil slice ใช้ `len`, `append` และ `range` ได้เลย แต่**เขียนลง nil map จะ panic**
- **`var x int` กับ `x := 0`:** ได้ผลเหมือนกัน แต่ `:=` ใช้ได้เฉพาะในฟังก์ชัน และต้องมีตัวแปรใหม่อย่างน้อย 1 ตัวทางซ้าย `var` ใช้ได้ทั้งระดับ package และตอนที่อยากได้ zero value ชัด ๆ

---

## ข้อ 02 — Pointer และ Pointer Receiver

**โค้ดเต็ม:** [solutions/02_pointers/pointers.go](solutions/02_pointers/pointers.go)
<!-- code: solutions/02_pointers/pointers.go -->

### บั๊กที่ 1: Value receiver

```go
func (a Account) Deposit(amount int) error { a.Balance += amount } // ❌ แก้แค่สำเนา
func (a *Account) Deposit(amount int) error { a.Balance += amount } // ✅
```

Go ส่งทุกอย่างแบบ **pass by value** (copy) receiver ก็เช่นกัน ถ้าต้องการแก้ค่าตัวจริงต้องใช้ pointer
เรียก `acc.Deposit(50)` ได้ทั้งที่ `acc` ไม่ใช่ pointer เพราะ Go แปลงเป็น `(&acc).Deposit(50)` ให้เองเมื่อ `acc` addressable

### บั๊กที่ 2: range ได้สำเนา

```go
for _, a := range accounts { a.Balance += bonus }   // ❌ a คือสำเนา
for i := range accounts { accounts[i].Balance += bonus } // ✅
```

### Transfer ต้องเป็น atomic (ในเชิง logic)

ถอนก่อน ถ้าถอนไม่ผ่านให้ return ทันที บัญชีไหนก็ยังไม่เปลี่ยน ส่วน `Deposit` หลังจากนั้นไม่มีทาง fail เพราะ `amount` ผ่านการตรวจใน `Withdraw` แล้ว

### คำถามต่อยอด

- **ใช้ pointer receiver เมื่อ:** (1) ต้องแก้ค่าใน struct (2) struct ใหญ่ ไม่อยาก copy (3) struct มี field ที่ห้าม copy เช่น `sync.Mutex`
  ถ้ามี method ไหนเป็น pointer receiver แล้ว ควรให้ทุก method ของ type นั้นเป็นแบบเดียวกันเพื่อความสม่ำเสมอ
- **Return pointer ของตัวแปร local ได้ไหม:** ได้ (ต่างจาก C) compiler ทำ **escape analysis** แล้วย้ายตัวแปรไปไว้บน heap ให้เอง ดูผลได้ด้วย `go build -gcflags=-m`
- **Go ไม่มี pointer arithmetic** (ยกเว้นใช้ `unsafe`) จึงปลอดภัยกว่า C

---

## ข้อ 03 — Slice, Map, Sort

**โค้ดเต็ม:** [solutions/03_collections/collections.go](solutions/03_collections/collections.go)
<!-- code: solutions/03_collections/collections.go -->

### แนวคิดสำคัญ

- **ลำดับของ map สุ่มโดยตั้งใจ:** Go สุ่มลำดับการวนลูป map ทุกครั้ง ถ้าต้องการผลที่แน่นอนต้องดึง key ออกมาแล้ว sort เอง
- **Sort หลายเงื่อนไข:** เทียบ count ก่อน ถ้าเท่ากันค่อยเทียบ word
- **Set ใน Go:** ใช้ `map[T]struct{}` เพราะ `struct{}` มีขนาด 0 byte
- **Full slice expression `s[i:j:k]`:** จำกัด `cap = k - i`

```go
nums := []int{1, 2, 3, 4}
a := nums[0:2]        // len 2, cap 4 → ยังแชร์ array กับ nums
a = append(a, 99)     // เขียนทับ nums[2]!  nums = [1 2 99 4]

b := nums[0:2:2]      // len 2, cap 2
b = append(b, 99)     // cap ไม่พอ → สร้าง array ใหม่ nums ไม่เปลี่ยน
```

### คำถามต่อยอด

- **Array กับ Slice:** array มีขนาดคงที่และขนาดเป็นส่วนหนึ่งของ type (`[3]int` ≠ `[4]int`) และ copy ทั้งก้อนเมื่อส่งต่อ
  slice เป็น header 3 ค่า `{ptr, len, cap}` ที่ชี้ไป array ข้างหลัง copy slice จึงยังแชร์ข้อมูลเดิม
- **`append` เมื่อ cap เต็ม:** Go จะสร้าง array ใหม่ (ประมาณ ×2 ตอนเล็ก และ ×1.25 ตอนใหญ่) แล้ว copy ข้อมูลไป
- **Map ไม่ thread-safe:** อ่านและเขียนพร้อมกันจะเจอ fatal error (ดูข้อ 08)

---

## ข้อ 04 — Struct, Interface, Polymorphism

**โค้ดเต็ม:** [solutions/04_interfaces/shapes.go](solutions/04_interfaces/shapes.go)
<!-- code: solutions/04_interfaces/shapes.go -->

### แนวคิด

- **Implicit interface:** ไม่ต้องประกาศ `implements` ขอแค่ type มี method ครบก็ใช้เป็น interface นั้นได้ ข้อดีคือ interface ประกาศฝั่ง "ผู้ใช้" ได้ (consumer-side interface) ทำให้ mock ง่าย
- **Compile-time check:** `var _ Shape = Rectangle{}` ถ้า Rectangle มี method ไม่ครบ จะ compile ไม่ผ่านทันที
- **Type switch:**

```go
switch v := s.(type) {
case Rectangle: // v เป็น Rectangle
case Circle:    // v เป็น Circle
default:
}
```

### คำถามต่อยอด

- **Go ไม่มี method/function overloading** ใช้แทนได้ด้วย (1) ตั้งชื่อต่างกัน เช่น `ParseInt` / `ParseFloat` (2) variadic `...T` (3) interface (4) generics (5) functional options pattern
- **nil interface กับ interface ที่เก็บ nil pointer** (กับดักยอดฮิต):

```go
func find() error {
    var e *MyErr = nil
    return e          // interface{type=*MyErr, value=nil} ≠ nil
}
find() != nil // true! 😱
```

interface จะเป็น nil ก็ต่อเมื่อทั้ง **type และ value** เป็น nil ดังนั้นถ้าไม่มี error ให้ `return nil` ตรง ๆ

---

## ข้อ 05 — Error Handling และ Custom Error

**โค้ดเต็ม:** [solutions/05_errors/store.go](solutions/05_errors/store.go)
<!-- code: solutions/05_errors/store.go -->

### รูปแบบ error 3 แบบในข้อนี้

| แบบ | ตัวอย่าง | ผู้เรียกตรวจด้วย | ใช้เมื่อ |
|---|---|---|---|
| Sentinel error | `var ErrNotFound = errors.New(...)` | `errors.Is(err, ErrNotFound)` | สนใจแค่ว่า "เป็น error ชนิดนี้ไหม" |
| Custom type | `type ValidationError struct{...}` | `errors.As(err, &ve)` | ต้องการข้อมูลเพิ่ม เช่น field ไหนผิด |
| Wrapping | `fmt.Errorf("get user %d: %w", id, err)` | ทั้ง Is และ As ยังทำงาน | เพิ่ม context ระหว่างส่ง error ขึ้นไป |

### จุดที่มักพลาด

- ใช้ `err == ErrNotFound` ซึ่งจะพังทันทีที่มีใคร wrap error ให้ใช้ `errors.Is` เสมอ
- `errors.As` ต้องส่ง **pointer ไปยังตัวแปร** ที่เป็น type เดียวกับตัวที่ implement `error` ในข้อนี้คือ `var ve *ValidationError; errors.As(err, &ve)`
- ข้อความ error ควรขึ้นต้นด้วยตัวพิมพ์เล็กและไม่ลงท้ายด้วยจุด (Go convention) เพราะมักถูกนำไปต่อกับข้อความอื่น

### คำถามต่อยอด

- **ทำไม Go ไม่มี exception:** เพื่อให้เส้นทางของ error **เห็นชัดในโค้ด** ทุกจุดที่ fail ได้ต้องจัดการหรือส่งต่อ ไม่มี control flow ซ่อนอยู่
- **`panic` ใช้เมื่อ:** เกิดสิ่งที่ไม่ควรเกิดขึ้นเลย (programmer error) หรือตอน init แล้ว config พัง ไม่ใช้กับ error ปกติ ส่วน `recover` ใช้ได้เฉพาะใน `defer` เช่น middleware ที่กัน server ล่ม

---

## ข้อ 06 — JSON และ Struct Tags

**โค้ดเต็ม:** [solutions/06_json/orders.go](solutions/06_json/orders.go)
<!-- code: solutions/06_json/orders.go -->

### แนวคิด

```go
type Order struct {
    ID       string `json:"id"`
    Customer string `json:"customer_name"`   // ชื่อ key ต่างจากชื่อ field
    Items    []Item `json:"items"`
    Note     string `json:"note,omitempty"`  // ค่าว่างจะไม่ออกใน JSON
}
```

- `encoding/json` ใช้ reflection และมองเห็น**เฉพาะ exported field** (ตัวพิมพ์ใหญ่)
- ถ้าไม่มี tag ตอน **Unmarshal** จะจับคู่ชื่อแบบไม่สนตัวพิมพ์ (`"id"` เข้า `ID` ได้) แต่ `customer_name` ไม่มีทางเข้า `Customer` ได้
- ถ้าไม่มี tag ตอน **Marshal** key จะออกมาเป็นชื่อ field ตรง ๆ (`"Customer"`)
- Test ใช้ **raw string literal** (backtick) เก็บ JSON หลายบรรทัดได้โดยไม่ต้อง escape

### คำถามต่อยอด

- **JSON ใหญ่มาก:** ใช้ `json.NewDecoder(r)` อ่านแบบ stream และใช้ `dec.Token()` + `dec.More()` วน decode ทีละ element แทนการโหลดทั้งไฟล์เข้า memory
- **Unmarshal ลง `map[string]any`:** ตัวเลขทุกตัวเป็น `float64` (ใช้ `dec.UseNumber()` ถ้าต้องการความแม่นยำ)
- **ปฏิเสธ field ที่ไม่รู้จัก:** `dec.DisallowUnknownFields()`

---

## ข้อ 07 — Worker Pool

**โค้ดเต็ม:** [solutions/07_workerpool/workerpool.go](solutions/07_workerpool/workerpool.go)
<!-- code: solutions/07_workerpool/workerpool.go -->

### แนวคิด

```
        jobs (index)            results[i]
main ──► [0,1,2,...] ──► worker 1 ──┐
                     ──► worker 2 ──┼──► slice ที่จองไว้ล่วงหน้า
                     ──► worker 3 ──┘
close(jobs) → worker ออกจาก range → wg.Done() → wg.Wait() คืนผล
```

- **คงลำดับได้:** เพราะส่ง index ไปด้วย แต่ละ worker เขียนลง `results[i]` คนละช่อง จึง**ไม่มี data race** (เขียน element ต่างกันของ slice พร้อมกันได้)
- **จำกัด concurrency:** จำนวน goroutine = `workers` ตัวพอดี
- **`close(jobs)`** ทำให้ `for i := range jobs` จบ ถ้าลืม worker จะค้างและ `wg.Wait()` จะ deadlock
- Go 1.25+ ใช้ `wg.Go(func(){...})` แทน `wg.Add(1)` + `defer wg.Done()` ได้

### คำถามต่อยอด

- **Goroutine กับ OS thread:** goroutine เริ่มที่ stack ~2KB และโตได้ สร้างได้เป็นแสนตัว runtime จัดตารางแบบ M:N (ดูทฤษฎีข้อ 45) ส่วน OS thread มี stack ~1–8MB และการ context switch แพงกว่า
- **Unbuffered channel:** ผู้ส่งจะรอจนมีผู้รับ (synchronization) ส่วน **buffered** ส่งได้เลยจน buffer เต็ม
- **ถ้า fn panic:** panic ใน goroutine จะทำให้ทั้งโปรแกรมล่ม ต้อง `defer func(){ if r := recover(); r != nil {...} }()` ในแต่ละ worker
- **ต้องคืน error:** ใช้ `golang.org/x/sync/errgroup` ซึ่งมี `g.SetLimit(n)` จำกัด concurrency และคืน error ตัวแรกได้

---

## ข้อ 08 — Race Condition และ Mutex

**โค้ดเต็ม:** [solutions/08_mutex/counter.go](solutions/08_mutex/counter.go)
<!-- code: solutions/08_mutex/counter.go -->

### แนวคิด

- `c.m[key]++` คือ **อ่าน → บวก → เขียน** ไม่ใช่ atomic เมื่อหลาย goroutine ทำพร้อมกัน ค่าจะหาย และ runtime ของ Go จะตรวจเจอการเขียน map พร้อมกันแล้ว crash ทันที
- `mu.Lock()` + `defer mu.Unlock()` ทำให้แน่ใจว่า unlock เสมอแม้จะ return หรือ panic กลางทาง
- **Snapshot ต้อง copy:** map เป็น reference type ถ้าคืน `c.m` ออกไปตรง ๆ คนนอกจะแก้ข้อมูลได้โดยไม่ผ่าน lock ซึ่งเป็น race อีกแบบ (Go 1.21+ ใช้ `maps.Clone` ได้)
- `go test -race` ใช้ **race detector** ตรวจการเข้าถึง memory พร้อมกัน ควรเปิดใน CI เสมอ

### คำถามต่อยอด

- **`RWMutex`:** อ่านพร้อมกันได้หลายตัว (`RLock`) แต่เขียนได้ทีละตัว เหมาะกับงานที่อ่านเยอะกว่าเขียนมาก
- **ทางเลือกอื่น:** `atomic.Int64` สำหรับตัวนับตัวเดียว, `sync.Map` สำหรับ key ที่เขียนครั้งเดียวแต่อ่านบ่อย หรือ goroutine เดียวที่เป็นเจ้าของ state แล้วให้ตัวอื่นสื่อสารผ่าน channel ("Don't communicate by sharing memory; share memory by communicating")
- **ห้าม copy Mutex:** copy แล้วจะได้ lock คนละตัว `go vet` (copylocks) จะเตือน นี่คือเหตุผลที่ struct ที่มี Mutex ต้องใช้ pointer receiver

---

## ข้อ 09 — Pipeline, Fan-in และ Goroutine Leak

**โค้ดเต็ม:** [solutions/09_pipeline/pipeline.go](solutions/09_pipeline/pipeline.go)
<!-- code: solutions/09_pipeline/pipeline.go -->

### ทำไมโค้ดเดิม leak

```go
for _, n := range nums {
    out <- n    // ถ้าไม่มีใครอ่านแล้ว บรรทัดนี้จะบล็อกตลอดไป
}
```

goroutine ที่บล็อกค้างจะ**ไม่ถูก GC เก็บ** เพราะ runtime ไม่รู้ว่าจะมีคนมาอ่านหรือไม่ ทั้ง stack และตัวแปรที่มันถืออยู่จึงค้างใน memory ตลอดไป

### วิธีแก้: select กับ ctx.Done()

```go
select {
case out <- n:
case <-ctx.Done():
    return   // defer close(out) ทำงาน → stage ถัดไปก็จบตาม
}
```

### Merge (fan-in)

- 1 goroutine ต่อ 1 input channel คอย forward ค่าไปที่ `out`
- **ต้องแยก goroutine** ออกมาทำ `wg.Wait(); close(out)` ถ้าเรียก `wg.Wait()` ใน Merge ตรง ๆ จะ deadlock เพราะยังไม่ได้คืน `out` ให้ใครอ่าน

### คำถามต่อยอด

- **หา leak ใน production:** `net/http/pprof` → `/debug/pprof/goroutine?debug=2` ดูว่า goroutine ค้างที่บรรทัดไหน หรือดู metric `runtime.NumGoroutine()` ว่าโตขึ้นเรื่อย ๆ หรือไม่ ใน test ใช้ `go.uber.org/goleak`
- **ผู้ส่งเป็นคน close:** ส่งค่าลง channel ที่ close แล้วจะ **panic** ส่วนการอ่านจาก channel ที่ close แล้วได้ zero value โดยไม่ panic ผู้ส่งเป็นคนเดียวที่รู้ว่า "ไม่มีข้อมูลแล้ว"

---

## ข้อ 10 — Timeout ด้วย Context

**โค้ดเต็ม:** [solutions/10_timeout/timeout.go](solutions/10_timeout/timeout.go)
<!-- code: solutions/10_timeout/timeout.go -->

### หัวใจของข้อนี้: buffered channel

```go
ch := make(chan result, 1)   // ✅ buffer 1
go func() { v, err := fn(ctx); ch <- result{v, err} }()

select {
case r := <-ch:       return r.val, r.err
case <-ctx.Done():    return "", ctx.Err()   // เราออกไปก่อนแล้ว...
}
```

ถ้า channel เป็น **unbuffered** เมื่อเรา timeout ออกไป จะไม่มีใครอ่าน `ch` อีก goroutine ที่ทำ `ch <- ...` จะบล็อกตลอดไป (leak)
buffer 1 ช่องทำให้มันส่งได้ทันทีแล้วจบ ส่วนค่าที่ค้างอยู่ใน channel GC จะเก็บให้

### FirstSuccess

- `context.WithCancel` + `defer cancel()` เมื่อได้คำตอบแล้ว return ทำให้ตัวที่เหลือถูก cancel อัตโนมัติ
- buffer = `len(fns)` ทุก goroutine ส่งผลได้โดยไม่ค้าง แม้เราจะเลิกอ่านแล้ว
- `errors.Join` (Go 1.20+) รวมหลาย error ไว้ด้วยกัน และ `errors.Is` ยังหาได้ทุกตัว

### คำถามต่อยอด

- **ต้อง `defer cancel()` ทุกครั้ง** ไม่งั้น timer และ goroutine ภายในของ context จะค้างจนกว่า parent จะถูก cancel (`go vet` เตือนเรื่อง lostcancel)
- **ctx เป็น parameter ตัวแรก** เป็น convention เพื่อให้ cancellation และ deadline ไหลผ่านทุกชั้นของ call stack อย่างชัดเจน
- **ไม่ควรเก็บ ctx ไว้ใน struct** เพราะ context ผูกกับ request แต่ละครั้ง ควรส่งผ่าน parameter

---

## ข้อ 11 — LRU Cache (Generics)

**โค้ดเต็ม:** [solutions/11_lru/lru.go](solutions/11_lru/lru.go)
<!-- code: solutions/11_lru/lru.go -->

### โครงสร้างข้อมูล

```
map[K]*list.Element            doubly linked list
  "a" ──────────────────►  [front] a ⇄ c ⇄ b [back]
  "b" ──────────────────────────────────────┘  ↑ ตัวที่จะโดนเตะ
  "c" ─────────────────────────┘
```

| Operation | ขั้นตอน | Big-O |
|---|---|---|
| Get | หาใน map → `MoveToFront` | O(1) |
| Put (มีอยู่แล้ว) | อัปเดตค่า → `MoveToFront` | O(1) |
| Put (ใหม่) | `PushFront` → ถ้าเกินให้ลบ `Back()` **และลบจาก map** | O(1) |

- ต้องเก็บ **key ไว้ใน element ด้วย** ตอนเตะตัวท้ายจะได้รู้ว่าต้อง `delete` key ไหนออกจาก map
- `var zero V` ใช้คืน zero value ของ generic type ใด ๆ

### คำถามต่อยอด

- **Thread-safe:** ใส่ `sync.Mutex` แล้ว lock ทั้ง Get และ Put (**Get ก็ต้อง lock** เพราะมันแก้ list ด้วย `MoveToFront` ใช้ RWMutex ไม่ได้ช่วยในกรณีนี้)
- **TTL:** เก็บ `expiresAt` ใน entry ตรวจตอน Get แล้วลบถ้าหมดอายุ (lazy expiration) อาจมี goroutine คอยกวาดเป็นระยะด้วย
- **`comparable` กับ `any`:** key ของ map ต้องเทียบ `==` ได้ จึงต้องเป็น `comparable` ส่วน value เป็นอะไรก็ได้ (`any`)

---

## ข้อ 12 — เขียน Test หาบั๊ก

**โค้ดเต็ม:** [solutions/12_testing/palindrome.go](solutions/12_testing/palindrome.go) · [palindrome_test.go](solutions/12_testing/palindrome_test.go)
<!-- code: solutions/12_testing/palindrome.go -->
<!-- code: solutions/12_testing/palindrome_test.go -->

### บั๊กในโค้ดเดิม

1. **เทียบทีละ byte:** `"été"` = `c3 a9 74 c3 a9` เทียบ byte แรก `c3` กับ byte สุดท้าย `a9` ไม่เท่ากัน จึงได้ false ทั้งที่เป็น palindrome
2. **ไม่ข้ามช่องว่างและเครื่องหมาย:** `"A man, a plan..."` จึงได้ false

### เครื่องมือ test ของ Go ที่ควรรู้

| อย่าง | ชื่อฟังก์ชัน | รันด้วย |
|---|---|---|
| Unit test | `TestXxx(t *testing.T)` | `go test` |
| Subtest | `t.Run(name, func(t *testing.T){...})` | `go test -run TestX/name` |
| Benchmark | `BenchmarkXxx(b *testing.B)` + `b.Loop()` | `go test -bench=. -benchmem` |
| Example | `ExampleXxx()` + `// Output:` | `go test` (และแสดงใน `go doc`) |
| Fuzz | `FuzzXxx(f *testing.F)` | `go test -fuzz=FuzzXxx` |
| Coverage | — | `go test -cover` / `-coverprofile=c.out` |

### คำถามต่อยอด

- **Mock:** รับ dependency เป็น interface เล็ก ๆ (ประกาศฝั่งผู้ใช้) แล้วส่ง fake struct เข้าไปตอน test สำหรับ HTTP ใช้ `httptest.NewServer`
- **`t.Helper()`** ทำให้บรรทัดที่รายงาน error ชี้ไปที่ผู้เรียก helper **`t.Cleanup()`** ลงทะเบียนงานเก็บกวาดหลัง test จบ **`t.Parallel()`** ให้ subtest รันขนานกัน

---

## ข้อ 13 — ใช้ Memory อย่างมีประสิทธิภาพ

**โค้ดเต็ม:** [solutions/13_memory/memory.go](solutions/13_memory/memory.go)
<!-- code: solutions/13_memory/memory.go -->

### 1. Preallocate

`var out []int` แล้ว append ไป 1000 ครั้ง ต้องขยาย array ~10 รอบ (ทุกรอบ allocate ใหม่และ copy) แต่ `make([]int, 0, n)` allocate ครั้งเดียวจบ

### 2. strings.Builder

string ใน Go **immutable** ทุกครั้งที่ `s += x` จะสร้าง string ใหม่และ copy ของเดิมทั้งหมด ได้ O(n²)
`strings.Builder` ต่อลง `[]byte` ภายในแล้วแปลงเป็น string ตอนท้าย**โดยไม่ copy** ส่วน `strconv.AppendInt` เขียนตัวเลขลง buffer บน stack ไม่ต้องสร้าง string ชั่วคราว

### 3. Sub-slice ค้าง array ใหญ่

```go
big := readFile()      // 100 MB
head := big[:10]       // ❌ head ยังชี้ array 100MB → GC เก็บไม่ได้
head := make([]byte, 10); copy(head, big) // ✅ เหลือแค่ 10 byte
```

### คำถามต่อยอด

- **Stack กับ Heap:** stack ต่อ goroutine เร็วมากและคืนอัตโนมัติเมื่อฟังก์ชันจบ heap ต้องให้ GC จัดการ compiler ตัดสินว่าจะวางตัวแปรที่ไหนด้วย **escape analysis** (ถ้าตัวแปรถูกอ้างถึงหลังฟังก์ชันจบ เช่น return pointer จะ escape ไป heap) ดูผลได้ด้วย `go build -gcflags=-m`
- **Profile memory:** `go test -memprofile mem.out` หรือ `import _ "net/http/pprof"` แล้วใช้ `go tool pprof -http=: mem.out`
- **`sync.Pool`:** ใช้ reuse object ชั่วคราวที่สร้างบ่อย เช่น buffer ลดภาระของ GC

---

## ข้อ 14 — REST API ด้วย net/http

**โค้ดเต็ม:** [solutions/14_http/server.go](solutions/14_http/server.go)
<!-- code: solutions/14_http/server.go -->

### แนวคิด

- **Routing ของ Go 1.22+:** `mux.HandleFunc("GET /todos/{id}", h)` แล้วอ่าน path parameter ด้วย `r.PathValue("id")` ถ้า method ไม่ตรง mux จะตอบ 405 ให้เอง
- **Concurrency:** `net/http` สร้าง goroutine ใหม่ต่อ request จึงต้องใช้ `sync.RWMutex` ป้องกัน map
- **`[]` กับ `null`:** nil slice ถูก encode เป็น `null` ต้อง `make([]Todo, 0)`
- **ลำดับการเขียน header:** ต้อง `w.Header().Set(...)` **ก่อน** `w.WriteHeader(status)` เพราะหลังจากนั้น header จะถูกส่งไปแล้ว
- **Test ด้วย `httptest`:** `httptest.NewRequest` + `httptest.NewRecorder` เรียก `ServeHTTP` ตรง ๆ ได้โดยไม่ต้องเปิด port

### คำถามต่อยอด

- **Middleware:**

```go
func Logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

- **Graceful shutdown:** `signal.NotifyContext(ctx, os.Interrupt)` → เมื่อได้สัญญาณ เรียก `srv.Shutdown(ctx)` เพื่อหยุดรับ request ใหม่และรอ request ที่ค้างอยู่ให้เสร็จ
- **Timeout ที่ควรตั้ง:** `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout` (ค่า default คือไม่มี timeout เลย เสี่ยงโดนโจมตีแบบ Slowloris)
- **Layer:** handler (HTTP ↔ struct) → service (business logic) → repository (interface ของ storage) test แต่ละชั้นแยกกันด้วย fake repository
