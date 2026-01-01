import socket

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(("localhost", 8000))
while True:
    msg = input("Enter a message: ") + '\n'
    sock.send(bytes(msg, 'utf-8'))
