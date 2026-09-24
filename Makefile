# ใช้: make test N=01
N ?= 01

.PHONY: test solution reset test-all solution-all html

test: ## รัน test ของโจทย์ข้อ N
	go test -race -v ./exercises/$(N)_*/

solution: ## รัน test ของเฉลยข้อ N
	go test -race -v ./solutions/$(N)_*/

reset: ## ล้างโค้ดข้อ N กลับเป็นโจทย์เริ่มต้น
	git checkout -- exercises/$(N)_*/

test-all:
	go test -race ./exercises/...

solution-all:
	go test -race ./solutions/...

html: ## สร้าง golang-interview.html จาก PROBLEMS.md, SOLUTIONS.md, THEORY.md
	python3 tools/build_html.py
