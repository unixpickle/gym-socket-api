"""
Command-line tool for serving gym-socket-api.
"""

from argparse import ArgumentParser
import server

def main():
    """
    Parse command-line arguments and invoke server.
    """
    parser = ArgumentParser()
    parser.add_argument('-p', '--port', action='store', type='int',
                        dest='port', default=5001)
    options = parser.parse_args()[0]
    server.serve(options.port)

if __name__ == '__main__':
    main()
