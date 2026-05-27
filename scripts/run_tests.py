#!/usr/bin/env python3
"""
Zscaler Go SDK — Interactive Test Runner

Presents a menu to select which product suite to test, runs the Go test
suite, parses results in real time, and displays a summary table.

Usage:
    python3 scripts/run_tests.py
    # or
    chmod +x scripts/run_tests.py && ./scripts/run_tests.py
"""

import json
import os
import platform
import re
import subprocess
import sys
import time
from collections import defaultdict
from datetime import datetime

# ─── Configuration ───────────────────────────────────────────────────────────

REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))

PRODUCTS = {
    "1": {"name": "ZPA", "path": "./zscaler/zpa/services/..."},
    "2": {"name": "ZIA", "path": "./zscaler/zia/services/..."},
    "3": {"name": "ZCC", "path": "./zscaler/zcc/services/..."},
    "4": {"name": "ZDX", "path": "./zscaler/zdx/services/..."},
    "5": {"name": "ZTW", "path": "./zscaler/ztw/services/..."},
    "6": {"name": "ZID", "path": "./zscaler/zid/services/..."},
}

CLIENT_MODES = {
    "1": {"name": "OneAPI (default)", "env_value": ""},
    "2": {"name": "ZPA Legacy", "env_value": "zpa_legacy"},
}

# ─── Acceptance Thresholds ───────────────────────────────────────────────────
# Tiered grading for integration test suites. Integration tests are inherently
# flakier than unit tests (network, rate limits, API transient errors), so a
# strict 100% gate is impractical for periodic validation runs.
#
# Override via environment variables if needed:
#   ZSCALER_TEST_THRESHOLD_PASS=98
#   ZSCALER_TEST_THRESHOLD_CONDITIONAL=95
#   ZSCALER_TEST_THRESHOLD_UNSTABLE=85

THRESHOLDS = {
    "pass":        float(os.getenv("ZSCALER_TEST_THRESHOLD_PASS", "98")),
    "conditional": float(os.getenv("ZSCALER_TEST_THRESHOLD_CONDITIONAL", "95")),
    "unstable":    float(os.getenv("ZSCALER_TEST_THRESHOLD_UNSTABLE", "85")),
}

def compute_verdict(pass_rate, total_tests):
    """Return (label, color_fn, description) based on pass rate."""
    if total_tests == 0:
        return ("NO TESTS", yellow, "No tests were executed")
    if pass_rate >= 100.0:
        return ("PASS", green, "All tests passed")
    if pass_rate >= THRESHOLDS["pass"]:
        return ("PASS (warnings)", green,
                "Negligible failures — likely transient/flaky")
    if pass_rate >= THRESHOLDS["conditional"]:
        return ("CONDITIONAL PASS", yellow,
                "Acceptable for integration suites — review failures")
    if pass_rate >= THRESHOLDS["unstable"]:
        return ("UNSTABLE", lambda t: f"\033[38;5;208m{t}\033[0m",
                "Needs attention — pattern of failures emerging")
    return ("FAIL", red, "Significant regression — action required")


# ─── Helpers ─────────────────────────────────────────────────────────────────


def get_go_version():
    try:
        out = subprocess.check_output(["go", "version"], text=True).strip()
        match = re.search(r"go(\d+\.\d+(\.\d+)?)", out)
        return match.group(1) if match else out
    except Exception:
        return "unknown"


def get_sdk_version():
    version_file = os.path.join(REPO_ROOT, "zscaler", "oneapiclient.go")
    try:
        with open(version_file) as f:
            for line in f:
                m = re.search(r'VERSION\s*=\s*"([^"]+)"', line)
                if m:
                    return m.group(1)
    except Exception:
        pass
    return "unknown"


def get_git_branch():
    try:
        return subprocess.check_output(
            ["git", "rev-parse", "--abbrev-ref", "HEAD"],
            cwd=REPO_ROOT, text=True
        ).strip()
    except Exception:
        return "unknown"


def clear_screen():
    os.system("cls" if os.name == "nt" else "clear")


def bold(text):
    return f"\033[1m{text}\033[0m"


def green(text):
    return f"\033[92m{text}\033[0m"


def red(text):
    return f"\033[91m{text}\033[0m"


def yellow(text):
    return f"\033[93m{text}\033[0m"


def cyan(text):
    return f"\033[96m{text}\033[0m"


def dim(text):
    return f"\033[2m{text}\033[0m"


# ─── Branding ─────────────────────────────────────────────────────────────────

_ZSCALER_ART = [
    "███████╗███████╗ ██████╗ █████╗ ██╗     ███████╗██████╗ ",
    "╚══███╔╝██╔════╝██╔════╝██╔══██╗██║     ██╔════╝██╔══██╗",
    "  ███╔╝ ███████╗██║     ███████║██║     █████╗  ██████╔╝",
    " ███╔╝  ╚════██║██║     ██╔══██║██║     ██╔══╝  ██╔══██╗",
    "███████╗███████║╚██████╗██║  ██║███████╗███████╗██║  ██║",
    "╚══════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝",
]
_TAGLINE = "Go SDK — Interactive Test Runner"
_RESET = "\x1b[0m"


def _supports_truecolor():
    if not sys.stdout.isatty():
        return False
    if os.environ.get("NO_COLOR"):
        return False
    if os.environ.get("COLORTERM", "").lower() in ("truecolor", "24bit"):
        return True
    term = os.environ.get("TERM", "").lower()
    return "256color" in term or "kitty" in term or "iterm" in term


def _rgb(r, g, b):
    return f"\x1b[38;2;{r};{g};{b}m"


def show_banner():
    width = max(len(line) for line in _ZSCALER_ART)
    pad = 2
    inner = width + pad * 2

    if not _supports_truecolor():
        print()
        print(f"  +{'-' * inner}+")
        for line in _ZSCALER_ART:
            print(f"  |{' ' * pad}{line.ljust(width)}{' ' * pad}|")
        print(f"  +{'-' * inner}+")
        print(f"  {_TAGLINE}")
        print()
        return

    start = (0x55, 0xCC, 0xFF)
    end = (0x00, 0x3D, 0x99)
    border = _rgb(0x33, 0x55, 0x99)
    shadow_color = _rgb(0x00, 0x3D, 0x99)

    def gradient_line(text):
        out = []
        last = None
        padded = text.ljust(width)
        for i, ch in enumerate(padded):
            if ch == " ":
                out.append(" ")
                continue
            t = i / max(width - 1, 1)
            r = int(start[0] + (end[0] - start[0]) * t)
            g = int(start[1] + (end[1] - start[1]) * t)
            b = int(start[2] + (end[2] - start[2]) * t)
            color = (r, g, b)
            if color != last:
                out.append(_rgb(*color))
                last = color
            out.append(ch)
        out.append(_RESET)
        return "".join(out)

    blank = " " * width
    print()
    print(f"  {border}╭{'─' * inner}╮{_RESET}")
    print(f"  {border}│{_RESET}{' ' * pad}{blank}{' ' * pad}{border}│{_RESET}")
    for line in _ZSCALER_ART:
        print(f"  {border}│{_RESET}{' ' * pad}{gradient_line(line)}{' ' * pad}{border}│{_RESET}")
    shadow = "░" * width
    print(f"  {border}│{_RESET}{' ' * pad}{shadow_color}{shadow}{_RESET}{' ' * pad}{border}│{_RESET}")
    print(f"  {border}│{_RESET}{' ' * pad}{blank}{' ' * pad}{border}│{_RESET}")
    print(f"  {border}╰{'─' * inner}╯{_RESET}")
    print(f"  {_TAGLINE}")
    print()


# ─── Keyboard Interrupt Handling ─────────────────────────────────────────────


def graceful_exit(msg="Interrupted by user."):
    """Print a clean message and exit without traceback."""
    print(f"\n\n  {yellow(msg)} Exiting.\n")
    sys.exit(130)


def safe_input(prompt):
    """input() wrapper that catches Ctrl+C / Ctrl+D cleanly."""
    try:
        return input(prompt)
    except KeyboardInterrupt:
        graceful_exit()
    except EOFError:
        graceful_exit("End of input.")


# ─── Menu ────────────────────────────────────────────────────────────────────


def select_product():
    print(bold("  Select product to test:"))
    print()
    for key, product in PRODUCTS.items():
        print(f"    {cyan(key)}) {product['name']:6s}  {dim(product['path'])}")
    print()
    print(f"    {cyan('a')}) Run ALL products")
    print(f"    {cyan('q')}) Quit")
    print()

    choice = safe_input("  Enter choice: ").strip().lower()
    if choice == "q":
        sys.exit(0)
    if choice == "a":
        return list(PRODUCTS.values())
    if choice in PRODUCTS:
        return [PRODUCTS[choice]]
    print(red("\n  Invalid selection. Try again.\n"))
    return select_product()


def select_client_mode():
    print()
    print(bold("  Select client mode:"))
    print()
    for key, mode in CLIENT_MODES.items():
        print(f"    {cyan(key)}) {mode['name']}")
    print()

    choice = safe_input("  Enter choice [1]: ").strip() or "1"
    if choice in CLIENT_MODES:
        return CLIENT_MODES[choice]
    print(red("\n  Invalid selection. Defaulting to OneAPI.\n"))
    return CLIENT_MODES["1"]


def ask_exclusions(product_name):
    print()
    resp = safe_input(f"  Exclude any packages from {product_name}? (comma-separated keywords, or Enter to skip): ").strip()
    if resp:
        return [x.strip() for x in resp.split(",") if x.strip()]
    return []


def ask_verbose():
    print()
    resp = safe_input("  Verbose output (-v)? [Y/n]: ").strip().lower()
    return resp != "n"


# ─── Test Execution ──────────────────────────────────────────────────────────


def resolve_packages(test_path, exclusions):
    """Use go list to expand the path, then filter exclusions."""
    try:
        out = subprocess.check_output(
            ["go", "list", test_path],
            cwd=REPO_ROOT, text=True, stderr=subprocess.DEVNULL
        )
        packages = [line.strip() for line in out.splitlines() if line.strip()]
    except subprocess.CalledProcessError:
        return [test_path]

    if exclusions:
        packages = [
            p for p in packages
            if not any(exc.lower() in p.lower() for exc in exclusions)
        ]
    return packages


def run_tests(packages, client_mode, verbose):
    """Run go test -json and parse results in real time."""
    env = os.environ.copy()
    if client_mode["env_value"]:
        env["ZSCALER_SDK_TEST_CLIENT"] = client_mode["env_value"]

    cmd = ["go", "test", "-json", "-count=1", "-cover"] + packages
    if verbose:
        pass  # -json already captures verbose info

    results = {
        "passed": 0,
        "failed": 0,
        "skipped": 0,
        "total_tests": 0,
        "packages_tested": 0,
        "packages_passed": 0,
        "packages_failed": 0,
        "coverage": {},
        "failed_tests": [],
        "durations": {},
    }

    package_status = defaultdict(lambda: "pass")
    test_count_by_pkg = defaultdict(int)

    print()
    print(bold("  Running tests..."))
    print(f"  {dim('Command:')} {' '.join(cmd[:6])} ... ({len(packages)} packages)")
    print()

    start_time = time.time()

    proc = subprocess.Popen(
        cmd,
        cwd=REPO_ROOT,
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    current_counts = {"pass": 0, "fail": 0, "skip": 0}

    try:
        for line in proc.stdout:
            line = line.strip()
            if not line:
                continue
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                continue

            action = event.get("Action", "")
            test_name = event.get("Test", "")
            pkg = event.get("Package", "")
            elapsed = event.get("Elapsed", 0)
            output_line = event.get("Output", "")

            # Track coverage from output lines
            if output_line and "coverage:" in output_line:
                m = re.search(r"coverage:\s+([\d.]+)%", output_line)
                if m:
                    results["coverage"][pkg] = float(m.group(1))

            if test_name:
                if action == "pass":
                    results["passed"] += 1
                    current_counts["pass"] += 1
                    test_count_by_pkg[pkg] += 1
                elif action == "fail":
                    results["failed"] += 1
                    current_counts["fail"] += 1
                    test_count_by_pkg[pkg] += 1
                    results["failed_tests"].append(f"{pkg}: {test_name}")
                    package_status[pkg] = "fail"
                elif action == "skip":
                    results["skipped"] += 1
                    current_counts["skip"] += 1
            elif action == "fail" and not test_name:
                package_status[pkg] = "fail"
            elif action == "pass" and not test_name and pkg:
                results["durations"][pkg] = elapsed

            # Live progress
            total_so_far = current_counts["pass"] + current_counts["fail"] + current_counts["skip"]
            sys.stdout.write(
                f"\r  Progress: {green(str(current_counts['pass']))} passed | "
                f"{red(str(current_counts['fail']))} failed | "
                f"{yellow(str(current_counts['skip']))} skipped | "
                f"{total_so_far} total"
            )
            sys.stdout.flush()

    except KeyboardInterrupt:
        print(f"\n\n  {yellow('Tests interrupted by user.')} Terminating child process...")
        proc.terminate()
        try:
            proc.wait(timeout=5)
        except subprocess.TimeoutExpired:
            proc.kill()
            proc.wait()
        print(f"  {dim('Partial results collected up to interruption point.')}")

    proc.wait()
    end_time = time.time()

    print("\n")

    results["total_tests"] = results["passed"] + results["failed"] + results["skipped"]
    results["elapsed_seconds"] = end_time - start_time
    results["packages_tested"] = len(set(test_count_by_pkg.keys()) | set(results["durations"].keys()))
    results["packages_failed"] = sum(1 for s in package_status.values() if s == "fail")
    results["packages_passed"] = results["packages_tested"] - results["packages_failed"]
    results["exit_code"] = proc.returncode

    return results


# ─── Results Display ─────────────────────────────────────────────────────────


def format_duration(seconds):
    if seconds < 60:
        return f"{seconds:.1f}s"
    minutes = int(seconds // 60)
    secs = seconds % 60
    return f"{minutes}m {secs:.0f}s"


def print_table(product_name, client_mode, results):
    """Print a formatted results table."""
    go_version = get_go_version()
    sdk_version = get_sdk_version()
    branch = get_git_branch()
    runner = platform.node() or "local"
    runner_os = f"{platform.system()} {platform.machine()}"
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    avg_coverage = 0.0
    if results["coverage"]:
        avg_coverage = sum(results["coverage"].values()) / len(results["coverage"])

    pass_rate = 0.0
    if results["total_tests"] > 0:
        pass_rate = (results["passed"] / results["total_tests"]) * 100

    verdict_label, verdict_color, verdict_desc = compute_verdict(pass_rate, results["total_tests"])

    # Table width
    W = 68

    def sep():
        return "  +" + "-" * (W - 2) + "+"

    def row(label, value):
        label_str = f"  | {label:<22}"
        value_str = f"{value}"
        padding = W - 2 - len(label) - 22 - len(value) + len(label)
        # Simpler: fixed width columns
        return f"  | {label:<22}| {value:<{W - 27}}|"

    print(bold("  ╔══════════════════════════════════════════════════════════════════╗"))
    print(bold(f"  ║{'TEST RESULTS SUMMARY':^66}║"))
    print(bold("  ╚══════════════════════════════════════════════════════════════════╝"))
    print()

    # Header section
    print(sep())
    print(f"  | {'ENVIRONMENT':<64}|")
    print(sep())
    print(f"  | {'Product':<22}| {product_name:<41}|")
    print(f"  | {'Client Mode':<22}| {client_mode['name']:<41}|")
    print(f"  | {'Go Version':<22}| {go_version:<41}|")
    print(f"  | {'SDK Version':<22}| {'v' + sdk_version:<41}|")
    print(f"  | {'Git Branch':<22}| {branch:<41}|")
    print(f"  | {'Runner':<22}| {runner:<41}|")
    print(f"  | {'OS / Arch':<22}| {runner_os:<41}|")
    print(f"  | {'Timestamp':<22}| {timestamp:<41}|")
    print(sep())
    print()

    # Verdict section
    print(sep())
    print(f"  | {'VERDICT':<64}|")
    print(sep())
    print(f"  | {'Status':<22}| {verdict_color(verdict_label):<50}|")
    print(f"  | {'Reason':<22}| {verdict_desc:<41}|")
    print(f"  | {'Pass Rate':<22}| {pass_rate:.1f}%{'':<38}|")
    print(f"  | {'Threshold':<22}| {'>='+str(THRESHOLDS['conditional'])+'% = Conditional Pass':<41}|")
    print(sep())
    print()

    # Results section
    print(sep())
    print(f"  | {'TEST RESULTS':<64}|")
    print(sep())
    print(f"  | {'Total Tests':<22}| {results['total_tests']:<41}|")
    print(f"  | {'Passed':<22}| {green(str(results['passed'])):<50}|")
    print(f"  | {'Failed':<22}| {red(str(results['failed'])):<50}|")
    print(f"  | {'Skipped':<22}| {yellow(str(results['skipped'])):<50}|")
    print(f"  | {'Pass Rate':<22}| {pass_rate:.1f}%{'':<38}|")
    print(sep())
    print()

    # Package section
    cov_display = f"{avg_coverage:.1f}%" if avg_coverage > 0 else "n/a"
    print(sep())
    print(f"  | {'PACKAGES':<64}|")
    print(sep())
    print(f"  | {'Total Packages':<22}| {results['packages_tested']:<41}|")
    print(f"  | {'Packages Passed':<22}| {green(str(results['packages_passed'])):<50}|")
    print(f"  | {'Packages Failed':<22}| {red(str(results['packages_failed'])):<50}|")
    print(f"  | {'Code Coverage':<22}| {cov_display:<41}|")
    print(f"  | {'Total Duration':<22}| {format_duration(results['elapsed_seconds']):<41}|")
    print(sep())
    print()

    # Failed tests detail
    if results["failed_tests"]:
        print(sep())
        print(f"  | {red('FAILED TESTS'):<73}|")
        print(sep())
        for ft in results["failed_tests"][:20]:
            display = ft if len(ft) <= 62 else "..." + ft[-(62 - 3):]
            print(f"  | {display:<64}|")
        if len(results["failed_tests"]) > 20:
            remaining = len(results["failed_tests"]) - 20
            print(f"  | {dim(f'... and {remaining} more'):<73}|")
        print(sep())
        print()

    # Per-package coverage breakdown (top 15 by duration)
    if results["durations"]:
        sorted_pkgs = sorted(results["durations"].items(), key=lambda x: x[1], reverse=True)[:15]
        print(sep())
        print(f"  | {'TOP PACKAGES BY DURATION':<64}|")
        print(sep())
        print(f"  | {'Package':<44}| {'Time':<8}| {'Cov':<8}|")
        print(f"  |{'-' * 45}|{'-' * 9}|{'-' * 9}|")
        for pkg, dur in sorted_pkgs:
            short_pkg = pkg.split("/services/")[-1] if "/services/" in pkg else pkg.split("/")[-1]
            cov = results["coverage"].get(pkg, 0)
            cov_str = f"{cov:.0f}%" if cov > 0 else "-"
            print(f"  | {short_pkg:<44}| {dur:<7.1f}s| {cov_str:<8}|")
        print(sep())
        print()

    # Acceptance criteria legend
    print(sep())
    print(f"  | {'ACCEPTANCE CRITERIA':<64}|")
    print(sep())
    print(f"  | {green('PASS'):<50}| {'100% pass rate':<22}|")
    print(f"  | {green('PASS (warnings)'):<50}| {'>= ' + str(THRESHOLDS['pass']) + '% pass rate':<22}|")
    print(f"  | {yellow('CONDITIONAL PASS'):<50}| {'>= ' + str(THRESHOLDS['conditional']) + '% pass rate':<22}|")
    orange = lambda t: f"\033[38;5;208m{t}\033[0m"
    print(f"  | {orange('UNSTABLE'):<50}| {'>= ' + str(THRESHOLDS['unstable']) + '% pass rate':<22}|")
    print(f"  | {red('FAIL'):<50}| {'< ' + str(THRESHOLDS['unstable']) + '% pass rate':<22}|")
    print(sep())
    print()


# ─── Main ────────────────────────────────────────────────────────────────────


def exit_code_for_verdict(pass_rate, total_tests):
    """Map verdict to a process exit code for CI integration.

    Exit codes:
        0 — PASS or PASS (warnings)
        1 — FAIL
        2 — UNSTABLE
        3 — CONDITIONAL PASS (success, but worth reviewing)
    """
    if total_tests == 0:
        return 1
    if pass_rate >= 100.0:
        return 0
    if pass_rate >= THRESHOLDS["pass"]:
        return 0
    if pass_rate >= THRESHOLDS["conditional"]:
        return 3
    if pass_rate >= THRESHOLDS["unstable"]:
        return 2
    return 1


def main():
    os.chdir(REPO_ROOT)
    clear_screen()
    show_banner()

    products = select_product()
    client_mode = select_client_mode()
    verbose = ask_verbose()

    worst_exit = 0

    for product in products:
        exclusions = ask_exclusions(product["name"])

        print()
        print(bold(f"  {'═' * 60}"))
        print(bold(f"  Testing: {product['name']}"))
        print(bold(f"  Client:  {client_mode['name']}"))
        if exclusions:
            exc_str = ", ".join(exclusions)
            print(f"  {dim('Excluding: ' + exc_str)}")
        print(bold(f"  {'═' * 60}"))

        packages = resolve_packages(product["path"], exclusions)
        if not packages:
            print(red(f"\n  No packages found for {product['name']} after exclusions.\n"))
            continue

        results = run_tests(packages, client_mode, verbose)
        print_table(product["name"], client_mode, results)

        pass_rate = 0.0
        if results["total_tests"] > 0:
            pass_rate = (results["passed"] / results["total_tests"]) * 100
        code = exit_code_for_verdict(pass_rate, results["total_tests"])
        worst_exit = max(worst_exit, code)

    print(bold("  Done."))
    print()
    sys.exit(worst_exit)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        graceful_exit()
    except BrokenPipeError:
        sys.exit(0)
