import struct

class ProtoException(Exception):
    pass

def read_byte(sock):
    flagByte = sock.read(1)
    if len(flagByte) != 1:
        raise ProtoException('EOF')
    return struct.unpack('<B', flagByte)[0]

def read_flags(sock):
    return read_byte(sock)

def read_packet_type(sock):
    typeVal = read_byte(sock)
    raise ProtoException('unknown packet type: ' + str(typeVal))

def read_field(sock):
    lenBytes = sock.read(4)
    if len(lenBytes) != 4:
        raise ProtoException('EOF reading length field')
    lenVal = struct.unpack('<I', lenBytes)[0]
    res = sock.read(lenVal)
    if len(res) != lenVal:
        raise ProtoException('EOF reading field value')
    return res

def read_field_str(sock):
    return read_field(sock).decode('utf-8')

def write_field(sock, field):
    lenField = struct.pack('<I', len(field))
    sock.write(lenField)
    sock.write(field)

def write_field_str(sock, field):
    write_field(sock, field.encode('utf-8'))
