import http.server
import socketserver
import json
import os
from socketserver import ThreadingMixIn
import subprocess  # For opening the .cpp file
import re
import shutil

PORT = 54321  # Replace with your custom port number

# Define the handler to process the data sent by Competitive Companion
class Handler(http.server.SimpleHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])  # Get the size of the data
        body = self.rfile.read(content_length)  # Read the data
        data = json.loads(body)  # Parse the JSON data
        
        # Extract the problem name and the number of tests
        problem_name = data['name']  # Name of the problem (for lowercase use - .lower())
        test_cases = data['tests']  # List of test cases
        problem_link = data['url']  # URL link for the problem

        if 'y' in output.lower():
            # Replace spaces with hyphens but avoid adding them around punctuation
            problem_name = problem_name.strip().replace(" ", "-")  # For USACO problems only       
        else:
            problem_name = problem_name[0]  # For contest problems only

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
        template_path = '/Users/woofwoof/scripts/template.cpp'

        if not os.path.exists(cpp_filename):
            with open(cpp_filename, 'w') as f:
                # Add author and problem link 
                f.write("/*\n")
                f.write(" * Author: Calypsoo\n")
                f.write(f" * Problem: {problem_name}\n")
                f.write(f" * P-link: {problem_link}\n")
                f.write(" */\n\n")
                
                # Copy the cp template file content   
                if os.path.exists(template_path):
                    with open(template_path, 'r') as template_file:
                        shutil.copyfileobj(template_file, f)  
                                                
        # Open the .cpp file using Sublime Text (or change to your preferred editor)
        subprocess.run(['/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl', cpp_filename])        
        self.send_response(200)  # Send a response back to the client
        self.end_headers()



# Use ThreadingMixIn to handle requests in separate threads
class ThreadedTCPServer(ThreadingMixIn, socketserver.TCPServer):
    """Handle requests in a separate thread."""


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