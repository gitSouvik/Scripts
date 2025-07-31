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
// Overloading the << operator for ALL types
template<typename A, typename B>                                                           // For pair 
ostream& operator<<(ostream &os, const pair<A, B> &p) {
    return os << '(' << p.first << ", " << p.second << ')';
}
template<typename... Args>                                                                 // For tuple
ostream& operator<<(ostream& os, const tuple<Args...>& t) {
    os << '(';
    apply([&os](const Args&... args) {
        size_t n = 0;
        ((os << args << (++n != sizeof...(Args) ? ", " : "")), ...);
    }, t);
    return os << ')' /*<< '\n'*/ ;
}
// For containers (vector, set, map, deque, string, etc.)
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type>
ostream& operator<<(ostream &os, const T_container &v) {
    os << '{';
    string sep;
    for (const T &x : v) os << sep << x, sep = ", ";
    return os << '}' /*<< '\n'*/ ;
}
template <typename T>                                                                      // For Stack
ostream& operator<<(ostream& os, stack<T> st) {
    vector<T> temp;
    while (!st.empty()) {
        temp.push_back(st.top());
        st.pop();
    }
    reverse(temp.begin(), temp.end());
    return os << temp;
}
template <typename T>                                                                      // For Queue
ostream& operator<<(ostream& os, queue<T> q) {
    vector<T> temp;
    while (!q.empty()) {
        temp.push_back(q.front());
        q.pop();
    }
    return os << temp;
}
template <typename T>                                                                       // For Priority Queue (max-heap by default)
ostream& operator<<(ostream& os, priority_queue<T> pq) {
    vector<T> temp;
    while (!pq.empty()) {
        temp.push_back(pq.top());
        pq.pop();
    }
    reverse(temp.begin(), temp.end());  // max-heap: reverse to print in ascending order
    return os << temp;
}
template <typename T, typename Container, typename Compare>                                 // For Priority_Queue (min-heap)
ostream& operator<<(ostream& os, priority_queue<T, Container, Compare> pq) {
    vector<T> temp;
    while (!pq.empty()) {
        temp.push_back(pq.top());
        pq.pop();
    }
    return os << temp;  // min-heap: already in ascending order
}
template <typename T, size_t N>                                                             // For C-style 1D arrays
static void dbgg(T (&arr)[N]) {
    cerr << '{';
    for (size_t i = 0; i < N; ++i) {
        cerr << arr[i];
        if (i != N - 1) cerr << ", "; // Use cerr << setw(3) << arr[i] ;
    }
    cerr << "}\n";
}
template <typename T, size_t N, size_t M>                                                   // For C-style 2D arrays 
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
template <typename T>                                                                       // For 2D vectors
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
// Debugging utilities
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
