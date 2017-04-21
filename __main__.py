from optparse import OptionParser
import server

if __name__ == '__main__':
    parser = OptionParser()
    parser.add_option('-p', '--port', action='store', type='int',
        dest='port', default=5001)
    (options, args) = parser.parse_args()
    server.serve(options.port)
