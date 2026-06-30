/**
 * debug++.h — Competitive Programming Debug Macros (extended)
 * ============================================================
 * Extended version of debug.h. Provides `dbg(...)` and `dbgg(...)` for
 * printing virtually any C++ type to stderr during local testing. Enabled
 * only when compiled with -DWOOF_; a complete no-op on the judge.
 *
 * Supported types (in addition to debug.h basics):
 *   pair<A, B>              →  (a, b)
 *   tuple<...>              →  (a, b, c, …)
 *   any STL container       →  {x, y, z}
 *   stack<T>                →  {bottom, …, top}  (non-destructive)
 *   queue<T>                →  {front, …, back}  (non-destructive)
 *   priority_queue<T>       →  {…, max}   (max-heap, ascending order)
 *   priority_queue<T,C,Cmp> →  {min, …}   (min-heap)
 *   T arr[N]                →  {x, y, z}  (C-style 1D array)
 *   T mat[N][M]             →  formatted grid (see dbgg below)
 *   vector<vector<T>>       →  formatted grid (see dbgg below)
 *
 * Installation:
 *   Copy as bits/debug.h in your compiler's include path. Example (macOS):
 *     sudo cp debug++.h /usr/local/include/bits/debug.h
 *
 * Usage:
 *   dbg(x, v, m);        // standard debug print for any supported type
 *   dbgg(arr);           // 1D C-style array
 *   dbgg(mat, 1, 0);     // 2D array/matrix with index labels, 0-based
 *   dbgg(mat, 1, 1);     // 2D array/matrix with index labels, 1-based
 */

#include <iostream>
#include <vector>
#include <set>
#include <map>
#include <unordered_map>
#include <unordered_set>
#include <iomanip>
#include <deque>
#include <queue>
#include <stack>
#include <tuple>
#include <type_traits>
using namespace std;

// --- Stream operator overloads ----------------------------------------------

// pair<A, B>  →  (first, second)
template<typename A, typename B>
ostream& operator<<(ostream &os, const pair<A, B> &p) {
    return os << '(' << p.first << ", " << p.second << ')';
}

// tuple<Args...>  →  (a, b, c, …)
template<typename... Args>
ostream& operator<<(ostream& os, const tuple<Args...>& t) {
    os << '(';
    apply([&os](const Args&... args) {
        size_t n = 0;
        ((os << args << (++n != sizeof...(Args) ? ", " : "")), ...);
    }, t);
    return os << ')';
}

// Any STL container (vector, set, deque, etc., but NOT string)  →  {x, y, z}
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type>
ostream& operator<<(ostream &os, const T_container &v) {
    os << '{';
    string sep;
    for (const T &x : v) os << sep << x, sep = ", ";
    return os << '}';
}

// stack<T>  →  {bottom, …, top}  (copy to avoid mutating the original)
template <typename T>
ostream& operator<<(ostream& os, stack<T> st) {
    vector<T> temp;
    while (!st.empty()) { temp.push_back(st.top()); st.pop(); }
    reverse(temp.begin(), temp.end());
    return os << temp;
}

// queue<T>  →  {front, …, back}
template <typename T>
ostream& operator<<(ostream& os, queue<T> q) {
    vector<T> temp;
    while (!q.empty()) { temp.push_back(q.front()); q.pop(); }
    return os << temp;
}

// priority_queue<T> (max-heap)  →  {…, max}  (ascending order)
template <typename T>
ostream& operator<<(ostream& os, priority_queue<T> pq) {
    vector<T> temp;
    while (!pq.empty()) { temp.push_back(pq.top()); pq.pop(); }
    reverse(temp.begin(), temp.end());
    return os << temp;
}

// priority_queue<T, Container, Compare> (min-heap)  →  {min, …}
template <typename T, typename Container, typename Compare>
ostream& operator<<(ostream& os, priority_queue<T, Container, Compare> pq) {
    vector<T> temp;
    while (!pq.empty()) { temp.push_back(pq.top()); pq.pop(); }
    return os << temp;
}

// --- dbgg: grid / array printers --------------------------------------------

// C-style 1D array  →  {x, y, z}
template <typename T, size_t N>
static void dbgg(T (&arr)[N]) {
    cerr << '{';
    for (size_t i = 0; i < N; ++i) {
        cerr << arr[i];
        if (i != N - 1) cerr << ", ";
    }
    cerr << "}\n";
}

// C-style 2D array — optional column/row index labels.
// t=1 enables index display; one_based=1 starts indices from 1.
template <typename T, size_t N, size_t M>
static void dbgg(T (&matrix)[N][M], int t = 0, int one_based = 0) {
    if (t) {
        cerr << left << setw(4) << " ";
        for (int i = 0; i < (int)M; ++i) cerr << left << setw(4) << i + one_based;
        cerr << "\n";
    }
    for (size_t i = 0; i < N; ++i) {
        if (t) cerr << left << setw(4) << (to_string(i + one_based) + " -");
        for (size_t j = 0; j < M; ++j) cerr << left << setw(4) << matrix[i][j];
        if (i != N - 1) cerr << '\n';
    }
}

// vector<vector<T>> — same optional index labels as the 2D array version.
template <typename T>
static void dbgg(const vector<vector<T>>& matrix, int t = 0, int one_based = 0) {
    size_t N = matrix.size();
    size_t M = (N > 0) ? matrix[0].size() : 0;
    if (t) {
        cerr << left << setw(4) << " ";
        for (size_t i = 0; i < M; ++i) cerr << left << setw(4) << i + one_based;
        cerr << "\n";
    }
    for (size_t i = 0; i < N; ++i) {
        if (t) cerr << left << setw(4) << (to_string(i + one_based) + " -");
        for (size_t j = 0; j < M; ++j) cerr << left << setw(4) << matrix[i][j];
        if (i != N - 1) cerr << '\n';
    }
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
