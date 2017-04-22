"""
Low-level API for protocol-specific encoding/decoding.
"""

import struct
import json
from gym import spaces
import numpy as np

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
    elif type_id == 2:
        return 'get_space'
    elif type_id == 3:
        return 'sample_action'
    elif type_id == 4:
        return 'monitor'
    elif type_id == 5:
        return 'render'
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

def write_obs(sock, env, obs):
    """
    Encode and send an observation.
    """
    if isinstance(obs, np.ndarray):
        if obs.dtype == 'uint8':
            write_obs_byte_list(sock, obs)
            return
    jsonable = env.observation_space.to_jsonable(obs)
    write_obs_json(sock, jsonable)

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

def read_bool(sock):
    """
    Read a boolean.
    """
    flag = read_byte(sock)
    if flag == 0:
        return False
    elif flag == 1:
        return True
    raise ProtoException('invalid boolean: ' + str(flag))

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

def write_action(sock, env, action):
    """
    Write an action object.
    """
    jsonable = env.action_space.to_jsonable(action)
    sock.write(struct.pack('<B', 0))
    write_field_str(sock, json.dumps(jsonable))

def write_space(sock, space):
    """
    Encode and write a gym.Space.
    """
    write_field_str(sock, json.dumps(space_json(space)))

def space_json(space):
    """
    Encode a gym.Space as JSON.
    """
    if isinstance(space, spaces.Box):
        # JSON doesn't support infinity.
        bound = 1e30
        return {
            'type': 'Box',
            'shape': space.shape,
            'low': np.clip(space.low, -bound, bound).flatten().tolist(),
            'high': np.clip(space.high, -bound, bound).flatten().tolist()
        }
    elif isinstance(space, spaces.Discrete):
        return {
            'type': 'Discrete',
            'n': space.n
        }
    elif isinstance(space, spaces.MultiBinary):
        return {
            'type': 'MultiBinary',
            'n': space.n
        }
    elif isinstance(space, spaces.MultiDiscrete):
        return {
            'type': 'MultiDiscrete',
            'low': space.low.tolist(),
            'high': space.high.tolist()
        }
    elif isinstance(space, spaces.Tuple):
        return {
            'type': 'Tuple',
            'subspaces': [space_json(sub) for sub in space.spaces]
        }
    else:
        raise ProtoException('unknown space type: ' + str(type(space)))

def read_space_id(sock):
    """
    Read a space ID and convert it to a string.
    """
    space_id = read_byte(sock)
    if space_id == 0:
        return 'action'
    elif space_id == 1:
        return 'observation'
    raise ProtoException('unknown space ID: ' + str(space_id))
