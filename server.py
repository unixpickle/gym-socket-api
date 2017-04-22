"""
Listen for client connections and dispatch handlers.
"""

import os
import sys
import subprocess

if sys.version_info >= (3, 0):
    import socketserver
else:
    import SocketServer as socketserver

def serve(port):
    """
    Run a server on the given port.
    """
    server = Server(('127.0.0.1', port), Handler)
    print('Listening on port ' + str(port) + '...')
    server.serve_forever()

class Server(socketserver.ThreadingMixIn, socketserver.TCPServer):
    """
    The connection server.
    """
    allow_reuse_address = True

class Handler(socketserver.BaseRequestHandler):
    """
    The connection handler.
    """
    def handle(self):
        script_file = os.path.join(os.path.dirname(__file__), 'handler.py')
        args = [
            sys.executable,
            script_file,
            '-u',
            str(self.client_address)
        ]

        if sys.version_info >= (3, 0):
            sock = self.request.makefile('rwb', buffering=0)
        else:
            sock = self.request.makefile('rwb', 0)

        try:
            print('Connection from ' + str(self.client_address))
            proc = subprocess.Popen(args, stdin=sock, stdout=sock,
                                    stderr=sys.stderr)
            proc.wait()
        finally:
            sock.close()
            print('Disconnected from ' + str(self.client_address))
