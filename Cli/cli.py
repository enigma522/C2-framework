import argparse
import requests
import os

def login(server_url):
    try:
        secret = "e7bcc0ba5fb1dc9cc09460baaa2a6986"  # Replace with actual secret key input
        response = requests.post(f'{server_url}/login', json={'implantID': '47e040dd-cd6b-4463-a8de-00423b8b9c21', 'secret': secret})
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

def main():
    parser = argparse.ArgumentParser(description="A Python CLI to interact with the C2-server.")
    
    parser.add_argument(
        'action',
        nargs='?',
        choices=['implants', 'tasks', 'results'],
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
