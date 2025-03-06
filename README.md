# Competitive Programming Setup

## Initial Setup 

**Version 1.0**
Originally, the system required running [fetch-server.sh](https://github.com/gitSouvik/Scripts/blob/main/fetch-server.py) manually before each session. Inputs and outputs had to be added one by one using **Competitive Companion**, making the process time-consuming and inefficient. This version stored input and output in the following files:

* input.txt <br>
* output.txt <br>
* expected.txt <br>
* error.txt

>  **Drawback :** Cannot handle multiple file input and output at once (Since the previous file gets re-written and the input/output are lost). Need to fetch them every time individually for each problem.

 **Version 2.0**
 
A revised version of fetch.py improved test case organization by storing them in separate **.in** and **.out** files and also creating and automatically opening of **.cpp** file with that problem name in the **Editor**. 

  To solve the issue of storing multiple test cases in a single file, a **new-fetch mechanism** was implemented. This version saved test cases separately using problem-specific .in and .out files:
  
* A.in &nbsp;/&nbsp; A.out <br>
* B.in &nbsp;/&nbsp; B.out <br>
* C.in &nbsp;/&nbsp; C.out 

>  **Drawback :** If there were multiple test-cases for a single problem, it can only store **1** test case - 'the first one'. Also filenames generated from problem titles occasionally contain **extra spaces** which cannot be accessed in terminal (e.g., “cow tipping.cpp” may need to be manually corrected to “cow-tipping.cpp”).

**Version 3.0** (Latest)  

The latest iteration of the fetch system [fetch-unique-name.sh](https://github.com/gitSouvik/Scripts/blob/main/fetch-unique-name.py) resolves previous limitations. Previously, the script only retrieved the first test case from a problem statement. This issue has now been corrected. The modified script can now handle multiple test cases and save them as individual .in and .out files, such as :

* A-1.in &nbsp;/&nbsp; A-1.out <br>
* A-2.in &nbsp;/&nbsp; A-2.out <br>
* and **so on**

Spaces in filenames are now replaced with "-" to ensure smooth compatibility with **USACO** and other platforms.

* "***Problem 2. Cow Tipping.cpp***" will be renamed to "***Problem-2.-Cow-Tipping.cpp***"

> **Bonus** : As the name suggests, it was introduced to **prevent filename conflict**s. If a file named A.cpp already exists, the system automatically creates A1.cpp, then A2.cpp, and so forth, ensuring that existing files are not overwritten. <br>

## Automated File Creation

A new script, [create.sh](https://github.com/gitSouvik/Scripts/blob/main/create.sh), has been added to automate the creation of files with different extensions. This allows users to quickly generate any number of .cpp, .py, .txt, or other necessary files.

## Runsamples & Debug

**runsamples.sh & runsamples-debug.sh**
Both the scripts serves the same function i.e. to show output in a specific format in the terminal. But there functionalities and purpose differs a bit.

 > **Important** : Let us, denote ***runsamples.sh*** as **[run.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples.sh)** and ***runsamples-debug.sh*** as **[rrun.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples-debug.sh)**.

**Key Differences:**

1. **Test Case Handling:** <br> <br>
•  **[run.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples.sh)** :  This script loops through all test cases (A-1.in, A-2.in, ...) and executes the compiled program for each input file while checking memory and time usage. <br> 
•  **[rrun.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples-debug.sh)** : This script allows executing a specific test case (by passing the test number as an argument) or running all test cases if N=0.

2. **Output Handling:** <br> <br>
•  **[run.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples.sh)** : Stores the program’s output in A-1-output.out, A-2-output.out, etc., and compares it with expected outputs. <br>
•  **[rrun.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples-debug.sh)** : Displays the output in the terminal but does not compare it with expected outputs.

3. **Debugging and Error Handling:** <br> <br>
•  **[run.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples.sh)** : This script explicitly checks for segmentation faults and highlights them. <br>
•  **[rrun.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples-debug.sh)** : This stores error logs in a temporary file (temp_debug.log) and displays them separately.

4. **Execution Modes:** <br> <br>
• **[run.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples.sh)** is designed for automated testing and detailed logging. <br>
•  **[rrun.sh](https://github.com/gitSouvik/Scripts/blob/main/runsamples-debug.sh)** is designed for quick execution with an option to run specific test cases.

----------
