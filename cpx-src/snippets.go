package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ── Code Snippets ─────────────────────────────────────────────────────────────

type snippet struct {
	name string
	code string
}

var defaultSnippets = []snippet{
	{
		name: "DSU (Union-Find)",
		code: `struct DSU {
    vector<int> p, rnk;
    DSU(int n) : p(n), rnk(n, 0) { iota(p.begin(), p.end(), 0); }
    int find(int x) { return p[x] == x ? x : p[x] = find(p[x]); }
    bool unite(int a, int b) {
        a = find(a); b = find(b);
        if (a == b) return false;
        if (rnk[a] < rnk[b]) swap(a, b);
        p[b] = a;
        if (rnk[a] == rnk[b]) rnk[a]++;
        return true;
    }
    bool same(int a, int b) { return find(a) == find(b); }
};`,
	},
	{
		name: "Segment Tree (range sum)",
		code: `struct SegTree {
    int n;
    vector<long long> t;
    SegTree(int n) : n(n), t(2 * n, 0) {}
    void update(int i, long long val) {
        for (t[i += n] = val; i > 1; i >>= 1)
            t[i >> 1] = t[i] + t[i ^ 1];
    }
    long long query(int l, int r) { // [l, r)
        long long res = 0;
        for (l += n, r += n; l < r; l >>= 1, r >>= 1) {
            if (l & 1) res += t[l++];
            if (r & 1) res += t[--r];
        }
        return res;
    }
};`,
	},
	{
		name: "Fenwick Tree (BIT)",
		code: `struct BIT {
    int n;
    vector<long long> t;
    BIT(int n) : n(n), t(n + 1, 0) {}
    void update(int i, long long delta) {
        for (++i; i <= n; i += i & -i) t[i] += delta;
    }
    long long query(int i) { // prefix sum [0, i]
        for (++i; i > 0; i -= i & -i) s += t[i];
        return s;
    }
    long long range(int l, int r) { return query(r) - (l ? query(l - 1) : 0); }
};`,
	},
	{
		name: "Dijkstra (weighted, 0-indexed)",
		code: `using pli = pair<long long, int>;
vector<long long> dijkstra(int src, int n, vector<vector<pli>>& adj) {
    vector<long long> dist(n, LLONG_MAX);
    priority_queue<pli, vector<pli>, greater<>> pq;
    dist[src] = 0;
    pq.push({0, src});
    while (!pq.empty()) {
        auto [d, u] = pq.top(); pq.pop();
        if (d > dist[u]) continue;
        for (auto [w, v] : adj[u])
            if (dist[u] + w < dist[v]) {
                dist[v] = dist[u] + w;
                pq.push({dist[v], v});
            }
    }
    return dist;
}`,
	},
	{
		name: "Binary Search template",
		code: `// Finds smallest x in [lo, hi] where cond(x) is true.
// Replace cond with your predicate.
auto bin_search = [&](long long lo, long long hi) -> long long {
    while (lo < hi) {
        long long mid = lo + (hi - lo) / 2;
        if (cond(mid)) hi = mid;
        else lo = mid + 1;
    }
    return lo; // lo == hi is the answer
};`,
	},
	{
		name: "Modular Arithmetic",
		code: `const long long MOD = 1e9 + 7;
long long power(long long b, long long e, long long m = MOD) {
    long long res = 1; b %= m;
    for (; e > 0; e >>= 1) {
        if (e & 1) res = res * b % m;
        b = b * b % m;
    }
    return res;
}
long long inv(long long a, long long m = MOD) { return power(a, m - 2, m); }`,
	},
	{
		name: "Sparse Table (RMQ, O(1) query)",
		code: `struct SparseTable {
    int n, LOG;
    vector<vector<int>> t;
    vector<int> lg;
    SparseTable(vector<int>& a) : n(a.size()), LOG(__lg(a.size()) + 2),
        t(LOG, vector<int>(a.size())), lg(a.size() + 1) {
        t[0] = a;
        for (int i = 2; i <= n; i++) lg[i] = lg[i / 2] + 1;
        for (int j = 1; j < LOG; j++)
            for (int i = 0; i + (1 << j) <= n; i++)
                t[j][i] = min(t[j-1][i], t[j-1][i + (1 << (j-1))]);
    }
    int query(int l, int r) { // [l, r] inclusive
        int k = lg[r - l + 1];
        return min(t[k][l], t[k][r - (1 << k) + 1]);
    }
};`,
	},
}



var snippets []snippet

func loadAllSnippets(cwd string) {
	snippets = append([]snippet{}, defaultSnippets...)
	file := filepath.Join(cwd, ".cpx_snippets.json")
	if data, err := os.ReadFile(file); err == nil {
		var customs []snippet
		if err := json.Unmarshal(data, &customs); err == nil {
			snippets = append(snippets, customs...)
		}
	}
}

func saveCustomSnippet(cwd, name, code string) error {
	file := filepath.Join(cwd, ".cpx_snippets.json")
	var customs []snippet
	if data, err := os.ReadFile(file); err == nil {
		json.Unmarshal(data, &customs)
	}
	customs = append(customs, snippet{name: name, code: code})
	data, err := json.MarshalIndent(customs, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(file, data, 0644); err != nil {
		return err
	}
	loadAllSnippets(cwd) // reload
	return nil
}
