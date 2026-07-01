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

// pair 
template<typename A, typename B>
ostream& operator<<(ostream &os, const pair<A, B> &p) {
    return os << '(' << p.first << ", " << p.second << ')';
}

// tuple
template<typename... Args>
ostream& operator<<(ostream& os, const tuple<Args...>& t) {
    os << '(';
    apply([&os](const Args&... args) {
        size_t n = 0;
        ((os << args << (++n != sizeof...(Args) ? ", " : "")), ...);
    }, t);
    return os << ')';
}

// containers (vector, set, map, deque, string, etc.)
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type>
ostream& operator<<(ostream &os, const T_container &v) {
    os << '{';
    string sep;
    for (const T &x : v) os << sep << x, sep = ", ";
    return os << '}';
}

// stack
template <typename T>
ostream& operator<<(ostream& os, stack<T> st) {
    vector<T> temp;
    while (!st.empty()) {
        temp.push_back(st.top());
        st.pop();
    }
    reverse(temp.begin(), temp.end());
    return os << temp;
}

// queue
template <typename T>
ostream& operator<<(ostream& os, queue<T> q) {
    vector<T> temp;
    while (!q.empty()) {
        temp.push_back(q.front());
        q.pop();
    }
    return os << temp;
}

// priority queue (max-heap by default)
template <typename T>
ostream& operator<<(ostream& os, priority_queue<T> pq) {
    vector<T> temp;
    while (!pq.empty()) {
        temp.push_back(pq.top());
        pq.pop();
    }
    reverse(temp.begin(), temp.end());
    return os << temp;
}

// priority queue (min-heap)
template <typename T, typename Container, typename Compare>
ostream& operator<<(ostream& os, priority_queue<T, Container, Compare> pq) {
    vector<T> temp;
    while (!pq.empty()) {
        temp.push_back(pq.top());
        pq.pop();
    }
    return os << temp;
}

// C-style 1D arrays
template <typename T, size_t N>
static void dbgg(T (&arr)[N]) {
    cerr << '{';
    for (size_t i = 0; i < N; ++i) {
        cerr << arr[i];
        if (i != N - 1) cerr << ", ";
    }
    cerr << "}\n";
}

// C-style 2D arrays 
template <typename T, size_t N, size_t M>
static void dbgg(T (&matrix)[N][M], int t = 0, int one_based = 0) {
    if(t){
        cerr << left << setw(4) << " ";
        for(int i = 0; i < M; ++i) cerr << left << setw(4) << i+one_based;
        cerr << "\n";
    }
    for (size_t i = 0; i < N; ++i) {
        if(t) cerr << left << setw(4) << (to_string(i+one_based) + " -");
        for (size_t j = 0; j < M; ++j) {
            cerr << left << setw(4) << matrix[i][j];
        }
        if (i != N - 1) cerr << ('\n'); 
    }
}

// 2D vectors
template <typename T>
static void dbgg(const vector<vector<T>>& matrix, int t = 0, int one_based = 0) {
    size_t N = matrix.size();
    size_t M = (N > 0) ? matrix[0].size() : 0;
    if (t) {
        cerr << left << setw(4) << " ";
        for (size_t i = 0; i < M; ++i) cerr << left << setw(4) << i+one_based;
        cerr << "\n";
    }
    for (size_t i = 0; i < N; ++i) {
        if (t) cerr << left << setw(4) << (to_string(i+one_based) + " -");
        for (size_t j = 0; j < M; ++j) {
            cerr << left << setw(4) << matrix[i][j];
        }
        if (i != N - 1) cerr << '\n';
    }
}

// debug utilities
void dbg_out() { cerr << endl; }
template<typename Head, typename... Tail>
void dbg_out(Head H, Tail... T) {
    cerr << ' ' << H;
    dbg_out(T...);
}

#ifdef WOOF_
#define dbg(...) cerr << '[' << __LINE__ << "] (" << #__VA_ARGS__ << "):", dbg_out(__VA_ARGS__)
#else
#define dbg(...)
#endif
