import socket

HOST = '127.0.0.1'
PORT = 8081
s = socket.socket()
s.connect((HOST, PORT))

s.send("Hello".encode())
data = s.recv(1024)
print(data.decode())
