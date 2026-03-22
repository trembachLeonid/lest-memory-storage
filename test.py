import socket
from time import sleep
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('127.0.0.1', 8090))

while True:
    sock.sendall(b"PING\n")

    response = sock.recv(1024)
    print(f"Received from server: {response.decode('utf-8').strip()}")
    sleep(10)
    
