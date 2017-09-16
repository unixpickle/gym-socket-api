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
    parser.add_argument('-p', '--port', action='store', type=int,
                        dest='port', default=5001)
    parser.add_argument('-u', '--universe', action='store_true',
                        dest='universe')
    parser.add_argument('-s', '--setup', action='store', type=str,
                        dest='setup_code')
    options = parser.parse_args()
    server.serve(**vars(options))

if __name__ == '__main__':
    main()
