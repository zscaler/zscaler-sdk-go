#!/usr/bin/env python3
"""
coverage_summary.py
===================

Reads a Go cover profile (atomic / count / set mode) produced by `go test
-coverprofile=...` and prints a per-product summary table for the Zscaler
Go SDK.

Usage
-----

    # Aggregate (all products) — what `make test:unit:coverage` shows.
    python3 scripts/coverage_summary.py [path/to/coverage.out]

    # Per-service breakdown of a single product — what each
    # `make test:unit:<product>` target shows.
    python3 scripts/coverage_summary.py \
        --product zpa unit-zpa-coverage.out

Defaults to `unit-coverage.out` when no path is given.

Coverage profile format
-----------------------

Each non-header line has the shape

    <import-path>/<file>.go:<startLine>.<col>,<endLine>.<col> <stmts> <count>

Aggregate mode buckets lines by the first path segment under `zscaler/`
— e.g. `zscaler/zpa/...` → product `zpa`. Anything directly under
`zscaler/` (oneapiclient.go, zparequests.go, etc.) is bucketed as `core`
because it is shared infrastructure exercised by every product.

Per-product mode (--product <name>) buckets lines by the first segment
under `zscaler/<name>/services/` — e.g. `services/appconnectorgroup/...`
→ service `appconnectorgroup`. Files directly under `zscaler/<name>/`
(v2_client.go, v2_config.go, etc.) are bucketed as `(root)` and are
excluded from the per-service table and TOTAL — they are legacy client
infrastructure, not individual service packages under `services/`.

Both modes produce the same column layout so the per-product table that
follows `make test:unit:zpa` is visually consistent with the aggregate
table from `make test:unit:coverage`.

Columns reported:
  - Coverage %   : covered / total statements
  - Covered/Total: the raw counts behind the percentage
  - Files        : number of unique SDK files that contributed statements
  - Notes        : aggregate mode only — family/role annotation
"""

from __future__ import annotations

import os
import re
import sys
from collections import defaultdict

# Match the SDK source path layout. Group 1 captures the top-level slot
# directly under `zscaler/` (a product name like `zpa`, `zia`, … or a
# bare filename like `oneapiclient.go` when the path is shared core).
PRODUCT_RE = re.compile(
    r"github\.com/zscaler/zscaler-sdk-go/v3/zscaler/([^/:]+)"
)

# Family / role annotations shown in the rightmost column. Keeps the
# table self-documenting for someone reviewing it for the first time.
NOTES = {
    "core": "shared transport / routing (oneapiclient, *requests.go, cache, …)",
    "zpa": "OneAPI — Zero Trust Private Access",
    "zia": "OneAPI — Internet Access",
    "zcc": "OneAPI — Client Connector",
    "zdx": "OneAPI — Digital Experience",
    "ztw": "OneAPI — Cloud / Branch Connector",
    "zid": "OneAPI — Zidentity",
    "zwa": "Standalone — Workflow Automation (not OneAPI-routed)",
    "common": "shared zscaler/common helpers",
    "errorx": "error type helpers",
}

# ANSI colors. Disabled automatically when stdout is not a TTY (CI logs).
USE_COLOR = sys.stdout.isatty() and os.environ.get("NO_COLOR") is None
def _c(code: str) -> str:
    return code if USE_COLOR else ""

RESET = _c("\033[0m")
BOLD = _c("\033[1m")
DIM = _c("\033[2m")
CYAN = _c("\033[36m")
GREEN = _c("\033[32m")
YELLOW = _c("\033[33m")
RED = _c("\033[31m")


def pct_color(pct: float) -> str:
    """Color-code a coverage percentage for quick visual scanning."""
    if pct >= 50:
        return GREEN
    if pct >= 25:
        return YELLOW
    return RED


def parse_profile(path: str):
    """Return (per_product_total, per_product_covered, per_product_files).

    Block deduplication
    -------------------

    When `-coverpkg=...` instruments packages outside the package under
    test, every test binary in the run emits its own copy of each block
    in the merged profile (the file just appends entries; the merge
    happens later inside `go tool cover`). For us that means the same
    block — say `jmespath.go:21.82,23.2 1` — typically shows up once per
    test binary in `./tests/unit/...`, with possibly-different counts.

    The atomic-mode merge rule is: same location → SUM counts, single
    statement total. We replicate that here so the per-product
    percentages match `go tool cover -func`. Without this, statements
    get multiplied by the number of test binaries, which is what made
    the first version of this script report ~13% when the real number
    was ~52%.
    """
    # block_key -> (stmts, accumulated_count, product, file_path)
    blocks: dict[str, list] = {}

    with open(path) as fh:
        for i, raw in enumerate(fh):
            line = raw.strip()
            if not line:
                continue
            if i == 0 and line.startswith("mode:"):
                continue

            # Format: "<path>:<block> <stmts> <count>"
            try:
                location, stmts_str, count_str = line.rsplit(" ", 2)
                stmts = int(stmts_str)
                count = int(count_str)
            except ValueError:
                # Tolerate malformed lines instead of crashing the
                # whole post-processing step.
                continue

            m = PRODUCT_RE.search(location)
            if not m:
                continue
            slot = m.group(1)

            # Files directly under zscaler/ — e.g. zscaler/oneapiclient.go
            # — have a ".go" extension on the slot. Re-bucket as "core".
            product = "core" if slot.endswith(".go") else slot
            file_path = location.split(":", 1)[0]

            existing = blocks.get(location)
            if existing is None:
                blocks[location] = [stmts, count, product, file_path]
            else:
                # Atomic-mode merge: sum execution counts; statement
                # count is identical across duplicates by definition.
                existing[1] += count

    total = defaultdict(int)
    covered = defaultdict(int)
    files = defaultdict(set)
    for stmts, count, product, file_path in blocks.values():
        total[product] += stmts
        if count > 0:
            covered[product] += stmts
        files[product].add(file_path)

    return total, covered, files


def render_table(total, covered, files) -> None:
    products = sorted(
        total.keys(),
        key=lambda p: ((covered[p] / total[p]) if total[p] else 0.0),
        reverse=True,
    )

    # Column widths chosen so the table fits comfortably in an 80-col
    # terminal even with the longest "Notes" entry.
    headers = ("Product", "Coverage", "Covered / Total", "Files", "Notes")
    widths = (10, 10, 18, 6, 60)

    bar = "=" * (sum(widths) + len(widths) * 1)
    sep = "-" * (sum(widths) + len(widths) * 1)

    print()
    print(f"{BOLD}{CYAN}{bar}{RESET}")
    print(f"{BOLD}{CYAN}  UNIT TEST COVERAGE BY PRODUCT{RESET}")
    print(f"{BOLD}{CYAN}{bar}{RESET}")
    print(
        f"{BOLD}{headers[0]:<{widths[0]}}"
        f"{headers[1]:<{widths[1]}}"
        f"{headers[2]:<{widths[2]}}"
        f"{headers[3]:<{widths[3]}}"
        f"{headers[4]:<{widths[4]}}{RESET}"
    )
    print(sep)

    total_n = 0
    total_c = 0
    for p in products:
        n = total[p]
        c = covered[p]
        if n == 0:
            continue
        pct = (c / n) * 100.0
        total_n += n
        total_c += c
        color = pct_color(pct)
        pct_str = f"{color}{pct:5.1f}%{RESET}"
        note = NOTES.get(p, "")
        print(
            f"{BOLD}{p:<{widths[0]}}{RESET}"
            f"{pct_str:<{widths[1] + len(color) + len(RESET)}}"
            f"{c:>6} / {n:<8}  "
            f"{len(files[p]):<{widths[3]}}"
            f"{DIM}{note}{RESET}"
        )

    print(sep)
    if total_n:
        overall = (total_c / total_n) * 100.0
        ocolor = pct_color(overall)
        pct_str = f"{ocolor}{overall:5.1f}%{RESET}"
        print(
            f"{BOLD}{'TOTAL':<{widths[0]}}{RESET}"
            f"{pct_str:<{widths[1] + len(ocolor) + len(RESET)}}"
            f"{total_c:>6} / {total_n:<8}"
        )
    else:
        print(f"{YELLOW}No coverage data in profile.{RESET}")
    print(f"{BOLD}{CYAN}{bar}{RESET}")
    print()


def parse_profile_for_product(path: str, product: str):
    """Bucket coverage by service within a single product.

    Mirrors parse_profile() but instead of grouping by top-level product
    slot, it groups by the first directory under
    `zscaler/<product>/services/`. Files directly under
    `zscaler/<product>/` (v2_client.go, v2_config.go, …) go into a
    `(root)` bucket so the table is exhaustive without us hand-listing
    every top-level file.

    The same atomic-mode block deduplication applies here as in
    parse_profile() — every test binary in `./tests/unit/...` emits its
    own copy of each block, and naive summation would inflate
    denominators by the binary count.
    """
    prefix = f"github.com/zscaler/zscaler-sdk-go/v3/zscaler/{product}/"
    # block_key -> [stmts, accumulated_count, bucket, file_path]
    blocks: dict[str, list] = {}

    with open(path) as fh:
        for i, raw in enumerate(fh):
            line = raw.strip()
            if not line:
                continue
            if i == 0 and line.startswith("mode:"):
                continue

            try:
                location, stmts_str, count_str = line.rsplit(" ", 2)
                stmts = int(stmts_str)
                count = int(count_str)
            except ValueError:
                continue

            file_path = location.split(":", 1)[0]
            if not file_path.startswith(prefix):
                continue

            tail = file_path[len(prefix):]
            parts = tail.split("/")
            # services/<name>/<file>[/<sub>/...] -> bucket as <name>.
            # We roll up nested service packages (e.g.
            # cloudbrowserisolation/cbibannercontroller) into the
            # top-level service name to keep the table readable; the
            # per-file detail is still available via `go tool cover
            # -func` for anyone who wants to drill in.
            #
            # Require len(parts) >= 3 so a stray top-level file under
            # services/ (e.g. zwa/services/service.go, zdx/services/
            # service.go) doesn't get bucketed as a fake service. Those
            # belong with the product-level "(root)" files.
            if parts[0] == "services" and len(parts) >= 3:
                bucket = parts[1]
            else:
                bucket = "(root)"

            existing = blocks.get(location)
            if existing is None:
                blocks[location] = [stmts, count, bucket, file_path]
            else:
                existing[1] += count

    total = defaultdict(int)
    covered = defaultdict(int)
    files = defaultdict(set)
    for stmts, count, bucket, file_path in blocks.values():
        total[bucket] += stmts
        if count > 0:
            covered[bucket] += stmts
        files[bucket].add(file_path)

    return total, covered, files


def render_service_table(product: str, total, covered, files) -> None:
    # Legacy product-level files (v2_client.go, v2_config.go, …) are not
    # service packages — exclude them from the per-service benchmark table.
    excluded_root = "(root)"
    services = sorted(
        (s for s in total.keys() if s != excluded_root),
        key=lambda s: ((covered[s] / total[s]) if total[s] else 0.0),
        reverse=True,
    )

    # Column widths chosen so even long service names like
    # `cloudbrowserisolation` fit without wrapping.
    headers = ("Service", "Coverage", "Covered / Total", "Files")
    widths = (32, 10, 18, 6)

    bar = "=" * (sum(widths) + len(widths) * 1)
    sep = "-" * (sum(widths) + len(widths) * 1)

    print()
    print(f"{BOLD}{CYAN}{bar}{RESET}")
    print(f"{BOLD}{CYAN}  UNIT TEST COVERAGE — {product.upper()}{RESET}")
    print(f"{BOLD}{CYAN}{bar}{RESET}")
    print(
        f"{BOLD}{headers[0]:<{widths[0]}}"
        f"{headers[1]:<{widths[1]}}"
        f"{headers[2]:<{widths[2]}}"
        f"{headers[3]:<{widths[3]}}{RESET}"
    )
    print(sep)

    total_n = 0
    total_c = 0
    rendered = 0
    for s in services:
        n = total[s]
        c = covered[s]
        if n == 0:
            continue
        rendered += 1
        pct = (c / n) * 100.0
        total_n += n
        total_c += c
        color = pct_color(pct)
        pct_str = f"{color}{pct:5.1f}%{RESET}"
        # Truncate overly-long service names so the columns stay aligned.
        name = s if len(s) <= widths[0] - 1 else s[: widths[0] - 2] + "…"
        print(
            f"{BOLD}{name:<{widths[0]}}{RESET}"
            f"{pct_str:<{widths[1] + len(color) + len(RESET)}}"
            f"{c:>6} / {n:<8}  "
            f"{len(files[s]):<{widths[3]}}"
        )

    print(sep)
    if total_n:
        overall = (total_c / total_n) * 100.0
        ocolor = pct_color(overall)
        pct_str = f"{ocolor}{overall:5.1f}%{RESET}"
        print(
            f"{BOLD}{'TOTAL':<{widths[0]}}{RESET}"
            f"{pct_str:<{widths[1] + len(ocolor) + len(RESET)}}"
            f"{total_c:>6} / {total_n:<8}  "
            f"{sum(len(v) for v in files.values()):<{widths[3]}}"
        )
        # Footer summarises how many distinct services contributed
        # statements — handy gut-check when running a freshly-scaffolded
        # product with only a handful of services covered.
        print(f"{DIM}  {rendered} service(s) instrumented{RESET}")
        if excluded_root in total and total[excluded_root] > 0:
            root_n = total[excluded_root]
            root_c = covered[excluded_root]
            root_pct = (root_c / root_n) * 100.0 if root_n else 0.0
            print(
                f"{DIM}  (root) excluded from total: "
                f"{root_c}/{root_n} stmts ({root_pct:.1f}%) — "
                f"product-level client files, not services/{RESET}"
            )
    else:
        print(
            f"{YELLOW}No coverage data for product '{product}' in "
            f"profile.{RESET}"
        )
    print(f"{BOLD}{CYAN}{bar}{RESET}")
    print()


def main(argv: list[str]) -> int:
    # Tiny hand-rolled arg parser. We avoid argparse so the script stays
    # zero-dependency and trivially callable from the Makefile.
    product: str | None = None
    path: str | None = None

    args = list(argv[1:])
    while args:
        a = args.pop(0)
        if a in ("-h", "--help"):
            print(__doc__)
            return 0
        if a == "--product":
            if not args:
                print("--product requires a value", file=sys.stderr)
                return 2
            product = args.pop(0)
        elif a.startswith("--product="):
            product = a.split("=", 1)[1]
        elif a.startswith("-"):
            print(f"unknown option: {a}", file=sys.stderr)
            return 2
        else:
            path = a

    if path is None:
        path = "unit-coverage.out"

    if not os.path.exists(path):
        print(f"coverage profile not found: {path}", file=sys.stderr)
        return 1

    if product:
        total, covered, files = parse_profile_for_product(path, product)
        render_service_table(product, total, covered, files)
    else:
        total, covered, files = parse_profile(path)
        render_table(total, covered, files)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
