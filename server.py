"""
Listen for client connections and dispatch handlers.
"""

import os
import sys
import subprocess
import socket

if sys.version_info >= (3, 0):
    import socketserver
else:
    import SocketServer as socketserver

def serve(port=5001, universe=False, setup_code=''):
    """
    Run a server on the given port.
    """
    server = Server(('127.0.0.1', port), Handler)
    server.universe = universe
    server.setup_code = setup_code
    print('Listening on port ' + str(port) + '...')
    server.serve_forever()

class Server(socketserver.ThreadingMixIn, socketserver.TCPServer):
    """
    The connection server.
    """
    allow_reuse_address = True
    universe = False
    setup_code = ''

class Handler(socketserver.BaseRequestHandler):
    """
    The connection handler.
    """
    def handle(self):
        script_file = os.path.join(os.path.dirname(__file__), 'handler.py')
        args = [
            sys.executable,
            script_file,
            '--addr',
            str(self.client_address),
            '--fd',
            str(self.request.fileno()),
            '--setup',
            str(self.server.setup_code)
        ]

        if self.server.universe:
            args.append('--universe')

        # Greatly reduces latency on Linux.
        if sys.platform in ['linux', 'linux2', 'darwin']:
            self.request.setsockopt(socket.IPPROTO_TCP, socket.TCP_NODELAY, 1)

        try:
            print('Connection from ' + str(self.client_address))
            if sys.version_info >= (3, 2):
                proc = subprocess.Popen(args,
                                        stdin=sys.stdin,
                                        stdout=sys.stdout,
                                        stderr=sys.stderr,
                                        pass_fds=(self.request.fileno(),))
            else:
                proc = subprocess.Popen(args,
                                        stdin=sys.stdin,
                                        stdout=sys.stdout,
                                        stderr=sys.stderr,
                                        close_fds=False)
            proc.wait()
        finally:
            print('Disconnected from ' + str(self.client_address))
