#include <bits/stdc++.h>
using namespace std;

#define all(x) (x).begin(), (x).end()
#define umap unordered_map
#define mset multiset
#define pb push_back
#define eb emplace_back
#define len(s) (s).length()
#define nl cout << "\n"
#define _ << " " <<
#define rep(i, a, b) for(int i = a; i < b; i++)
#define fun(ret, ...) std::function<ret(__VA_ARGS__)>
// #define int long long

using ll = long long;
using ld = long double;
constexpr int mod = 1000000007; // 998244353;

inline void apv(int ans = 1, int cap = 0) {
    cout << (ans == -1 ? "-1" : (cap ? (ans ? "YES" : "NO") : (ans ? "Yes" : "No"))) << "\n";
}

void USACO(const string& file) {
    freopen((file + ".in").c_str(), "r", stdin), freopen((file + ".out").c_str(), "w", stdout);
}

#ifdef WOOF_
#include <bits/debug.h>
#else
#define dbg(...) 25
#endif

/*  In the name of GOD, here we go!
    * Think (SIMPLE)
    * Complicated ? (START AGAIN) from SCRATCH
    * Spend about the (SAME AMOUNT OF TIME) that you would be able to DURING A REAL CONTEST
*/

const int N = 2e5+4;

void solve() {
    int n; cin >> n;
    vector<int> a(n);
    for(int i = 0; i < n; ++i) {
        cin >> a[i];
    }
}

int32_t main() {
    ios_base::sync_with_stdio(false);
    cin.tie(0); cout.tie(0);
    // USACO("");
    int tc = 1, _tc = 1;
    // cin >> tc;
    while(tc-- > 0) {
        // cerr << "\nCase #" << _tc++ << " :\n";
        solve();
    }
    return 0;
}