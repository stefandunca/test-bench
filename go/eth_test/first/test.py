import json
from typing import Tuple
import subprocess
import os
import time

def get_ganache_account(ganache_out: str, index: int) -> Tuple[str, str]:
    # Find the "Available Accounts" line index
    accounts_line_start_index = ganache_out.find("Available Accounts")
    if accounts_line_start_index == -1:
        raise Exception("Could not find 'Available Accounts' in ganache output")
    address_line = ganache_out[accounts_line_start_index:].splitlines()[2 + index]
    address_start = address_line.find("0x")
    address = address_line[address_start:address_line.find(" ", address_start)]

    private_key_line_index = ganache_out.find("Private Keys")
    if private_key_line_index == -1:
        raise Exception("Could not find 'Private Keys' in ganache output")
    private_key_line = ganache_out[private_key_line_index:].splitlines()[2 + index]
    private_key = private_key_line[private_key_line.find("0x"):]
    return address, private_key

def generate_test_file() -> str:
    address_list = []
    private_key_list = []
    for i in range(2):
        address, key = get_ganache_account(ganache_out.decode("utf-8"), i)
        address_list.append(address)
        private_key_list.append(key)

    json_str = json.dumps({'addresses': address_list, 'private_keys': private_key_list})

    json_file = os.path.join(script_dir, 'test_data.json')
    with open(json_file, 'w') as f:
        f.write(json_str)
    return json_file

ganache_args = ['ganache-cli', '--networkId', '1337', '-m', '"much repair shock carbon improve miss forget sock include bullet interest solution"']
ganache = subprocess.Popen(ganache_args, stdout=subprocess.PIPE)
ganache_out = b''
done_reading_ganache = False

script_dir = os.path.dirname(os.path.realpath(__file__))

tests_args = ['go', 'test', '-v', f'{os.path.join(script_dir, "/*.go")}']
tests = None
tests_out = b''
tests_done_time = None



try:
    while ganache.poll() is None:# and (tests_done_time is None or (time.time() - tests_done_time) < 1):

        # Handle ganache
        gOut = ganache.stdout.read1()
        if gOut != b'':
            ganache_out += gOut

        if not done_reading_ganache:
            if b'Listening on' in gOut:
                done_reading_ganache = True
                os.set_blocking(ganache.stdout.fileno(), False)

        # Handle tests
        if tests is None:
            if done_reading_ganache:
                json_file_path = generate_test_file()

                os.environ['TEST_DATA_FILE'] = json_file_path
                print(f'Test data "{json_file_path}"')

                tests = subprocess.Popen(tests_args, stdout=subprocess.PIPE, cwd=script_dir, env=os.environ)
        else:
            tOut = tests.stdout.readline()
            if tOut != b'':
                print(tOut.decode("utf-8"), end='')

            if tests.poll() is not None and tests_done_time is None:
                tests_done_time = time.time()

        time.sleep(0.01)
finally:
    print(ganache_out.decode("utf-8"), end='')

ganache.terminate()
tests.terminate()