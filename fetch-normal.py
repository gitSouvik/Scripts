import http.server
import socketserver
import json
import os
from socketserver import ThreadingMixIn
import subprocess  # For opening the .cpp file
import re

PORT = 54321  # Replace with your custom port number

# Define the handler to process the data sent by Competitive Companion
class Handler(http.server.SimpleHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])  # Get the size of the data
        body = self.rfile.read(content_length)  # Read the data
        data = json.loads(body)  # Parse the JSON data
        
        # Extract the problem name and the number of tests
        problem_name = data['name']  # Name of the problem (for lowercase use - .lower())



        if 'y' in output.lower():
            # Change the output if 'no'
            # Replace spaces with hyphens but avoid adding them around punctuation
            problem_name = re.sub(r'\s+(?=\w)', '-', problem_name.strip())  # For USACO problems only       
        else:
            # Change the output if 'no'
            problem_name = problem_name[0]  # For contest problems only


 
        test_cases = data['tests']  # List of test cases
        problem_link = data['url']  # URL link for the problem



        # Use the current working directory
        cwd = os.getcwd()

        # Iterate over each test case and save them with numbered filenames
        for i, test_case in enumerate(test_cases, start=1):
            # Generate filenames with test case numbers
            input_filename = os.path.join(cwd, f"{problem_name}-{i}.in")
            output_filename = os.path.join(cwd, f"{problem_name}-{i}.out")
            
            # Get input and output data
            input_data = test_case['input']
            output_data = test_case['output']

            # Save the input data
            with open(input_filename, 'w') as f:
                f.write(input_data)
            
            # Save the output data
            with open(output_filename, 'w') as f:
                f.write(output_data)



        # Save only the first input/output as {problem}-1.in and {problem}-1.out             ***** PROBLEM : ONLY FIRST TEST CASE *****

        # input_filename = os.path.join(cwd, f"{problem_name}-1.in")
        # output_filename = os.path.join(cwd, f"{problem_name}-1.out")
        # input_data = test_cases[0]['input']
        # output_data = test_cases[0]['output']

        # # Save the input data
        # with open(input_filename, 'w') as f:
        #     f.write(input_data)
        
        # # Save the output data
        # with open(output_filename, 'w') as f:
        #     f.write(output_data)
        



        # Create a .cpp file for the problem if it doesn't exist
        cpp_filename = os.path.join(cwd, f"{problem_name}.cpp")
        if not os.path.exists(cpp_filename):
            with open(cpp_filename, 'w') as f:
                # Replace with your name and the problem link
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
                    f.write("#define _ << \" \" <<\n")
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
                    f.write("    int tc = 1, _tc = 1;\n")
                    f.write("    // cin >> tc;\n")
                    f.write("    while(tc-- > 0) {\n")
                    f.write("        // cerr << \"\\nCase #\" << _tc++ << \" :\\n\";\n")
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
                                
        # Open the .cpp file using Sublime Text (or change to your preferred editor)
        subprocess.run(['/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl', cpp_filename])        
        self.send_response(200)  # Send a response back to the client
        self.end_headers()

# Use ThreadingMixIn to handle requests in separate threads
class ThreadedTCPServer(ThreadingMixIn, socketserver.TCPServer):
    """Handle requests in a separate thread."""

# # Start the server
# with ThreadedTCPServer(("", PORT), Handler) as httpd:
#     print(f"Serving on port {PORT}")
#     output = input("usaco ? (y/n) : ")
#     httpd.serve_forever()  # Keep the server running

# Start the server [LATEST]
with ThreadedTCPServer(("", PORT), Handler) as httpd:
    print(f"Serving on port {PORT}")
    output = input("usaco ? (y/n) : ")
    try:
        httpd.serve_forever()  # Keep the server running
    except KeyboardInterrupt:
        print("\nServer stopped by user.")
    finally:
        httpd.server_close()  # Clean up
        print("Server closed.")