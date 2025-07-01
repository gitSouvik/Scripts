#include <bits/stdc++.h>
using namespace std;
#define all(x) (x).begin(), (x).end()
#define umap unordered_map
#define pb push_back
#define nline "\n"
#define rep(i, a, b) for(int i = a; i < b; ++i)
#define fun(ret, ...) std::function<ret(__VA_ARGS__)>
#define int long long
#ifdef WOOF_
#include <bits/debug.h>
#else
#define dbg(...) 25
#endif

mt19937_64 RNG(chrono::steady_clock::now().time_since_epoch().count());
int rand(int l, int r) { return uniform_int_distribution<int>(l, r)(RNG); }
void apv(int a = 1, int c = 1) {
    cout << (a == -1 ? "-1" : a ? (c ? "YES" : "Yes") : (c ? "NO" : "No")) << '\n';
}
void USACO(const string& file) {
    freopen((file + ".in").c_str(), "r", stdin), freopen((file + ".out").c_str(), "w", stdout); 
}
void chmin(int &a, int b) { a = min(a, b); }
void chmax(int &a, int b) { a = max(a, b); }

const int mod = 1000000007; // 998244353;
const int N = 1000006, maxN = 3005, inf = 1e9+5;
int a[N], pre[N];


void solve() {

}

int32_t main() {
    auto start = std::chrono::high_resolution_clock::now();
    ios_base::sync_with_stdio(0); cin.tie(0); 
    // USACO("code");
    int tc = 1; 
    // cin >> tc;
    for(int i = 1; i <= tc; ++i) {
        // cerr << "\nCase #" << i << " :\n";
        solve();
    }
    auto end = std::chrono::high_resolution_clock::now();
    std::chrono::duration<long double> elapsed = end - start;
    cerr << "Time measured: " << fixed << setprecision(7) << elapsed.count() << " sec";
    return 0;
}