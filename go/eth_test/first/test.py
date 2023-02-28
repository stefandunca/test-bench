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
tests_args = ['go', 'test']
tests = None
tests_out = b''
tests_start_time = None
tests_end_time = None

home_dir = os.path.dirname(os.path.expanduser('~/'))
stop_when_tests_end = True

try:
    while ganache.poll() is None and not stop_when_tests_end or tests_end_time is None:

        # Handle ganache
        gOut = ganache.stdout.read1()
        if gOut != b'':
            ganache_out += gOut
            print(f'@dd GANACHE: {gOut.decode("utf-8")}', end='')

        if not done_reading_ganache:
            if b'Listening on' in gOut:
                done_reading_ganache = True
                os.set_blocking(ganache.stdout.fileno(), False)

        # Handle tests
        if tests is None:
            if done_reading_ganache:
                # ganache.terminate()
                json_file_path = generate_test_file()

                os.environ['TEST_DATA_FILE'] = json_file_path
                os.environ['GOPATH'] = os.path.join(home_dir, "go")
                os.environ['GOBIN'] = os.path.join(os.environ['GOPATH'], 'bin')
                os.environ['PATH'] = os.environ['PATH'] + ":" + os.environ['GOBIN']
                go_root = subprocess.Popen(['go', 'env', 'GOROOT'], stdout=subprocess.PIPE, env=os.environ)
                os.environ['GOROOT'] = go_root.stdout.readline().decode("utf-8").strip()

                tests_start_time = time.time()
                tests = subprocess.Popen(tests_args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, cwd=script_dir, env=os.environ)
                os.set_blocking(tests.stdout.fileno(), False)
                os.set_blocking(tests.stderr.fileno(), False)
        else:
            tOut = tests.stdout.read1()
            if tOut != b'':
                print(f'@dd TESTS  : {tOut.decode("utf-8")}', end='')
            tOut = tests.stderr.read1()
            if tOut != b'':
                print(f'@dd TESTS  : {tOut.decode("utf-8")}', end='')

            if tests.poll() is not None and tests_end_time is None:
                tests_end_time = time.time()

            if stop_when_tests_end and tests_end_time is not None and (time.time() - tests_end_time) > 10:
                break

        time.sleep(0.01)
finally:
    print(f'DONE! tests run time (s): {tests_end_time - tests_start_time}')

if ganache:
    ganache.terminate()
if tests:
    tests.terminate()