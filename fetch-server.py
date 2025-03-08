import http.server
import socketserver
import json
import os
from socketserver import ThreadingMixIn
import shutil

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