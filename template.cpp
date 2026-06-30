/**
 * template.cpp — C++20 Competitive Programming Template
 * ======================================================
 * Base template loaded automatically by the cpx fetch server when a
 * problem is received from the Competitive Companion browser extension.
 *
 * Key features:
 *   • #define int long long       — all int variables are 64-bit by default
 *   • Fast I/O                    — ios_base::sync_with_stdio(0), cin.tie(0)
 *   • Seeded mt19937_64 RNG       — for randomised algorithms / stress tests
 *   • dbg(...) macro              — debug print to cerr (no-op on judge)
 *   • Execution timing            — elapsed time printed to cerr at end of run
 *
 * Macros and helpers:
 *   rep(i, a, b)     for(int i = a; i < b; ++i)
 *   all(x)           (x).begin(), (x).end()
 *   pb               push_back
 *   umap             unordered_map
 *   nline            "\n"
 *   fun(ret, ...)    std::function<ret(...)>
 *   apv(a, c)        print YES/NO (c=1) or Yes/No (c=0) or -1 (a==-1)
 *   USACO(file)      redirect stdin→file.in, stdout→file.out
 *   chmin(a, b)      a = min(a, b)
 *   chmax(a, b)      a = max(a, b)
 *   rand(l, r)       uniform random integer in [l, r]
 *
 * Debug header (bits/debug.h):
 *   Active when compiled with -DWOOF_ (set automatically by cpx).
 *   Automatically a no-op when submitted to the judge.
 *   See debug++.h for the full list of supported types.
 *
 * Usage:
 *   Write your solution inside solve().
 *   For multi-test problems, uncomment "cin >> tc;" in main().
 *   For USACO file I/O, uncomment "USACO("code");" and set the filename.
 */

#include <bits/stdc++.h>
using namespace std;

// --- Common macros ----------------------------------------------------------
#define all(x)  (x).begin(), (x).end()
#define umap    unordered_map
#define pb      push_back
#define nline   "\n"
#define rep(i, a, b) for(int i = a; i < b; ++i)
#define fun(ret, ...) std::function<ret(__VA_ARGS__)>
#define int long long   // promote all int to 64-bit

// --- Debug macro (active with -DWOOF_, no-op otherwise) --------------------
#ifdef WOOF_
#include <bits/debug.h>
#else
#define dbg(...) 25
#endif

// --- Random number generation -----------------------------------------------
mt19937_64 RNG(chrono::steady_clock::now().time_since_epoch().count());

/** Return a uniform random integer in [l, r]. */
int rand(int l, int r) { return uniform_int_distribution<int>(l, r)(RNG); }

// --- Output helpers ---------------------------------------------------------

/**
 * Print a YES/NO-style answer to stdout.
 *
 * @param a  If -1, prints "-1". If truthy, prints "YES"/"Yes". Otherwise "NO"/"No".
 * @param c  If 1 (default), uses uppercase YES/NO. If 0, uses Yes/No.
 */
void apv(int a = 1, int c = 1) {
    cout << (a == -1 ? "-1" : a ? (c ? "YES" : "Yes") : (c ? "NO" : "No")) << '\n';
}

/**
 * Redirect stdin and stdout to USACO-style file I/O.
 * Call at the top of main() for USACO problems.
 *
 * @param file  Base filename without extension (e.g. "code" → "code.in" / "code.out").
 */
void USACO(const string& file) {
    freopen((file + ".in").c_str(),  "r", stdin);
    freopen((file + ".out").c_str(), "w", stdout);
}

/** Set a = min(a, b). */
void chmin(int &a, int b) { a = min(a, b); }

/** Set a = max(a, b). */
void chmax(int &a, int b) { a = max(a, b); }

// --- Global constants -------------------------------------------------------
const int mod  = 1000000007; // 10^9+7  (swap to 998244353 if required)
const int N    = 1000006;
const int maxN = 3005;
const int inf  = 1e9+5;

// Global arrays (declare problem-specific variables here or inside solve())
int a[N], pre[N];

// --- Solution ---------------------------------------------------------------

void solve() {
    // Write your solution here.
}

// --- Main -------------------------------------------------------------------

int32_t main() {
    auto start = std::chrono::high_resolution_clock::now();

    ios_base::sync_with_stdio(0);
    cin.tie(0);

    // USACO("code");   // ← uncomment for USACO file I/O

    int tc = 1;
    // cin >> tc;       // ← uncomment for multi-test problems

    for (int i = 1; i <= tc; ++i) {
        solve();
    }

    // Print elapsed time to stderr (visible locally, invisible on judge)
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<long double> elapsed = end - start;
    cerr << "Time measured: " << fixed << setprecision(7) << elapsed.count() << " sec";

    return 0;
}