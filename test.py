import socket
from time import sleep
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('127.0.0.1', 8090))

sock.sendall(b"SET key 10\n")

response = sock.recv(1024)
print(f"Received from server: {response.decode('utf-8').strip()}")

for i in range(3):

    sock.sendall(b"INC key\n")

    response = sock.recv(1024)
    print(f"Received from server: {response.decode('utf-8').strip()}")

    sock.sendall(b"INC key 10\n")
    response = sock.recv(1024)
    print(f"Received from server: {response.decode('utf-8').strip()}")
    sleep(2)

sock.sendall(b"QUIT\n")

response = sock.recv(1024)
print(f"Received from server: {response.decode('utf-8').strip()}")

