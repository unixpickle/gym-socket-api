"""
High-level code for listening for and handling client
connections.
"""

import os
import socket
import json

import proto
import gym
import numpy as np

def serve(port):
    """
    Run a server on the given port.
    """
    # TODO: use socketserver module.
    sock = socket.socket()
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.bind(('127.0.0.1', port))
    print('Listening on port %d...' % port)
    sock.listen(10)
    while True:
        conn, addr = sock.accept()
        pid = os.fork()
        if pid == 0:
            handle(conn, addr)
            exit()
        else:
            conn.close()

def handle(conn, addr):
    """
    Handle a connection from a client.
    """
    print('Connection from %s' % str(addr))
    sock_file = conn.makefile(mode='rwb')
    try:
        env = handshake(sock_file)
        try:
            loop(sock_file, env)
        finally:
            env.close()
    except proto.ProtoException as exc:
        print('%s gave error %s' % (str(addr), str(exc)))
    finally:
        print('Disconnect from %s' % str(addr))
        sock_file.close()

def handshake(sock):
    """
    Perform the initial handshake and return the resulting
    Gym environment.
    """
    flags = proto.read_flags(sock)
    if flags != 0:
        raise proto.ProtoException('unsupported flags: ' + str(flags))
    env_name = proto.read_field_str(sock)
    try:
        env = gym.make(env_name)
        proto.write_field_str(sock, '')
        sock.flush()
        return env
    except gym.error.Error as gym_exc:
        proto.write_field_str(sock, str(gym_exc))
        sock.flush()
        raise gym_exc

def loop(sock, env):
    """
    Handle commands from the client as they come in and
    apply them to the given Gym environment.
    """
    while True:
        pack_type = proto.read_packet_type(sock)
        if pack_type == 'reset':
            handle_reset(sock, env)
        elif pack_type == 'step':
            handle_step(sock, env)
        elif pack_type == 'get_space':
            handle_get_space(sock, env)
        elif pack_type == 'sample_action':
            handle_sample_action(sock, env)

def handle_reset(sock, env):
    """
    Reset the environment and send the result.
    """
    proto.write_obs(sock, env, env.reset())
    sock.flush()

def handle_step(sock, env):
    """
    Step the environment and send the result.
    """
    action = proto.read_action(sock)
    if isinstance(action, list):
        action = np.array(action)
    obs, rew, done, info = env.step(action)
    proto.write_obs(sock, env, obs)
    proto.write_reward(sock, rew)
    proto.write_bool(sock, done)
    proto.write_field_str(sock, json.dumps(info))
    sock.flush()

def handle_get_space(sock, env):
    """
    Get information about the action or observation space.
    """
    space_id = proto.read_space_id(sock)
    if space_id == 'action':
        proto.write_space(sock, env.action_space)
    elif space_id == 'observation':
        proto.write_space(sock, env.observation_space)
    sock.flush()

def handle_sample_action(sock, env):
    """
    Generate and send a random action.
    """
    action = env.action_space.sample()
    proto.write_action(sock, env, action)
    sock.flush()
