# Go Worker Pool + Channels — Revision Notes

## 1. Goroutine kya hai?

Goroutine = ek lightweight thread jo Go khud manage karta hai.

```go
go someFunction()
```

Bas `go` keyword laga do, function alag se **parallel** chalne lagta hai, aur `main()` aage badh jata hai (wait nahi karta).

**Problem:** Agar `main()` khatam ho gaya, to sab goroutines mar jaate hain — chahe unka kaam adhura ho.

Isiliye humein ek **synchronization mechanism** chahiye — ye kaam karta hai `sync.WaitGroup` aur **Channels**.

---

## 2. Channel kya hai?

Channel = ek **pipe** jisse goroutines aapas mein data bhej/le sakte hain, safely (bina race condition ke).

```go
ch := make(chan string)
```

Socho ek **conveyor belt**:
- Ek taraf se cheez daalo → `ch <- "server1.log"`
- Doosri taraf se cheez uthao → `file := <-ch`

**Key property:** Channel **blocking** hoti hai by default.
- Agar koi cheez daal raha hai aur koi le nahi raha → daalne wala **ruk jayega** (wait karega).
- Agar koi cheez lena chahta hai aur kuch hai nahi → lene wala **ruk jayega**.

Ye blocking hi worker pool ko naturally synchronize karti hai — koi manual locking nahi chahiye.

---

## 3. Buffered vs Unbuffered Channel

```go
ch := make(chan string)      // unbuffered — capacity 0
ch := make(chan string, 10)  // buffered — capacity 10
```

- **Unbuffered:** Sender tabhi aage badhega jab receiver ready ho (handshake jaisa).
- **Buffered:** Sender 10 items daal sakta hai bina kisi receiver ke wait kiye. 11th item par ruk jayega.

Worker pool mein aam taur pe **buffered channel** use hota hai jobs ke liye, taaki scanner rukta na rahe.

---

## 4. Worker Pool Pattern — Poora Mental Model

### Ingredients:
1. **Jobs Channel** — saari "to-do" files/tasks isme daali jaati hain.
2. **N Workers** — goroutines jo jobs channel se padhte rehte hain.
3. **Results Channel** — har worker apna output yahan bhejta hai.
4. **Main goroutine** — sab results collect/merge karta hai.

### Flow Diagram:

```
Files (scanner)
      ↓
  Jobs Channel  ──────┬────────┬────────┬────────┐
                       ↓        ↓        ↓        ↓
                   Worker1  Worker2  Worker3  Worker4
                       ↓        ↓        ↓        ↓
                       └────────┴────────┴────────┘
                                 ↓
                         Results Channel
                                 ↓
                         Main (merge results)
```

### Worker ka code-shape (concept, na ki final code):

```go
func worker(jobs <-chan string, results chan<- Result) {
    for file := range jobs {
        // ye line hi worker pool ka dil hai:
        // jab tak jobs channel se kuch milta rahega, loop chalta rahega
        // jab channel CLOSE ho jayega, loop khud khatam ho jayega
        report := ParseSingleFile(file)
        results <- report
    }
}
```

**Important:** `for job := range jobs` — ye tab tak chalega jab tak jobs channel **close** nahi hoti. Isliye scanner ko saari files daal ke channel **close karna zaroori hai**, warna workers hamesha wait karte rahenge (deadlock).

---

## 5. `<-chan` aur `chan<-` ka matlab?

- `<-chan string` → **receive-only** channel (sirf padh sakte ho)
- `chan<- string` → **send-only** channel (sirf likh sakte ho)
- `chan string` → dono kar sakte ho

Ye Go ka type-safety feature hai — worker function ko explicitly bolna "tu sirf jobs padhega, results sirf likhega" — taaki accidental galat use na ho.

---

## 6. `sync.WaitGroup` kyun chahiye?

Channels data transfer karte hain, lekin **"sab workers apna kaam khatam kar chuke"** ye track karne ke liye WaitGroup use hota hai.

```go
var wg sync.WaitGroup

for i := 0; i < numWorkers; i++ {
    wg.Add(1)           // "ek aur kaam pending hai"
    go func() {
        defer wg.Done()  // "mera kaam khatam"
        worker(jobs, results)
    }()
}

wg.Wait()  // jab tak sab Done() na bol dein, yahan ruko
close(results)  // sab workers khatam → ab results channel close karo
```

**Analogy:** WaitGroup ek counter hai. `Add(1)` = counter ++. `Done()` = counter --. `Wait()` = jab tak counter 0 na ho, ruko.

---

## 7. Poora Sequence — Step by Step

1. `jobs := make(chan string, 100)` → jobs channel banao
2. `results := make(chan Report, 100)` → results channel banao
3. `N` workers start karo (goroutines), sab `jobs` se padhna shuru kar dete hain (abhi channel khali hai to wait karenge)
4. Main goroutine `ReadDir()` karke saari filenames `jobs` channel mein daalta hai
5. Saari files daalne ke baad `close(jobs)` — ye workers ko signal deta hai "aur kuch nahi aayega"
6. Har worker apna `for range jobs` loop khatam karta hai jab channel empty + closed ho
7. `wg.Wait()` — main yahan ruk jata hai jab tak saare workers na khatam ho jayein
8. Sab khatam → `close(results)`
9. Ab main `for r := range results` karke sab results collect/merge karta hai → Final Report

---

## 8. Deadlock — Sabse Common Mistake

Agar tumne `close(jobs)` nahi kiya, to workers hamesha `for range jobs` mein wait karte rahenge, aur `wg.Wait()` kabhi khatam nahi hoga → **program hang ho jayega**.

**Golden rule:** Jo bhi channel mein data daal raha hai, wahi close karega. Receiver kabhi close nahi karta.

---

## 9. Tumhare Project Mapping (Log Analyzer)

| Current (Sequential) | Naya (Worker Pool) |
|---|---|
| `ReadDir()` → filenames list | Same |
| `for _, file := range files { ParseSingleFile(file) }` | Filenames → `jobs` channel mein daalo |
| — | N workers goroutines, jo `jobs` se padh ke `ParseSingleFile()` call karte hain |
| Direct `LogReport` return | Har worker apna partial report `results` channel mein bhejta hai |
| — | Main `results` collect karke `MergeReports()` se final `LogReport` banata hai |
| Writer/Database/Analytics | **Bilkul same rahega, kuch nahi badlega** |

---

## 10. Ek-line Summary (Yaad Rakhne Ke Liye)

> **Jobs channel ek queue hai jisme se free worker khud utha leta hai. Results channel wo pipe hai jisse har worker apna output wapas bhejta hai. WaitGroup batata hai sab workers kab khatam hue. Jo daalta hai wahi close karta hai.**

---

## Next Step (Jab Ready Ho)

Commit order jo tumne khud socha tha, wahi sahi hai:
1. `worker.go` mein `Job` + `Result` struct
2. Jobs channel banao (empty, no logic)
3. Ek worker jo sirf filename print kare (no parsing)
4. Worker mein `ParseSingleFile()` connect karo
5. N workers (goroutine loop + WaitGroup)
6. `MergeReports()` likho
7. Purane sequential loop ko hata do, naya flow connect karo

---
---

# PART 2 — Deep Dive (Interview + Project Level)

---

## 11. Goroutine vs OS Thread (Interview Favorite)

| Goroutine | OS Thread |
|---|---|
| Go runtime manage karta hai (not OS) | OS kernel manage karta hai |
| Starting stack size ~2KB (grows dynamically) | Fixed, usually 1-2MB |
| Lakhs goroutines chal sakti hain easily | Thousands hi practically possible |
| Context switch cheap hai (user-space) | Context switch costly hai (kernel-space) |
| M:N scheduling (M goroutines on N OS threads) | 1:1 with CPU |

**Interview line:** *"Goroutines are multiplexed onto a small number of OS threads by Go's scheduler using M:N scheduling — that's why they're so cheap."*

### GOMAXPROCS
```go
runtime.GOMAXPROCS(4) // kitne OS threads parallel chalenge
```
Default = number of CPU cores. Ye batata hai kitni goroutines **truly parallel** (not just concurrent) chal sakti hain.

**Concurrency vs Parallelism (bahut common interview Q):**
- **Concurrency** = multiple cheezein "manage" karna (structure) — ek CPU pe bhi ho sakta hai (interleaved).
- **Parallelism** = multiple cheezein **actually same time** chalna — multiple cores chahiye.

> *"Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once."* — Rob Pike

---

## 12. Channel Ke Saare Rules (Ratta Maar Lo)

| Operation | nil channel | Open channel | Closed channel |
|---|---|---|---|
| **Send** (`ch <- v`) | Blocks forever | Works (blocks if full/no receiver) | **panic!** |
| **Receive** (`<-ch`) | Blocks forever | Works | Returns zero-value immediately, `ok=false` |
| **Close** | panic! | Works | **panic!** (double close) |

### Comma-ok idiom (channel closed check karne ke liye):
```go
val, ok := <-ch
if !ok {
    // channel closed hai aur empty bhi
}
```

### Gotchas jo interview mein poochte hain:
1. **"Send on closed channel panics"** — ye sabse common panic hai worker pools mein.
2. **"Close of closed channel panics"** — kabhi bhi channel ko do baar close mat karo.
3. **Sirf sender close kare, receiver kabhi nahi** — golden rule.
4. **nil channel forever blocks** — kabhi kabhi jaanbujh ke `select` mein use hota hai (branch disable karne ke liye).

---

## 13. `select` Statement — Channels ka `switch`

Jab tumhe **multiple channels** ek saath handle karni ho:

```go
select {
case job := <-jobs:
    process(job)
case <-time.After(5 * time.Second):
    fmt.Println("timeout!")
case <-ctx.Done():
    fmt.Println("cancelled!")
default:
    fmt.Println("kuch ready nahi, aage badho")
}
```

Rules:
- Jo bhi case **pehle ready** ho jaye, wahi chalega.
- Agar multiple ready hain ek saath, **random** pick hota hai.
- `default` case ho to select kabhi block nahi karega (non-blocking check).

**Use case in your project:** Agar tum chahte ho parser 5 second se zyada ek file pe atka na rahe, to `select` + `time.After()` use karoge (timeout pattern).

---

## 14. sync Package — Channels Ke Alternatives/Helpers

### `sync.WaitGroup` — already discussed
Wait karne ke liye ki N goroutines khatam hui ya nahi.

### `sync.Mutex` — Shared memory protect karne ke liye
```go
var mu sync.Mutex
var counter int

mu.Lock()
counter++
mu.Unlock()
```

**Channel vs Mutex — bahut popular interview question:**

| Use Channel when... | Use Mutex when... |
|---|---|
| Data/ownership **transfer** karna hai goroutines ke beech | Sirf shared state **protect** karna hai (no transfer) |
| Pipeline/worker pool jaisa flow hai | Simple counter, cache, map jaisi cheez hai |
| "Communicating" | "Protecting" |

> Go proverb: **"Don't communicate by sharing memory; share memory by communicating."**

Tumhare project mein: agar sirf ek shared `errorCount int` badhana hota to Mutex kaafi tha. Lekin poora Report object transfer karna hai worker se main tak — isliye **channel sahi choice hai**.

### `sync.RWMutex`
Jab **reads zyada** hon, writes kam:
```go
var mu sync.RWMutex
mu.RLock()   // multiple readers ek saath allowed
mu.RUnlock()
mu.Lock()    // writer ko exclusive access
mu.Unlock()
```

### `sync.Once`
Kisi cheez ko **exactly ek baar** run karna ho (jaise singleton init):
```go
var once sync.Once
once.Do(func() {
    fmt.Println("sirf ek baar chalega chahe 100 goroutines se call ho")
})
```

---

## 15. Race Condition — Kya Hai Aur Kaise Pakdo

**Race condition** = jab 2+ goroutines **bina synchronization** ke same variable read/write karti hain.

```go
var counter int
for i := 0; i < 1000; i++ {
    go func() { counter++ }()  // RACE! counter++ atomic nahi hai
}
```

`counter++` actually 3 steps hai: read → increment → write. Do goroutines beech mein overlap ho sakti hain → wrong result.

**Detect karne ka tool:**
```bash
go run -race main.go
go test -race ./...
```

**Interview line:** *"go run -race uses the race detector which instruments memory accesses to detect concurrent unsynchronized access at runtime."*

**Fix options:** Mutex, Channel, ya `sync/atomic` package (`atomic.AddInt64` jaisa lightweight atomic ops ke liye).

---

## 16. Common Concurrency Patterns (Interview + Real Use)

### A) Worker Pool (already covered)
Fixed N workers, jobs channel se consume karte hain.

### B) Fan-Out, Fan-In
- **Fan-out**: ek channel se multiple goroutines padh rahi hain (jaise tumhare 4 workers ek jobs channel se).
- **Fan-in**: multiple channels ka output ek channel mein merge karna.

```go
func fanIn(ch1, ch2 <-chan int) <-chan int {
    merged := make(chan int)
    go func() {
        for v := range ch1 { merged <- v }
    }()
    go func() {
        for v := range ch2 { merged <- v }
    }()
    return merged
}
```

Tumhara worker pool actually **fan-out (jobs distribute) + fan-in (results collect)** dono hai!

### C) Pipeline Pattern
Ek stage ka output doosre stage ka input:
```
ReadDir() → jobs chan → Parse (stage1) → parsed chan → Aggregate (stage2) → results chan
```
Har stage apna goroutine + channel rakhta hai. Tumhara pura log analyzer isi pattern pe based hai.

### D) Context for Cancellation/Timeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

select {
case result := <-resultChan:
    // process
case <-ctx.Done():
    fmt.Println("cancelled:", ctx.Err())
}
```

`context.Context` Go mein **standard way** hai cancellation/timeout/deadline propagate karne ka across goroutines. Production code mein har long-running operation `ctx` leta hai as first param — convention hai.

---

## 17. Buffered Channel Ka Sizing — Practical Tip

```go
jobs := make(chan string, len(files))     // saari files ek saath fit ho jayein
results := make(chan Report, numWorkers)  // workers jitna buffer
```

Agar buffer chota hai to producer/consumer ek dusre ko block karte rahenge (which is fine functionally, thoda slower ho sakta hai). Interview mein pooch sakte hain: **"buffered channel ka size kaise decide karoge?"** — answer: workload aur memory tradeoff pe depend karta hai, generally producer/consumer count ke around.

---

## 18. Deadlock Scenarios (Sabse Zyada Poocha Jaane Wala Topic)

```go
// Scenario 1: Unbuffered channel, koi receiver nahi
ch := make(chan int)
ch <- 5  // fatal error: all goroutines are asleep - deadlock!

// Scenario 2: WaitGroup count mismatch
wg.Add(2)
go func() { wg.Done() }()  // sirf 1 Done() call hua, 2 chahiye
wg.Wait()  // hamesha ke liye ruk jayega

// Scenario 3: Mutex double lock (same goroutine)
mu.Lock()
mu.Lock()  // deadlock — khud apne lock ka wait kar raha hai
```

**Go runtime khud detect karta hai** agar **saari** goroutines deadlock ho jaayein (`fatal error: all goroutines are asleep`), lekin agar sirf **ek** goroutine stuck hai baaki chal rahi hain, to koi error nahi aayega — silently hang ho jayega. Ye debugging mushkil banata hai.

---

## 19. Interview Q&A — Rapid Fire

**Q: Channel aur Mutex mein kab kya use karoge?**
A: Data transfer/ownership pass karna ho → Channel. Sirf shared state protect karna ho → Mutex.

**Q: Unbuffered channel synchronous hai ya asynchronous?**
A: Synchronous — sender tab tak block rehta hai jab tak receiver ready na ho (rendezvous point).

**Q: `nil` channel pe receive karoge to kya hoga?**
A: Forever block ho jayega (deadlock, agar koi aur activity na ho).

**Q: Goroutine leak kya hota hai?**
A: Jab ek goroutine kabhi khatam nahi hoti (jaise channel pe block hui hai jisme kabhi data ya close nahi aata) — memory leak jaisa hi hai, permanently resources hold karta hai.

**Q: WaitGroup ko value se pass karoge ya pointer se?**
A: Hamesha **pointer** (`*sync.WaitGroup` ya struct ke andar embed) — value copy hone se `Add`/`Done` alag copies pe operate karenge, bug ban jayega.

**Q: `for range` channel pe kab tak chalta hai?**
A: Jab tak channel **close** na ho jaye aur usme koi pending value na ho.

**Q: Worker pool ka size kitna rakhna chahiye?**
A: CPU-bound kaam ke liye ~`runtime.NumCPU()`. I/O-bound (file read, network) ke liye usse zyada bhi chal sakta hai kyunki workers waiting mein CPU free karte rehte hain.

**Q: Panic ek goroutine mein aaye to kya poore program pe asar padega?**
A: Haan — agar recover nahi kiya gaya to poora program crash ho jata hai, sirf wo goroutine nahi marti. Isliye production worker pools mein har worker ke andar `defer recover()` lagate hain.

---

## 20. Production-Grade Worker Pool Checklist (Interview Mein Ye Bologe To Impress Hoga)

- [ ] Jobs channel **buffered** (producer block na ho unnecessarily)
- [ ] `sync.WaitGroup` se worker completion track
- [ ] Sender hi channel close kare (never receiver)
- [ ] Har worker ke andar `defer func() { recover() }()` — ek file corrupt/panic ho to poora pool na gire
- [ ] `context.Context` for cancellation/timeout support
- [ ] Race detector se test kiya (`go test -race`)
- [ ] Worker count configurable (`runtime.NumCPU()` based default)
- [ ] Results merge karte waqt bhi agar shared state hai to Mutex ya single-consumer pattern use karo

---

## 21. One-Page Cheat Sheet

```
GOROUTINE     → go func()                     lightweight thread
CHANNEL       → make(chan T) / make(chan T,n)  pipe between goroutines
SEND          → ch <- val
RECEIVE       → val := <-ch
CLOSE         → close(ch)                      only sender closes
RANGE         → for v := range ch              auto-stops on close+empty
SELECT        → multiplex multiple channels
WAITGROUP     → wg.Add/Done/Wait               wait for N goroutines
MUTEX         → mu.Lock/Unlock                 protect shared state
CONTEXT       → ctx.Done()/cancel()            cancellation/timeout
RACE DETECT   → go run -race                   find unsynced access
```

---

Jab yahan tak clear ho jaye aur tum khud implement karna chaho, bol dena — main sirf **review aur guide** karunga, code tum likhoge. 💪