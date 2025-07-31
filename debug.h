#include <bits/stdc++.h>
using namespace std;

// For pair 
template<typename A, typename B>                                                           
ostream& operator<<(ostream &os, const pair<A, B> &p) {
    return os << '(' << p.first << ", " << p.second << ')';
}
// For containers (vector, set, map, deque, string, etc.)
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type>
ostream& operator<<(ostream &os, const T_container &v) {
    os << '{';
    string sep;
    for (const T &x : v) os << sep << x, sep = ", ";
    return os << '}' /*<< '\n'*/ ;
}
// map<int, container<pair>>> // not that necessary
template<typename K, typename V>
ostream& operator<<(ostream& os, const map<K, V>& m) {
    os << '{';
    string sep;
    for (const auto& [key, val] : m) os << sep << key << ": " << val, sep = ", ";
    return os << '}';
}

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
