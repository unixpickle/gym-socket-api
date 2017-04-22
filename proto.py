"""
Low-level API for protocol-specific encoding/decoding.
"""

import struct
import json

class ProtoException(Exception):
    """
    Exception type used for all protocol-related errors.
    """
    pass

def read_byte(sock):
    """
    Read a byte from the socket.
    """
    data = sock.read(1)
    if len(data) != 1:
        raise ProtoException('EOF')
    return struct.unpack('<B', data)[0]

def read_flags(sock):
    """
    Read handshake flags from the socket.
    """
    return read_byte(sock)

def read_packet_type(sock):
    """
    Read packet type from the socket and turn it into a
    human-readable string.
    """
    type_id = read_byte(sock)
    if type_id == 0:
        return 'reset'
    elif type_id == 1:
        return 'step'
    raise ProtoException('unknown packet type: ' + str(type_id))

def read_field(sock):
    """
    Read a variable length data field.
    """
    len_data = sock.read(4)
    if len(len_data) != 4:
        raise ProtoException('EOF reading length field')
    length = struct.unpack('<I', len_data)[0]
    res = sock.read(length)
    if len(res) != length:
        raise ProtoException('EOF reading field value')
    return res

def read_field_str(sock):
    """
    Read a variable length string field.
    """
    return read_field(sock).decode('utf-8')

def write_field(sock, field):
    """
    Write a variable length data field.
    """
    sock.write(struct.pack('<I', len(field)))
    sock.write(field)

def write_field_str(sock, field):
    """
    Write a variable length string field.
    """
    write_field(sock, field.encode('utf-8'))

def write_obs_json(sock, jsonable):
    """
    Write a JSON observation object.
    """
    sock.write(struct.pack('<B', 0))
    write_field_str(sock, json.dumps(jsonable, separators=(',', ':')))

def write_obs_byte_list(sock, arr):
    """
    Write a byte list observation from a numpy array.
    """
    sock.write(struct.pack('<B', 1))
    dims = list(arr.shape)
    header = struct.pack('<I', len(dims))
    for dim in dims:
        header += struct.pack('<I', dim)
    payload = arr.tobytes()
    sock.write(struct.pack('<I', len(header)+len(payload)))
    sock.write(header)
    sock.write(payload)

def write_reward(sock, rew):
    """
    Write a reward value.
    """
    sock.write(struct.pack('<d', rew))

def write_bool(sock, flag):
    """
    Write a boolean.
    """
    num = 0
    if flag:
        num = 1
    sock.write(struct.pack('<B', num))

def read_action(sock):
    """
    Read an action object.
    """
    type_id = read_byte(sock)
    if type_id == 0:
        return json.loads(read_field_str(sock))
    raise ProtoException('unknown action type: ' + str(type_id))
