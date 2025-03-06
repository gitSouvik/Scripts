import http.server
import socketserver
import json
import os
from socketserver import ThreadingMixIn

PORT = 54321  # Replace with your custom port number

# Define the handler to process the data sent by Competitive Companion
class Handler(http.server.SimpleHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])  # Get the size of the data
        body = self.rfile.read(content_length)  # Read the data
        data = json.loads(body)  # Parse the JSON data
        
        # Debug: Print the incoming data to check its structure
        # print("Received Data:", json.dumps(data, indent=4))  # Pretty print the incoming JSON
        
        try:
            # Extract problem information directly from the top-level keys
            problem_name = data['name']
            problem_link = data['url']
            
            # Extract input and output from the tests
            input_data = data['tests'][0]['input']
            output_data = data['tests'][0]['output']
        except KeyError as e:
            print(f"Missing key: {e}")  # Log the missing key
            self.send_response(400)  # Bad Request
            self.end_headers()
            return
        
        # Define the path where files will be saved
        desktop = os.path.join(os.environ['HOME'], 'Desktop')        
        cp_folder = os.path.join(desktop, 'CP')
        io_folder = os.path.join(cp_folder, 'io')

        # Make sure the 'CP' folder exists
        os.makedirs(cp_folder, exist_ok=True)  # Create the folder if it doesn't exist
        os.makedirs(io_folder, exist_ok=True)  # Create the 'io' folder if it doesn't exist
        
        # Save the input data
        with open(os.path.join(io_folder, 'input.txt'), 'w') as f:
            f.write(input_data)
        
        # Save the output data
        with open(os.path.join(io_folder, 'expected.txt'), 'w') as f:
            f.write(output_data)
        
        # Generate a new C++ file name
        cpp_filename = self.get_unique_cpp_filename(cp_folder, problem_name)
        
        # Create the C++ file with the template if it doesn't exist
        with open(cpp_filename, 'w') as f:
            f.write("/*\n")
            f.write(" * Author: Calypsoo\n")
            f.write(f" * Problem: {problem_name}\n")
            f.write(f" * P-link: {problem_link}\n")
            f.write(" */\n\n")

            f.write("#include <bits/stdc++.h>\n")
            f.write("using namespace std;\n\n")
            f.write("#define all(x) (x).begin(), (x).end()\n")
            f.write("#define rall(x) (x).rbegin(), (x).rend()\n")
            f.write("#define umap unordered_map\n")
            f.write("#define uset unordered_set\n")
            f.write("#define pb push_back\n")
            f.write("#define eb emplace_back\n")
            f.write("#define lb lower_bound\n")
            f.write("#define ub upper_bound\n")
            f.write("#define contain(map, i) (map.find(i) != map.end())\n")
            f.write("#define maxele max_element\n")
            f.write("#define minele min_element\n")
            f.write("#define len(s) (s).length()\n")
            f.write("#define nl cout << \"\\n\"\n")
            f.write("#define rep(i, a, b) for(int i = a; i < b; i++)\n")
            f.write("// #define int long long\n\n")

            f.write("using ll = long long;\n")
            f.write("using ld = long double;\n")
            f.write("constexpr ll mod = 1e9 + 7;\n\n")

            f.write("inline void apv(int ans = 1, int cap = 1) {\n")
            f.write("    if(ans == -1) {\n")
            f.write("        cout << (-1) << endl;\n")
            f.write("    } else {\n")
            f.write("        cout << (cap ? (ans ? \"YES\" : \"NO\") : (ans ? \"Yes\" : \"No\"));\n")
            f.write("        cout << endl;\n")
            f.write("    }\n")
            f.write("}\n\n")

            f.write("#ifdef WOOF_\n")
            f.write("#include <bits/debug.h>\n")
            f.write("#else\n")
            f.write("#define dbg(...) 25\n")
            f.write("#endif\n\n")

            f.write("/*\n")
            f.write(" * In the name of GOD\n")
            f.write(" * Here we go!\n")
            f.write(" */\n\n")

            f.write("void solve() {\n")
            f.write("    int n; cin >> n;\n")
            f.write("    vector<int> a(n);\n")
            f.write("    for(int i = 0; i < n; i++) {\n")
            f.write("        cin >> a[i];\n")
            f.write("    }\n")
            f.write("}\n\n")

            f.write("int32_t main() {\n")
            f.write("    ios_base::sync_with_stdio(false);\n")
            f.write("    cin.tie(0); cout.tie(0);\n")
            f.write("    // freopen(\"cmd+D.in\", \"r\", stdin);\n")
            f.write("    // freopen(\"cmd+D.out\", \"w\", stdout);\n")
            f.write("    int tc = 1, _ = 1;\n")
            f.write("    // cin >> tc;\n")
            f.write("    while(tc-- > 0) {\n")
            f.write("        // cerr << \"\\nCase #\" << _++ << \" :\\n\";\n")
            f.write("        solve();\n")
            f.write("    }\n")
            f.write("    return 0;\n")
            f.write("}")

            f.write("""
/*  
 * 1. ALWAYS THINK SIMPLE 
 * 2. IF IT GETS COMPLICATED THINK AGAIN :) FROM SCRATCH !
 * 3. SPEND ABOUT THE SAME AMOUNT OF TIME THAT YOU WOULD BE ABLE 
 *    TO DURING A REAL CONTEST
 */\n
""")
        
        self.send_response(200)  # Send a response back to the client
        self.end_headers()

    def get_unique_cpp_filename(self, cp_folder, base_name):
        """Generate a unique C++ file name (A.cpp, A2.cpp, A3.cpp, ...)"""
        i = 1
        cpp_filename = os.path.join(cp_folder, f"{base_name}.cpp")
        
        # Check if the file exists, if so, increment the counter
        while os.path.exists(cpp_filename):
            i += 1
            cpp_filename = os.path.join(cp_folder, f"{base_name}{i}.cpp")
        
        return cpp_filename

# Use ThreadingMixIn to handle requests in separate threads
class ThreadedTCPServer(ThreadingMixIn, socketserver.TCPServer):
    """Handle requests in a separate thread."""

# Start the server
with ThreadedTCPServer(("", PORT), Handler) as httpd:
    print(f"Serving on port {PORT}")
    httpd.serve_forever()  # Keep the server running