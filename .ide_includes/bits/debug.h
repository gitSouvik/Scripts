/**
 * debug.h — Competitive Programming Debug Macros (basic)
 * =======================================================
 * Provides a `dbg(...)` variadic macro that prints variable names and values
 * to stderr with file line numbers. Enabled only when compiled with -DWOOF_.
 *
 * Supported types:
 *   pair<A, B>              →  (a, b)
 *   any STL container       →  {x, y, z}   (vector, set, deque, string, …)
 *   map<K, V>               →  {k: v, …}
 *
 * For extended support (tuple, stack, queue, priority_queue, 2D arrays),
 * use debug++.h instead.
 *
 * Installation:
 *   Copy this file to your compiler's bits/ include directory so that
 *   #include <bits/debug.h> resolves correctly. Example (macOS, GCC 14):
 *     sudo cp debug.h /usr/local/include/bits/debug.h
 *
 * Usage in code:
 *   int x = 5;
 *   vector<int> v = {1, 2, 3};
 *   dbg(x);        // [14] (x): 5
 *   dbg(x, v);     // [15] (x, v): 5 {1, 2, 3}
 *
 * Note: When -DWOOF_ is NOT defined (e.g. on the judge), dbg expands to
 * nothing and produces zero overhead.
 */

#include <bits/stdc++.h>
using namespace std;

// --- Stream operator overloads ----------------------------------------------

// pair<A, B>  →  (first, second)
template<typename A, typename B>
ostream& operator<<(ostream &os, const pair<A, B> &p) {
    return os << '(' << p.first << ", " << p.second << ')';
}

// Any STL container (vector, set, deque, etc., but NOT string)  →  {x, y, z}
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type>
ostream& operator<<(ostream &os, const T_container &v) {
    os << '{';
    string sep;
    for (const T &x : v) os << sep << x, sep = ", ";
    return os << '}';
}

// map<K, V>  →  {k: v, k: v, …}
template<typename K, typename V>
ostream& operator<<(ostream& os, const map<K, V>& m) {
    os << '{';
    string sep;
    for (const auto& [key, val] : m) os << sep << key << ": " << val, sep = ", ";
    return os << '}';
}

// --- Internal debug output helper -------------------------------------------

void dbg_out() { cerr << endl; }

template<typename Head, typename... Tail>
void dbg_out(Head H, Tail... T) {
    cerr << ' ' << H;
    dbg_out(T...);
}

// --- Public macro -----------------------------------------------------------

#ifdef WOOF_
// Active when compiled with -DWOOF_:
// Prints  [<line>] (<var names>): <values>
#define dbg(...) cerr << '[' << __LINE__ << "] (" << #__VA_ARGS__ << "):", dbg_out(__VA_ARGS__)
#else
// No-op when -DWOOF_ is not defined (judge build)
#define dbg(...)
#endif
