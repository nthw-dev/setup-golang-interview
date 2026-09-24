# คำถามทฤษฎี Golang 50 ข้อ (พร้อมคำตอบ)

แปลและเรียบเรียงจาก [roadmap.sh/questions/golang](https://roadmap.sh/questions/golang) พร้อมเพิ่มรายละเอียดที่ interviewer มักถามต่อ
**วิธีฝึก:** อ่านคำถาม ตอบออกเสียงเองก่อน แล้วค่อยกดดูคำตอบ

---

## ระดับเริ่มต้น (Beginner)

### 1. Go คืออะไร มีจุดเด่นอะไรบ้าง?

<details><summary>ดูคำตอบ</summary>

ภาษา open-source จาก Google (2009) เป็น statically typed และ compile เป็น native binary ไฟล์เดียว จุดเด่นคือ

- **เรียบง่าย:** keyword น้อย (25 ตัว) อ่านง่าย มี `gofmt` จัด format ให้เป็นมาตรฐานเดียวกัน
- **Concurrency ในตัวภาษา:** goroutine + channel
- **มี Garbage Collector** ไม่ต้องจัดการ memory เอง
- **Compile เร็ว** และ cross-compile ง่าย (`GOOS=linux GOARCH=arm64 go build`)
- **Standard library ครบ** (`net/http`, `encoding/json`, `testing`) พร้อมเครื่องมือในตัว เช่น `go test`, `go vet`, `pprof`, race detector

</details>

### 2. `package main` และ `func main()` มีไว้ทำอะไร?

<details><summary>ดูคำตอบ</summary>

- `package main` บอก compiler ว่า package นี้ build เป็น **โปรแกรมที่รันได้** (executable) ถ้าเป็นชื่ออื่นจะเป็น library
- `func main()` คือ **จุดเริ่มต้น** ของโปรแกรม ไม่รับ argument และไม่คืนค่า (อ่าน argument ผ่าน `os.Args` และกำหนด exit code ด้วย `os.Exit`)
- ลำดับการทำงาน: init ตัวแปรระดับ package → `func init()` ของทุก package → `main()` เมื่อ `main` จบ โปรแกรมจบทันทีโดย**ไม่รอ goroutine อื่น**

</details>

### 3. ประกาศตัวแปรได้กี่แบบ และ zero value คืออะไร?

<details><summary>ดูคำตอบ</summary>

```go
var a int          // zero value = 0
var b = "hi"       // อนุมาน type
c := 3.14          // short declaration (ใช้ได้เฉพาะในฟังก์ชัน)
const Pi = 3.14159
```

**Zero value** คือค่าเริ่มต้นเมื่อไม่ได้กำหนดค่า Go ไม่มีตัวแปรที่ "ไม่ได้ initialize"

| Type | Zero value |
|---|---|
| ตัวเลข | `0` |
| `bool` | `false` |
| `string` | `""` |
| pointer, slice, map, channel, func, interface | `nil` |
| struct | ทุก field เป็น zero value ของตัวเอง |

แนวทางที่ดีคือออกแบบ type ให้ zero value ใช้งานได้เลย เช่น `sync.Mutex`, `bytes.Buffer`

</details>

### 4. อธิบาย data type พื้นฐานของ Go

<details><summary>ดูคำตอบ</summary>

- **Boolean:** `bool`
- **Integer:** `int`, `int8/16/32/64`, `uint...`, `uintptr` (`int` มีขนาด 64 bit บนระบบ 64-bit)
- **Float:** `float32`, `float64`
- **Complex:** `complex64`, `complex128`
- **String:** byte sequence แบบ UTF-8 แก้ไขไม่ได้ (immutable)
- **Alias:** `byte` = `uint8`, `rune` = `int32` (ใช้แทน Unicode code point)

Go เป็น **strongly typed** จึงบวก `int` กับ `int64` ตรง ๆ ไม่ได้ ต้องแปลงเองเสมอ

</details>

### 5. Pointer คืออะไร?

<details><summary>ดูคำตอบ</summary>

ตัวแปรที่เก็บ **address** ของตัวแปรอื่น

```go
x := 10
p := &x     // p เป็น *int ชี้ไปที่ x
*p = 20     // dereference → x = 20
```

ใช้เมื่อต้องการ (1) แก้ค่าต้นฉบับในฟังก์ชัน (2) ไม่อยาก copy struct ขนาดใหญ่ (3) แทน "ไม่มีค่า" ด้วย `nil`
Go ไม่มี pointer arithmetic และ return pointer ของตัวแปร local ได้อย่างปลอดภัย ดูฝึกได้ในข้อ 02

</details>

### 6. Composite data type มีอะไรบ้าง?

<details><summary>ดูคำตอบ</summary>

- **Array** `[3]int{1,2,3}` ขนาดคงที่ (ขนาดเป็นส่วนหนึ่งของ type) เป็น value type
- **Slice** `[]int{1,2,3}` ขนาดยืดได้ เป็น header `{ptr, len, cap}` ที่ชี้ไปยัง array
- **Map** `map[string]int{"a": 1}` เป็น hash table
- **Struct** รวม field หลายชนิดเข้าด้วยกัน

```go
type User struct { Name string; Age int }
u := User{Name: "Tah", Age: 30}
fmt.Println(u.Name)
```

ดูฝึกได้ในข้อ 03

</details>

### 7. จัดการ dependency ใน Go ยังไง?

<details><summary>ดูคำตอบ</summary>

ใช้ **Go Modules**

```bash
go mod init github.com/me/app     # สร้าง go.mod
go get github.com/lib/pq@v1.10.9  # เพิ่ม/อัปเดต dependency
go mod tidy                       # ลบตัวที่ไม่ใช้ และเพิ่มตัวที่ขาด
```

- `go.mod` เก็บชื่อ module, เวอร์ชัน Go และ dependency ที่ใช้
- `go.sum` เก็บ checksum ไว้ยืนยันว่าโค้ดที่โหลดมาไม่ถูกแก้
- commit ทั้งสองไฟล์เข้า git เพื่อให้ build ได้ผลเหมือนกันทุกเครื่อง

</details>

### 8. package `testing` มีบทบาทอะไร?

<details><summary>ดูคำตอบ</summary>

เป็นเครื่องมือ test ในตัวภาษา ไม่ต้องติดตั้ง framework เพิ่ม

- ไฟล์ลงท้าย `_test.go` และฟังก์ชันชื่อ `TestXxx(t *testing.T)`
- รันด้วย `go test ./...`
- รองรับ benchmark (`BenchmarkXxx`), example (`ExampleXxx`), fuzz (`FuzzXxx`), coverage (`-cover`) และ race detector (`-race`)
- สไตล์ที่นิยมคือ **table-driven test** + `t.Run`

ดูฝึกได้ในข้อ 12

</details>

### 9. Go จัดการ memory อย่างไร?

<details><summary>ดูคำตอบ</summary>

- **อัตโนมัติ:** compiler ตัดสินว่าจะวางตัวแปรบน **stack** หรือ **heap** ด้วย escape analysis
- **Garbage Collector** แบบ concurrent tri-color mark-and-sweep คืน memory บน heap ที่ไม่มีใครอ้างถึงแล้ว โดยทำงานไปพร้อมกับโปรแกรม ทำให้ pause สั้นมาก (ระดับ sub-millisecond)
- ปรับจูนได้ด้วย `GOGC` (ความถี่ของ GC) และ `GOMEMLIMIT` (เพดาน memory แบบ soft limit)

</details>

### 10. การเรียกฟังก์ชันใน Go ทำงานอย่างไร?

<details><summary>ดูคำตอบ</summary>

- Go ส่ง argument แบบ **pass by value เสมอ** คือ copy ค่าไปให้ฟังก์ชัน
- ถ้าส่ง pointer ก็คือ copy address ทำให้แก้ค่าต้นฉบับได้
- slice, map และ channel มี pointer อยู่ข้างใน การ copy จึงยังเห็นข้อมูลชุดเดียวกัน (แต่ `append` slice จนเกิน cap แล้วผู้เรียกจะมองไม่เห็นผล)
- runtime สร้าง stack frame ให้ทุกการเรียก ฟังก์ชันคืนค่าได้หลายค่า และใช้เป็น first-class value หรือ closure ได้

</details>

### 11. Value type กับ Reference type ต่างกันอย่างไร?

<details><summary>ดูคำตอบ</summary>

| Value type | "Reference-like" type |
|---|---|
| `int`, `float`, `bool`, `string`, array, struct | slice, map, channel, pointer, func, interface |
| ส่งต่อแล้วได้**สำเนาอิสระ** | ส่งต่อแล้ว**ยังชี้ข้อมูลชุดเดิม** |

ทางเทคนิค Go ไม่มี "reference" แบบ C++ ทุกอย่างคือ copy แต่บาง type ข้างในมี pointer อยู่ การ copy จึงยังแชร์ข้อมูล

</details>

### 12. แปลง type ทำอย่างไร?

<details><summary>ดูคำตอบ</summary>

ต้องแปลง **explicit** เสมอ

```go
f := float64(10)                // ตัวเลข → ตัวเลข
s := strconv.Itoa(42)           // int → string "42"
n, err := strconv.Atoi("42")    // string → int (อาจ error)
b := []byte("hi")               // string ↔ []byte
r := []rune("สวัสดี")            // string → []rune
```

⚠️ `string(65)` ได้ `"A"` ไม่ใช่ `"65"` (`go vet` เตือน) ถ้าต้องการตัวเลขเป็นข้อความต้องใช้ `strconv`

ส่วน **type assertion** `v, ok := x.(string)` ใช้กับ interface ไม่ใช่การแปลง type

</details>

### 13. Map คืออะไร ใช้ทำอะไร?

<details><summary>ดูคำตอบ</summary>

hash table ที่เก็บคู่ key/value ค้นหา เพิ่ม และลบได้เฉลี่ย O(1)

```go
m := make(map[string]int)
m["go"]++
v, ok := m["rust"]   // comma-ok: ok=false ถ้าไม่มี key
delete(m, "go")
```

ใช้ทำ cache, config, นับความถี่, set (`map[T]struct{}`) และ index ข้อมูล

ข้อควรรู้: ลำดับการวนลูปสุ่ม, **ไม่ thread-safe**, เขียนลง nil map จะ panic และ key ต้องเป็น type ที่ `comparable`

</details>

### 14. Raw string literal คืออะไร?

<details><summary>ดูคำตอบ</summary>

string ที่ครอบด้วย backtick `` ` `` ข้างในไม่ตีความ escape sequence (`\n` เป็นตัวอักษร 2 ตัว) และขึ้นบรรทัดใหม่ได้

```go
re := `^\d{3}-\d{4}$`          // regex ไม่ต้อง escape \ ซ้ำ
js := `{"name": "Tah"}`        // JSON ไม่ต้อง escape "
```

เหมาะกับ regex, JSON, SQL, template และ path ของ Windows ดูตัวอย่างใน test ของข้อ 06

</details>

### 15. Automatic memory management ช่วยป้องกัน memory leak อย่างไร?

<details><summary>ดูคำตอบ</summary>

GC คืน memory ของ object ที่**ไม่มีใครอ้างถึงแล้ว**ให้อัตโนมัติ จึงไม่มีปัญหาลืม `free` หรือ double free แบบ C

แต่ **ยังเกิด leak ได้** ถ้ายังมีการอ้างถึงอยู่ เช่น

- goroutine ที่บล็อกค้าง (ข้อ 09, 10)
- map หรือ cache ที่โตไม่หยุด
- sub-slice ที่ถือ array ใหญ่ไว้ (ข้อ 13)
- `time.Ticker` ที่ไม่ได้ `Stop()` (ก่อน Go 1.23)
- global variable ที่สะสมข้อมูลเรื่อย ๆ

</details>

### 16. Memory allocation ใน Go ทำงานอย่างไร?

<details><summary>ดูคำตอบ</summary>

- **Stack:** ตัวแปรที่ใช้แค่ในฟังก์ชัน จองและคืนเร็วมาก ไม่ต้องใช้ GC
- **Heap:** ตัวแปรที่ "escape" ออกนอกฟังก์ชัน เช่น return pointer, เก็บใน interface หรือ closure, หรือใหญ่เกินกว่าจะวางบน stack

compiler ตัดสินเองด้วย **escape analysis** ตรวจดูได้ด้วย

```bash
go build -gcflags=-m ./...
# ./x.go:10:9: &User{...} escapes to heap
```

`new(T)` และ `make(...)` ไม่ได้บังคับว่าต้องอยู่บน heap เสมอ

</details>

### 17. Struct คืออะไร และเข้าถึง field อย่างไร?

<details><summary>ดูคำตอบ</summary>

```go
type Point struct {
    X, Y int
}
p := Point{X: 1, Y: 2}
p.X = 10
pp := &p
pp.Y = 20   // Go dereference ให้อัตโนมัติ ไม่ต้องเขียน (*pp).Y
```

- field ที่ขึ้นต้นด้วยตัวพิมพ์ใหญ่จะ export ออกนอก package
- **Embedding** (`type Admin struct { User }`) ใช้ทำ composition แทน inheritance
- struct เทียบกันด้วย `==` ได้ถ้าทุก field comparable

</details>

### 18. Global variable คืออะไร?

<details><summary>ดูคำตอบ</summary>

ตัวแปรที่ประกาศนอกฟังก์ชัน (package-level) ใช้ได้ทั้ง package และจาก package อื่นถ้า export

**ควรหลีกเลี่ยง** เพราะ

- test ยาก เพราะ state ข้าม test กัน
- เสี่ยง data race เมื่อหลาย goroutine ใช้พร้อมกัน
- เกิด hidden dependency
- อยู่ใน memory ตลอดอายุโปรแกรม

ทางเลือกคือส่ง dependency ผ่าน struct หรือ parameter (dependency injection) ส่วนที่ยอมรับได้คือค่าคงที่ sentinel error และ regex ที่ compile ครั้งเดียว

</details>

### 19. ทำงานกับ JSON อย่างไร?

<details><summary>ดูคำตอบ</summary>

ใช้ package `encoding/json`

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    Pass  string `json:"-"`            // ไม่ออกใน JSON
}
var u User
err := json.Unmarshal(data, &u)        // JSON → struct
out, err := json.Marshal(u)            // struct → JSON
json.NewDecoder(r.Body).Decode(&u)     // แบบ stream
```

มองเห็นเฉพาะ exported field ดูฝึกได้ในข้อ 06

</details>

### 20. Reserved keyword คืออะไร?

<details><summary>ดูคำตอบ</summary>

คำที่ภาษาจองไว้ ใช้เป็นชื่อตัวแปรไม่ได้ Go มี 25 คำ

`break case chan const continue default defer else fallthrough for func go goto if import interface map package range return select struct switch type var`

ส่วน `int`, `string`, `true`, `nil`, `len`, `append` เป็น **predeclared identifier** ไม่ใช่ keyword จึงตั้งชื่อทับได้ แต่ไม่ควรทำ

</details>

### 21. Go รองรับ method overloading ไหม?

<details><summary>ดูคำตอบ</summary>

**ไม่รองรับ** ในหนึ่ง scope ชื่อฟังก์ชันซ้ำกันไม่ได้ (ตั้งใจให้เรียบง่าย) ใช้แทนได้ด้วย

- ตั้งชื่อต่างกัน: `strconv.ParseInt`, `strconv.ParseFloat`
- Variadic: `func Sum(nums ...int)`
- Interface: รับ `Shape` แล้วส่งอะไรก็ได้ที่ implement
- Generics: `func Max[T cmp.Ordered](a, b T) T`
- Functional options: `NewServer(WithPort(8080), WithTLS())`

</details>

---

## ระดับกลาง (Intermediate)

### 22. Implement error handling อย่างไร?

<details><summary>ดูคำตอบ</summary>

error คือ interface `type error interface { Error() string }` คืนเป็นค่าสุดท้ายของฟังก์ชัน แล้วตรวจทันที

```go
f, err := os.Open(path)
if err != nil {
    return fmt.Errorf("open config %s: %w", path, err) // wrap + เพิ่ม context
}
defer f.Close()
```

แนวทางที่ดี:

- จัดการ error **ครั้งเดียว** เลือกว่าจะ log หรือจะ return อย่าทำทั้งสองอย่าง
- ใช้ `errors.Is` และ `errors.As` แทน `==`
- สร้าง custom error type เมื่อต้องการส่งข้อมูลเพิ่ม

ดูฝึกได้ในข้อ 05

</details>

### 23. Method signature คืออะไร?

<details><summary>ดูคำตอบ</summary>

"หน้าตา" ของ method ประกอบด้วย receiver, ชื่อ, parameter และ return type

```go
func (a *Account) Withdraw(amount int) error
//    ^receiver   ^ชื่อ     ^params      ^return
```

signature กำหนดว่า type นั้น implement interface ไหนได้ ต้องตรงกันทุกส่วน รวมถึงชนิดของ receiver ด้วย: method ที่มี pointer receiver จะนับอยู่ใน method set ของ `*T` เท่านั้น ไม่นับใน `T`

</details>

### 24. อธิบาย pointer แบบละเอียด

<details><summary>ดูคำตอบ</summary>

- `&x` คือ address ของ x และ `*p` คือค่าที่ p ชี้อยู่ (dereference)
- zero value ของ pointer คือ `nil` การ dereference nil จะ **panic**
- `new(T)` คืน `*T` ที่ชี้ไปยัง zero value
- ใช้ pointer receiver เมื่อต้องแก้ค่า struct ใหญ่ หรือมี Mutex
- **ข้อควรระวัง:** pointer ที่ชี้เข้าไปใน slice อาจกลายเป็นของเก่าหลัง `append` สร้าง array ใหม่
- ก่อน Go 1.22 ตัวแปรใน loop ถูกใช้ร่วมกันทุกรอบ (`&v` ใน range ได้ address เดิม) ตั้งแต่ 1.22 ได้ตัวแปรใหม่ทุกรอบ

</details>

### 25. Concurrency model ของ Go ทำงานอย่างไร?

<details><summary>ดูคำตอบ</summary>

อิงแนวคิด **CSP (Communicating Sequential Processes)**

- **Goroutine:** `go f()` เป็น "thread" เบา ๆ เริ่มที่ stack ~2KB สร้างได้หลักแสนตัว
- **Channel:** ท่อส่งข้อมูลระหว่าง goroutine แบบ type-safe และ synchronize ให้ในตัว
- **`select`:** รอหลาย channel พร้อมกัน
- **`sync`:** `Mutex`, `WaitGroup`, `Once` และ `sync/atomic` สำหรับกรณีที่ต้องแชร์ memory
- **`context`:** ใช้ส่งสัญญาณ cancel และ deadline

คติ: *"Don't communicate by sharing memory; share memory by communicating."*

ดูฝึกได้ในข้อ 07 ถึง 10

</details>

### 26. จัดการ dependency ในโปรเจกต์อย่างไร?

<details><summary>ดูคำตอบ</summary>

(ต่อจากข้อ 7)

- `go get pkg@version` เลือกเวอร์ชัน, `go get -u ./...` อัปเดตทั้งหมด
- `go list -m all` ดู dependency ทั้งหมด (รวมตัวที่มาทางอ้อม)
- `go mod why pkg` ดูว่าทำไมต้องมี package นี้
- `replace` ใน go.mod ชี้ไปยังโค้ด local หรือ fork
- `go mod vendor` copy dependency มาไว้ใน repo (ใช้ build แบบ offline)
- `GOPRIVATE` ใช้กับ repo ภายในองค์กร
- `govulncheck ./...` ตรวจช่องโหว่

</details>

### 27. Go runtime ทำหน้าที่อะไร?

<details><summary>ดูคำตอบ</summary>

ชั้นซอฟต์แวร์ที่ compile ติดไปใน binary ทุกตัว มีหน้าที่

- **Scheduler:** จัดตาราง goroutine ลงบน OS thread (M:N) ดูข้อ 45
- **Memory allocator + GC**
- **จัดการ stack ของ goroutine:** ขยายและหดได้
- **Channel, select, timer, network poller:** ใช้ epoll หรือ kqueue ทำให้ I/O ดูเหมือน blocking แต่ไม่บล็อก thread
- ตรวจ panic, map race และ deadlock (`all goroutines are asleep`)

ปรับได้ผ่าน package `runtime` และ env เช่น `GOMAXPROCS`, `GOGC`, `GODEBUG`

</details>

### 28. Go เรียกฟังก์ชันอย่างไร (ระดับล่าง)?

<details><summary>ดูคำตอบ</summary>

- ตั้งแต่ Go 1.17 ใช้ **register-based calling convention** (ABIInternal) ส่ง argument และผลลัพธ์ผ่าน register ก่อน ถ้าไม่พอจึงใช้ stack
- ทุก call จะมี stack frame
- ส่วนต้นของฟังก์ชันมี **stack check** ถ้า stack ไม่พอ runtime จะขยาย stack โดย copy ไปยังพื้นที่ที่ใหญ่ขึ้น
- ฟังก์ชันเล็ก ๆ มักถูก **inline** ทำให้ไม่มี overhead ของการเรียก
- `defer` ในปัจจุบันเกือบไม่มี overhead (open-coded defer)

</details>

### 29. Error handling ของ Go ต่างจาก exception อย่างไร?

<details><summary>ดูคำตอบ</summary>

| Go (error value) | Exception (Java/Python) |
|---|---|
| error เป็นค่าธรรมดา คืนจากฟังก์ชัน | throw แล้วกระโดดข้าม stack |
| เห็นทุกจุดที่ fail ได้ใน signature | ไม่รู้ว่าฟังก์ชันไหน throw บ้าง |
| บังคับให้คิดทุกจุด (verbose) | โค้ดสั้นกว่า แต่มี control flow ซ่อนอยู่ |
| ไม่มีต้นทุนแฝง | ต้นทุนเกิดตอน unwind stack |

Go ยังมี `panic`/`recover` แต่ใช้กับ **กรณีร้ายแรงที่ไม่คาดคิด** เท่านั้น ไม่ใช่ flow ปกติ

</details>

### 30. Channel ช่วยให้ goroutine สื่อสารกันอย่างไร?

<details><summary>ดูคำตอบ</summary>

```go
ch := make(chan int)      // unbuffered: ผู้ส่งกับผู้รับต้องมาเจอกัน
bch := make(chan int, 10) // buffered: ส่งได้จนกว่า buffer จะเต็ม
ch <- 1                   // ส่ง
v, ok := <-ch             // รับ (ok=false ถ้า channel ถูก close และหมดข้อมูลแล้ว)
close(ch)                 // ผู้ส่งเป็นคน close
for v := range ch {}      // วนจนกว่าจะถูก close

select {
case v := <-a:
case b <- x:
case <-time.After(time.Second):   // timeout
default:                          // non-blocking
}
```

**กฎที่ต้องจำ:**

- ส่งลง channel ที่ close แล้ว → panic
- close ซ้ำ → panic
- ส่งหรือรับจาก nil channel → บล็อกตลอดไป
- ใช้ `chan<- T` และ `<-chan T` จำกัดทิศทางใน signature

</details>

### 31. Method signature กับ polymorphism

<details><summary>ดูคำตอบ</summary>

Polymorphism ใน Go ทำผ่าน **interface**: ฟังก์ชันรับ interface แล้วส่ง type ใดก็ได้ที่มี method signature ตรง

```go
type Notifier interface { Notify(msg string) error }
func Alert(n Notifier) { n.Notify("server down") }  // Email, Slack, SMS ใช้ได้หมด
```

หลักที่นิยม: **"Accept interfaces, return structs"** และทำ interface ให้เล็ก (1 ถึง 3 method) เช่น `io.Reader` ดูฝึกได้ในข้อ 04

</details>

### 32. ขั้นตอนของ memory allocation และ garbage collection

<details><summary>ดูคำตอบ</summary>

1. **Compile time:** escape analysis ตัดสินว่าจะวางบน stack หรือ heap
2. **Allocate บน heap:** runtime มีแคชต่อ P (`mcache`) แบ่งตาม size class จึงจองเร็วโดยไม่ต้อง lock
3. **GC เริ่ม** เมื่อ heap โตถึงเป้า (ค่า default คือโต 100% จาก live heap ครั้งก่อน ตาม `GOGC=100`) หรือเมื่อใกล้ `GOMEMLIMIT`
4. **Mark (concurrent):** ไล่จาก root (global, stack) แล้วระบาย object เป็นขาว เทา ดำ (tri-color) พร้อมเปิด write barrier เพื่อให้โปรแกรมทำงานต่อได้ระหว่าง mark
5. **Sweep:** object ที่ยังเป็นสีขาวคือไม่มีใครใช้แล้ว จึงคืนพื้นที่ไปใช้ใหม่

มีช่วง stop-the-world สั้น ๆ เพียงเล็กน้อยตอนเริ่มและจบ mark

</details>

### 33. ฟังก์ชันที่คืนหลายค่า

<details><summary>ดูคำตอบ</summary>

```go
func divmod(a, b int) (q, r int) {  // named result
    q, r = a/b, a%b
    return
}
q, r := divmod(7, 2)
_, r = divmod(7, 2)                  // ใช้ _ ทิ้งค่าที่ไม่ต้องการ
v, err := strconv.Atoi("x")         // รูปแบบ (result, error)
v, ok := m["k"]                     // comma-ok idiom
```

ใช้แทน tuple และ exception ได้อย่างเรียบง่าย ส่วน named result มีประโยชน์กับ `defer` ที่ต้องแก้ค่า error ก่อนคืน

</details>

### 34. ทำไม type conversion ถึงสำคัญใน Go?

<details><summary>ดูคำตอบ</summary>

เพราะ Go **ไม่แปลงให้อัตโนมัติ** แม้ `int` กับ `int32` ก็ถือเป็นคนละ type

- ป้องกันบั๊กเงียบ เช่น ข้อมูลหายตอนแปลงโดยไม่รู้ตัว
- ทำให้เห็นชัดว่าจุดไหนเสี่ยง overflow หรือเสียความแม่นยำ
- แม้แต่ `type Celsius float64` ก็ใช้แทน `float64` ตรง ๆ ไม่ได้ ช่วยป้องกันการใช้หน่วยผิด

</details>

### 35. เขียน unit test ด้วย package testing

<details><summary>ดูคำตอบ</summary>

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name       string
        a, b, want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d,%d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

- `t.Errorf` รายงานแล้วรันต่อ ส่วน `t.Fatalf` หยุด test นั้นทันที
- คำสั่งที่ใช้บ่อย: `go test -v -run TestAdd -race -cover ./...`

ดูฝึกได้ในข้อ 12

</details>

### 36. GC ส่งผลต่อ memory อย่างไร?

<details><summary>ดูคำตอบ</summary>

- **ข้อดี:** memory usage คงที่ ไม่ต้อง free เอง และไม่มี dangling pointer
- **ต้นทุน:** ใช้ CPU (มักอยู่ราว 25% ช่วงที่ GC ทำงาน) มี pause สั้น ๆ และต้องใช้ memory เผื่อ (heap ใหญ่กว่า live data ประมาณ 2 เท่าที่ `GOGC=100`)
- **ลดภาระ GC** ได้ด้วยการลดจำนวน allocation เช่น preallocate, `sync.Pool`, ใช้ value แทน pointer และหลีกเลี่ยง `[]*T` ขนาดใหญ่
- ดูสถิติได้ด้วย `GODEBUG=gctrace=1` หรือ `runtime.ReadMemStats`

</details>

---

## ระดับสูง (Advanced)

### 37. Go จัดการ memory leak และหาคอขวดอย่างไร?

<details><summary>ดูคำตอบ</summary>

**เครื่องมือ**

- `import _ "net/http/pprof"` แล้วดู `/debug/pprof/heap`, `/goroutine`, `/profile` (CPU), `/block`, `/mutex`
- `go tool pprof -http=:8080 heap.out` ดู flame graph
- `go test -bench . -benchmem -memprofile mem.out -cpuprofile cpu.out`
- `go tool trace` ดู timeline ของ goroutine และ GC
- metric: `runtime.NumGoroutine()`, heap in-use

**สาเหตุ leak ที่พบบ่อย:** goroutine ค้าง, ไม่ `Close()` body ของ HTTP response, ไม่ `Stop()` ticker (ก่อน Go 1.23), cache ไม่มีขอบเขต และ sub-slice ถือ array ใหญ่

**วิธีดู:** เทียบ heap profile 2 ช่วงเวลา (`-base`) ว่าอะไรโตขึ้น

</details>

### 38. Dependency management ของ Go เทียบกับภาษาอื่น

<details><summary>ดูคำตอบ</summary>

- **go.mod อ่านง่าย** เทียบกับ `pom.xml` (Maven) หรือ `package.json` + lock file (npm)
- **Minimal Version Selection (MVS):** Go เลือก "เวอร์ชันต่ำสุดที่ตรงทุกเงื่อนไข" ผลลัพธ์จึง deterministic โดยไม่ต้องมี lock file แยก ต่างจาก npm ที่เลือกเวอร์ชันล่าสุดใน range
- **Import path คือ URL ของ repo** ไม่ต้องมี central registry (แต่มี proxy.golang.org และ checksum DB คอยช่วย)
- **Semantic import versioning:** v2 ขึ้นไปต้องใส่ `/v2` ใน import path จึงใช้ v1 กับ v2 พร้อมกันได้
- ไม่มี `node_modules` ซ้อนกันหลายชั้น

</details>

### 39. ฟังก์ชันโต้ตอบกับฟังก์ชันรอบข้างอย่างไร?

<details><summary>ดูคำตอบ</summary>

- **Stack frame:** ผู้เรียกเตรียม argument แล้ว call ผู้ถูกเรียกมี frame ของตัวเอง เมื่อ return จะคืนค่าให้ผู้เรียก
- **Closure:** ฟังก์ชันภายในจับตัวแปรภายนอกแบบ **by reference** (ตัวแปรนั้นจะ escape ไป heap)

```go
func counter() func() int {
    n := 0
    return func() int { n++; return n }
}
```

- **`defer`** ทำงานแบบ LIFO ตอนฟังก์ชันจบ (แม้จะ panic) โดยประเมินค่า argument ตั้งแต่บรรทัดที่เขียน `defer`
- **panic** จะ unwind stack ไปทีละ frame และรัน defer ของแต่ละ frame จนกว่าจะเจอ `recover`

</details>

### 40. จัดการ concurrency ให้ใช้ memory อย่างมีประสิทธิภาพ

<details><summary>ดูคำตอบ</summary>

- goroutine ถูกก็จริง แต่**ไม่ฟรี** ถ้าสร้าง 1 ล้านตัว จะใช้ stack อย่างน้อย ~2GB ควร**จำกัดจำนวน**ด้วย worker pool (ข้อ 07), semaphore (`chan struct{}` แบบ buffered) หรือ `errgroup.SetLimit`
- ทุก goroutine ต้อง**มีทางจบ** โดยส่ง context เข้าไป (ข้อ 09, 10)
- ตั้งขนาด buffer ของ channel ให้เหมาะสม ไม่ใหญ่เกินจำเป็น
- ใช้ `sync.Pool` reuse buffer ลด allocation
- ลด lock contention ด้วยการแบ่ง shard หรือใช้ atomic

</details>

### 41. การแปลงระหว่าง data type ต่าง ๆ

<details><summary>ดูคำตอบ</summary>

```go
var big int64 = 300
small := int8(big)        // 44 — overflow แบบเงียบ ๆ ไม่มี error!
f := 3.99
i := int(f)               // 3 — ตัดทศนิยมทิ้ง ไม่ปัด
n := -1
u := uint(n)              // 18446744073709551615 — wrap รอบ
```

- ถ้าต้องการความปลอดภัย ให้ตรวจช่วงก่อนแปลง หรือใช้ `strconv.ParseInt(s, 10, 8)` ที่คืน error เมื่อเกินช่วง
- `string` ↔ `[]byte` จะ copy ข้อมูล (compiler ปรับให้ไม่ copy ได้ในบางกรณี)
- interface → type ใช้ **type assertion** `v, ok := x.(T)` (ถ้าไม่ใช้ `ok` แล้ว type ไม่ตรงจะ panic)

</details>

### 42. Method signature กับการเขียนโค้ดที่ reuse ได้

<details><summary>ดูคำตอบ</summary>

signature ที่ดีคือ **สัญญา** ที่ชัดเจน

- รับ **interface เล็ก ๆ** เช่น `io.Reader` ทำให้ฟังก์ชันใช้ได้กับไฟล์, network, buffer หรือ gzip
- รับ `context.Context` เป็นตัวแรก และคืน `error` เป็นตัวสุดท้าย
- หลีกเลี่ยง parameter แบบ `bool` หลายตัว ใช้ struct ของ option หรือ functional options แทน
- Generics ช่วย reuse algorithm ข้าม type ได้ เช่น `slices.Contains`, `maps.Keys`

</details>

### 43. Pointer และ memory ในระดับล่าง

<details><summary>ดูคำตอบ</summary>

- pointer ใน Go เป็น address ขนาด 8 byte (บน 64-bit) ที่ **GC ติดตามได้** runtime จึงย้าย stack ได้อย่างปลอดภัยเพราะรู้ว่า pointer ตัวไหนชี้เข้าไปใน stack
- ไม่มี pointer arithmetic (ยกเว้น `unsafe.Pointer` + `unsafe.Add` ซึ่งเสี่ยง)
- **Struct layout และ alignment:** เรียง field จากใหญ่ไปเล็กช่วยลด padding

```go
type Bad  struct { a bool; b int64; c bool }  // 24 byte
type Good struct { b int64; a, c bool }       // 16 byte
```

- ใช้ pointer มากเกินไป เช่น `[]*T` จะเพิ่มงานของ GC และทำให้ cache locality แย่กว่า `[]T`

</details>

### 44. ออกแบบ custom error type

<details><summary>ดูคำตอบ</summary>

```go
type HTTPError struct {
    Code int
    Err  error
}
func (e *HTTPError) Error() string { return fmt.Sprintf("http %d: %v", e.Code, e.Err) }
func (e *HTTPError) Unwrap() error { return e.Err }   // ให้ errors.Is/As มองทะลุเข้าไปได้

var he *HTTPError
if errors.As(err, &he) && he.Code == 404 { ... }
```

- เพิ่ม `Unwrap()` เพื่อรองรับ error chain
- method `Is(target error) bool` ใช้กำหนดการเทียบเองได้
- ระวัง **typed nil**: ห้าม return `(*HTTPError)(nil)` เป็น `error`

ดูฝึกได้ในข้อ 05

</details>

### 45. Go runtime จัดตาราง goroutine อย่างไร?

<details><summary>ดูคำตอบ</summary>

ใช้โมเดล **G-M-P** (M:N scheduling)

- **G** = goroutine
- **M** = OS thread (machine)
- **P** = processor เชิงตรรกะ มีจำนวนเท่ากับ `GOMAXPROCS` (ค่า default คือจำนวน CPU) แต่ละ P มี **local run queue** ของตัวเอง

การทำงาน:

1. M ต้องถือ P จึงจะรัน G ได้ โดยหยิบ G จาก local queue ของ P ก่อน
2. ถ้า queue ว่าง จะ **work stealing** ขโมย G ครึ่งหนึ่งจาก P อื่น หรือหยิบจาก global queue
3. ถ้า G ทำ syscall ที่บล็อก M จะปล่อย P ให้ M ตัวอื่นรัน G ต่อได้ ส่วน network I/O ใช้ netpoller ทำให้ไม่บล็อก thread
4. **Preemption:** ตั้งแต่ Go 1.14 เป็นแบบ asynchronous (ใช้ signal) goroutine ที่วนลูปนาน ๆ จึงไม่ทำให้ตัวอื่นอดรัน

</details>

### 46. การประกาศตัวแปรกับประสิทธิภาพของ memory

<details><summary>ดูคำตอบ</summary>

- **จำกัด scope** ให้แคบที่สุด ตัวแปรจะได้ถูกคืนเร็ว และมีโอกาสอยู่บน stack
- **Preallocate:** `make([]T, 0, n)`, `make(map[K]V, n)`, `strings.Builder.Grow(n)` ดูข้อ 13
- **ปล่อย reference** ที่ไม่ใช้แล้ว เช่น ตั้งเป็น `nil` ตอนลบ element ออกจาก slice ของ pointer หรือใช้ `clear()`
- หลีกเลี่ยงการเก็บ sub-slice หรือ substring ของข้อมูลใหญ่ (ใช้ `strings.Clone` หรือ `bytes.Clone`)
- ใช้ value type แทน pointer เมื่อ struct เล็ก จะช่วยลดงานของ GC

</details>

### 47. Dependency management กับ version control

<details><summary>ดูคำตอบ</summary>

- เวอร์ชันของ module คือ **git tag** ตาม semver (`v1.2.3`) ถ้าไม่มี tag จะใช้ **pseudo-version** (`v0.0.0-20240101120000-abcdef123456`)
- `go.sum` เก็บ **cryptographic hash** ของทุก module ไว้ตรวจกับ checksum database (sum.golang.org) ป้องกัน supply-chain attack
- เปลี่ยน major version (v2+) ต้องเปลี่ยน import path
- `retract` ใน go.mod ใช้ประกาศว่าเวอร์ชันไหนห้ามใช้
- ใน CI ควรรัน `go mod verify` และ `go mod tidy -diff` เพื่อตรวจว่า go.mod และ go.sum ตรงกับโค้ด

</details>

### 48. Go จัดการ global variable อย่างไร และมีผลอย่างไร?

<details><summary>ดูคำตอบ</summary>

- initialize ตามลำดับ dependency ก่อน `init()` และ `main()` และอยู่ใน memory ตลอดอายุโปรแกรม เป็น root ของ GC สิ่งที่ global ชี้อยู่จึงไม่มีวันถูกเก็บ
- **เสี่ยง data race** ถ้าหลาย goroutine เขียน ต้องป้องกันด้วย `sync.Mutex`, `atomic` หรือ `sync.Once` (สำหรับ lazy init)
- ทำให้ test ขนานกันไม่ได้ และ package ผูกกันแน่น
- ถ้าจำเป็นต้องใช้ ให้ unexport แล้วเข้าถึงผ่านฟังก์ชัน

</details>

### 49. Pointer ไปยัง struct

<details><summary>ดูคำตอบ</summary>

```go
u := &User{Name: "Tah"}   // composite literal + & (นิยมใช้)
u.Name = "Bee"            // auto-dereference
func Update(u *User) { u.Age++ }   // แก้ตัวจริง ไม่ copy
```

- เหมาะกับ struct ใหญ่ หรือ struct ที่ต้องแก้ค่า หรือมี Mutex
- **ข้อควรระวัง:** nil pointer dereference, หลาย goroutine ถือ pointer เดียวกันแล้วเกิด race และ pointer ที่ชี้เข้าไปใน slice (`&s[i]`) อาจชี้ข้อมูลเก่าหลัง `append`
- struct เล็ก เช่น `Point{X,Y}` ส่งเป็น value มักเร็วกว่า เพราะไม่ escape และไม่เพิ่มงานให้ GC

</details>

### 50. วิเคราะห์โค้ดเพื่อหาปัญหา dependency และ memory

<details><summary>ดูคำตอบ</summary>

**Dependency**

- `go mod graph` และ `go mod why` ดูว่าอะไรดึงอะไรเข้ามา
- `govulncheck` ตรวจช่องโหว่
- pin เวอร์ชันและ commit `go.sum`
- ลด dependency ที่ไม่จำเป็น (standard library มักพอแล้ว)

**Memory** ตรวจโค้ดหา

- resource ที่ไม่ได้ `Close()` (file, `resp.Body`, rows)
- goroutine ที่ไม่มีทางจบ
- global map หรือ cache ที่ไม่มี limit
- `time.After` ในลูปที่ทำงานนาน ๆ (ก่อน Go 1.23 timer จะค้างจนกว่าจะหมดเวลา)
- `defer` ภายในลูป (ทำงานตอนฟังก์ชันจบ ไม่ใช่ตอนจบแต่ละรอบ)

แล้ว**วัดจริง**ด้วย pprof, benchmark และ `-race` อย่าเดา

</details>
