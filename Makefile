# Build dependencies
GO = tinygo
WASM_OPT = wasm-opt
W4 = w4

# Whether to build for debugging instead of release
DEBUG = 0

# Compilation flags
GOFLAGS = -target ./target.json -panic trap
ifeq ($(DEBUG), 1)
	GOFLAGS += -opt 1
else
	GOFLAGS += -opt z -no-debug
endif

# wasm-opt flags
WASM_OPT_FLAGS = -Oz --zero-filled-memory --strip-producers --enable-bulk-memory

# Output
OUTPUT = build/jumpnshoot.wasm

all: $(OUTPUT)

$(OUTPUT):
	@mkdir -p build
	$(GO) build $(GOFLAGS) -o $@ ./src
ifneq ($(DEBUG), 1)
ifeq (, $(shell command -v $(WASM_OPT)))
	@echo Tip: $(WASM_OPT) was not found. Install it from binaryen for smaller builds!
else
	$(WASM_OPT) $(WASM_OPT_FLAGS) $@ -o $@
endif
endif

run: $(OUTPUT)
	$(W4) run $<

watch: $(OUTPUT)
	$(W4) watch $<

bundle: $(OUTPUT)
	$(W4) bundle $< --title "Jump-N-Shoot" --html jumpnshoot.html

debug:
	$(MAKE) DEBUG=1

.PHONY: all clean run watch bundle debug

clean:
	rm -rf build