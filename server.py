import os
import socket
import traceback
import struct
import proto
import gym

def serve(port):
    # TODO: use socketserver module.
    s = socket.socket()
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.bind(('127.0.0.1', port))
    print('Listening on port %d...' % port)
    s.listen(10)
    while True:
        conn, addr = s.accept()
        pid = os.fork()
        if pid == 0:
            handle(conn, addr)
            exit()
        else:
            conn.close()

def handle(conn, addr):
    print('Connection from %s' % str(addr))
    f = conn.makefile(mode='rwb')
    try:
        env = handshake(f)
        try:
            loop(f, env)
        finally:
            env.close()
    except Exception as e:
        print('Error from %s' % str(addr))
        traceback.print_exc()
    finally:
        f.close()

def handshake(sock):
    flags = proto.read_flags(sock)
    envName = proto.read_field_str(sock)
    try:
        env = gym.make(envName)
        proto.write_field_str(sock, '')
        sock.flush()
        return env
    except gym.error.Error as e:
        proto.write_field_str(sock, str(e))
        sock.flush()
        raise e

def loop(sock, env):
    while True:
        packType = proto.read_packet_type(sock)
        # TODO: use packType here.
        print('Packet of type ' + packType)
